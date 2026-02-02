package server

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/docs"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository/postgres"
	"control-panel-service/internal/server/handler"
	"control-panel-service/internal/server/middleware"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg/telemetry"
	"fmt"
	"strconv"
	"time"

	pslq "control-panel-service/pkg/database/postgres"

	_ "control-panel-service/docs"

	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

func MetricsMiddleware() fiber.Handler {
	meter := telemetry.GetMeter()

	requestCounter, _ := meter.Int64Counter("http.requests.total")
	requestDuration, _ := meter.Int64Histogram("http.request.duration_ms")
	requestSize, _ := meter.Int64Histogram("http.request.size_bytes")
	responseSize, _ := meter.Int64Histogram("http.response.size_bytes")

	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestSizeBytes := int64(len(c.Request().Header.String()) + len(c.Body()))
		attrs := attribute.NewSet(
			attribute.String("method", c.Method()),
			attribute.String("route", c.Route().Path),
		)
		requestSize.Record(c.Context(), requestSizeBytes, metric.WithAttributeSet(attrs))

		err := c.Next()

		duration := time.Since(start)
		durationMs := duration.Milliseconds()

		statusCode := c.Response().StatusCode()
		attrs = attribute.NewSet(
			attribute.String("method", c.Method()),
			attribute.String("route", c.Route().Path),
			attribute.String("status_code", strconv.Itoa(statusCode)),
		)

		requestCounter.Add(c.Context(), 1, metric.WithAttributeSet(attrs))
		requestDuration.Record(c.Context(), durationMs, metric.WithAttributeSet(attrs))

		responseSizeBytes := int64(len(c.Response().Header.String()) + len(c.Response().Body()))
		responseSize.Record(c.Context(), responseSizeBytes, metric.WithAttributeSet(attrs))

		return err
	}
}

func Serve(ctx context.Context, cfg *config.Config) error {
	app := fiber.New(fiber.Config{
		BodyLimit:    cfg.Server.PostBodyLimit,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	// Global middlewares
	app.Use(cors.New())

	app.Use(MetricsMiddleware())

	// �🔒 Rate Limiter (GLOBAL)
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Server.RateLimitMaxRequest,         // max requests
		Expiration: cfg.Server.RateLimitExpirationDuration, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // rate limit per IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
	}))

	docs.SwaggerInfo.Host = cfg.Server.SwaggerHost
	docs.SwaggerInfo.Schemes = cfg.Server.SwaggerScheme

	swaggerHandler := fiberSwagger.New(fiberSwagger.Config{
		Title:                "Control Panel API",
		DeepLinking:          true,
		PersistAuthorization: true,
		DocExpansion:         "list",
		URL:                  cfg.Server.SwaggerDocJSON,
	})
	// Logging middleware should be early to capture all requests
	app.Use(middleware.LoggingMiddleware())

	eventProducer, err := provider.NewKafkaEventProducer(ctx, cfg)
	if err != nil {
		zap.L().Error("failed to create kafka event producer", zap.Error(err))

		return fmt.Errorf("create kafka event producer: %w", err)
	}
	defer func() {
		err := eventProducer.Close()
		if err != nil {
			zap.L().Error("failed to close event producer", zap.Error(err))
		}
	}()

	db, err := pslq.New(&pslq.DatabaseConfig{
		Host:               cfg.Postgres.Host,
		Port:               cfg.Postgres.Port,
		User:               cfg.Postgres.User,
		Password:           cfg.Postgres.Password,
		Database:           cfg.Postgres.Database,
		MaxOpenConnections: cfg.Postgres.MaxOpenConnections,
		LogLevel:           pslq.Silent,
	})
	if err != nil {
		zap.L().Error("failed to connect to postgres", zap.Error(err))

		return fmt.Errorf("could not connect to postgres: %w", err)
	}

	motherService := usecase.NewMotherService(db, postgres.NewMotherServiceRepository(db), eventProducer)
	motherHandler := handler.NewMotherServiceHandler(motherService)
	testCategoryHandler := handler.NewTestCategoryHandler(cfg, postgres.NewTestCategoryRepository(db))
	testScenarioUsecase := usecase.NewTestScenarioUsecase(
		db,
		postgres.NewTestScenarioRepository(db),
		postgres.NewTestCategoryRepository(db),
		postgres.NewTestServiceConfigRepository(db),
		postgres.NewMotherServiceRepository(db),
		usecase.NewInMemoryScenarioExecutorEngine(),
	)
	testScenarioHandler := handler.NewTestScenarioHandler(testScenarioUsecase)

	apiV1 := app.Group("/api/v1")

	// Mother service
	apiV1.Post("/mother-services", motherHandler.Create())
	apiV1.Get("/mother-services/:id", motherHandler.GetByID())
	apiV1.Post("/mother-services/search", motherHandler.GetPaginated())

	// test category
	apiV1.Get("/test-categories", testCategoryHandler.GetAll())
	apiV1.Get("/test-categories/:id", testCategoryHandler.GetByID())

	// test scenario
	apiV1.Post("/test-scenarios", testScenarioHandler.Create())
	apiV1.Get("/test-scenarios/:id", testScenarioHandler.GetByID())
	apiV1.Post("/test-scenarios/search", testScenarioHandler.GetPaginated())
	apiV1.Post("/test-scenarios/:id/start", testScenarioHandler.Start())

	// swagger endpoint
	apiV1.Get("/docs/*", swaggerHandler)

	zap.L().Info("fiber server starting",
		zap.Int("port", cfg.Server.Port),
	)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := app.Listen(addr); err != nil {
			errChan <- fmt.Errorf("fiber listen on %s failed: %w", addr, err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		zap.L().Info("shutting down server gracefully")
		// Gracefully shutdown the server
		if err := app.Shutdown(); err != nil {
			zap.L().Error("fiber shutdown failed", zap.Error(err))

			return fmt.Errorf("fiber shutdown failed: %w", err)
		}
		zap.L().Info("server shut down successfully")

		return nil
	case err := <-errChan:
		zap.L().Error("server error", zap.Error(err))

		return err
	}
}
