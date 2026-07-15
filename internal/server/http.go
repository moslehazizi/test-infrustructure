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
	"encoding/json"
	"fmt"
	"net/http"

	pslq "control-panel-service/pkg/database/postgres"
	kuber "control-panel-service/pkg/kubernetes"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/httprate"

	globalmiddleware "github.com/go-chi/chi/v5/middleware"
)

func Serve(ctx context.Context, cfg *config.Config) error {
	// app := fiber.New(fiber.Config{
	// 	BodyLimit:    cfg.Server.PostBodyLimit,
	// 	ReadTimeout:  cfg.Server.ReadTimeout,
	// 	WriteTimeout: cfg.Server.WriteTimeout,
	// })

	app := chi.NewRouter()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      app,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Global middlewares
	app.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// nolint
	app.Use(middleware.MetricsMiddleware)

	app.Use(globalmiddleware.ClientIPFromRemoteAddr)

	// Logging middleware should be early to capture all requests
	app.Use(middleware.LoggingMiddleware())

	app.Use(httprate.LimitBy(
		cfg.Server.RateLimitMaxRequest,
		cfg.Server.RateLimitExpirationDuration,
		func(r *http.Request) (string, error) {
			return httprate.CanonicalizeIP(globalmiddleware.GetClientIP(r.Context())), nil
		},
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too many requests, please try again later",
			})
		}),
	))

	docs.SwaggerInfo.Host = cfg.Server.SwaggerHost
	docs.SwaggerInfo.Schemes = cfg.Server.SwaggerScheme

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

	testScenarioUsecase := usecase.NewTestScenarioUsecase(
		testScenarioRepository,
		postgres.NewTestCategoryRepository(db),
		postgres.NewTestServiceConfigRepository(db),
		motherServiceRepository,
		postgres.NewTestServiceRepository(db),
	)
	testScenarioOperationUsecase := usecase.NewTestScenarioOperationUsecase(
		testScenarioRepository,
		stressTestExecutionManager,
	)
	motherHandler := handler.NewMotherServiceHandler(motherService)
	testCategoryHandler := handler.NewTestCategoryHandler(cfg, postgres.NewTestCategoryRepository(db))

	testScenarioHandler := handler.NewTestScenarioHandler(testScenarioUsecase)
	testScenarioOperationHandler := handler.NewTestScenarioOperationHandler(testScenarioOperationUsecase)
	databaseMetadataService := usecase.NewDatabaseMetadata(postgres.NewDatabaseMetadataRepository(db, cfg))
	databaseMetadataHandler := handler.NewDatabaseMetadataHandler(databaseMetadataService)

	app.Route("/api/v1", func(app chi.Router) {
		// Mother service
		app.Post("/mother-services", motherHandler.Create)
		app.Get("/mother-services/{id}", motherHandler.GetByID)
		app.Post("/mother-services/search", motherHandler.GetPaginated)
		app.Post("/mother-services/{id}/delete", motherHandler.Delete)

		// test category
		app.Get("/test-categories", testCategoryHandler.GetAll)
		app.Get("/test-categories/{id}", testCategoryHandler.GetByID)

		// test scenario
		app.Post("/test-scenarios", testScenarioHandler.Create)
		app.Get("/test-scenarios/{id}", testScenarioHandler.GetByID)
		app.Post("/test-scenarios/search", testScenarioHandler.GetPaginated)
		app.Post("/test-scenarios/update", testScenarioHandler.Update)

		// test scenario operationn-up
		app.Post("/test-scenarios/{id}/start", testScenarioOperationHandler.Start)
		app.Post("/test-scenarios/{id}/pause", testScenarioOperationHandler.Pause)
		app.Post("/test-scenarios/{id}/resume", testScenarioOperationHandler.Resume)
		app.Post("/test-scenarios/{id}/stop", testScenarioOperationHandler.Stop)
		app.Post("/test-scenarios/{id}/delete", testScenarioOperationHandler.Delete)

		// database-metadata
		app.Get("/databases", databaseMetadataHandler.GetAll)
		app.Post("/databases/tables", databaseMetadataHandler.GetTablesByDBNamePost)

		// swagger endpoint
		app.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL(cfg.Server.SwaggerDocJSON),
			httpSwagger.DeepLinking(true),
			httpSwagger.PersistAuthorization(true),
			httpSwagger.DocExpansion("list"),
		))

	})

	zap.L().Info("server starting",
		zap.Int("port", cfg.Server.Port),
	)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
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
