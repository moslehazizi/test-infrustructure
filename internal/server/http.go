package server

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/docs"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository/postgres"
	"control-panel-service/internal/server/handler"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg/telemetry"
	"fmt"
	"net/http"
	"strconv"
	"time"

	pslq "control-panel-service/pkg/database/postgres"
	kuber "control-panel-service/pkg/kubernetes"

	"github.com/gofiber/fiber/v2"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

		// nolint
		return err
	}
}

func Serve(ctx context.Context, cfg *config.Config) error {
	// app := fiber.New(fiber.Config{
	// 	BodyLimit:    cfg.Server.PostBodyLimit,
	// 	ReadTimeout:  cfg.Server.ReadTimeout,
	// 	WriteTimeout: cfg.Server.WriteTimeout,
	// })

	app := chi.NewRouter()

	// TODO
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      app,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Global middlewares
	// app.Use(cors.New())
	app.Use(cors.Handler(cors.Options{}))

	// nolint
	// app.Use(MetricsMiddleware())
	// TODO

	// �🔒 Rate Limiter (GLOBAL)
	// app.Use(limiter.New(limiter.Config{
	// 	Max:        cfg.Server.RateLimitMaxRequest,         // max requests
	// 	Expiration: cfg.Server.RateLimitExpirationDuration, // per minute
	// 	KeyGenerator: func(c *fiber.Ctx) string {
	// 		return c.IP() // rate limit per IP
	// 	},
	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
	// 			"error": "Too many requests, please try again later",
	// 		})
	// 	},
	// }))
	// TODO

	docs.SwaggerInfo.Host = cfg.Server.SwaggerHost
	docs.SwaggerInfo.Schemes = cfg.Server.SwaggerScheme

	// swaggerHandler := fiberSwagger.New(fiberSwagger.Config{
	// 	Title:                "Control Panel API",
	// 	DeepLinking:          true,
	// 	PersistAuthorization: true,
	// 	DocExpansion:         "list",
	// 	URL:                  cfg.Server.SwaggerDocJSON,
	// })
	// Logging middleware should be early to capture all requests
	// app.Use(middleware.LoggingMiddleware())
	// TODO

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

	kubernetes, err := kuber.New(ctx, &kuber.KubernConfig{NameSpace: cfg.Kubernetese.NameSpace})
	if err != nil {
		return fmt.Errorf("could not connect to kubernetes: %w", err)
	}

	provisionService := provider.NewProvisioningService(cfg, kubernetes)
	testServiceSDK := provider.NewSDKTestService(&http.Client{})
	testScenarioRepository := postgres.NewTestScenarioRepository(db)
	scenarioExecutorBuilder := usecase.NewScenarioExecutorBuilder()
	agentBuilder := usecase.NewTestAgentControllerToolBox(provisionService, testServiceSDK, cfg.Kubernetese.TestServiceAPPServe, cfg.Kubernetese.IngressHost, cfg.Kubernetese.IngressPort)
	stressTestExecutionManager := usecase.NewStressTestExecutionManager(agentBuilder, testScenarioRepository, scenarioExecutorBuilder)
	outboxRepository := postgres.NewOutboxRepository(db)
	motherServiceRepository := postgres.NewMotherServiceRepository(db, outboxRepository)

	motherService := usecase.NewMotherService(
		motherServiceRepository,
		testScenarioRepository,
		provider.NewProvisioningService(cfg, kubernetes),
		stressTestExecutionManager,
		cfg.Outbox.MaxAttempts,
	)

	// testScenarioUsecase := usecase.NewTestScenarioUsecase(
	// 	testScenarioRepository,
	// 	postgres.NewTestCategoryRepository(db),
	// 	postgres.NewTestServiceConfigRepository(db),
	// 	motherServiceRepository,
	// 	postgres.NewTestServiceRepository(db),
	// )
	// testScenarioOperationUsecase := usecase.NewTestScenarioOperationUsecase(
	// 	testScenarioRepository,
	// 	stressTestExecutionManager,
	// )
	motherHandler := handler.NewMotherServiceHandler(motherService)
	// testCategoryHandler := handler.NewTestCategoryHandler(cfg, postgres.NewTestCategoryRepository(db))

	// testScenarioHandler := handler.NewTestScenarioHandler(testScenarioUsecase)
	// testScenarioOperationHandler := handler.NewTestScenarioOperationHandler(testScenarioOperationUsecase)
	databaseMetadataService := usecase.NewDatabaseMetadata(postgres.NewDatabaseMetadataRepository(db, cfg))
	databaseMetadataHandler := handler.NewDatabaseMetadataHandler(databaseMetadataService)

	// apiV1 := app.Group("/api/v1")

	// TODO: Route or Group (Decision Needed)
	app.Route("/api/v1", func(app chi.Router) {
		// Mother service
		app.Post("/mother-services", motherHandler.Create)
		app.Get("/mother-services/{id}", motherHandler.GetByID)
		// app.Post("/mother-services/search", motherHandler.GetPaginated())
		// app.Post("/mother-services/:id/delete", motherHandler.Delete())

		// test category
		// app.Get("/test-categories", testCategoryHandler.GetAll())
		// app.Get("/test-categories/:id", testCategoryHandler.GetByID())

		// test scenario
		// app.Post("/test-scenarios", testScenarioHandler.Create())
		// app.Get("/test-scenarios/:id", testScenarioHandler.GetByID())
		// app.Post("/test-scenarios/search", testScenarioHandler.GetPaginated())
		// app.Post("/test-scenarios/update", testScenarioHandler.Update())

		// test scenario operationn-up
		// app.Post("/test-scenarios/:id/start", testScenarioOperationHandler.Start())
		// app.Post("/test-scenarios/:id/pause", testScenarioOperationHandler.Pause())
		// app.Post("/test-scenarios/:id/resume", testScenarioOperationHandler.Resume())
		// app.Post("/test-scenarios/:id/stop", testScenarioOperationHandler.Stop())
		// app.Post("/test-scenarios/:id/delete", testScenarioOperationHandler.Delete())

		// database-metadata
		app.Get("/databases", databaseMetadataHandler.GetAll)
		app.Post("/databases/tables", databaseMetadataHandler.GetTablesByDBNamePost)

		// swagger endpoint
		// app.Get("/docs/*", swaggerHandler)

	})

	zap.L().Info("server starting",
		zap.Int("port", cfg.Server.Port),
	)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		// if err := http.ListenAndServe(addr, app); err != nil {
		// 	errChan <- fmt.Errorf("fiber listen on %s failed: %w", addr, err)
		// }

		if err := srv.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("server listen on %s failed: %w", addr, err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		zap.L().Info("shutting down server gracefully")
		// Gracefully shutdown the server
		if err := srv.Shutdown(ctx); err != nil {
			zap.L().Error("server shutdown failed", zap.Error(err))

			return fmt.Errorf("server shutdown failed: %w", err)
		}
		zap.L().Info("server shut down successfully")

		return nil
	case err := <-errChan:
		zap.L().Error("server error", zap.Error(err))

		return err
	}
}
