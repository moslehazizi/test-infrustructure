package server

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository/postgres"
	"control-panel-service/internal/server/handler"
	"control-panel-service/internal/usecase"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Serve(ctx context.Context, cfg *config.Config) error {
	app := fiber.New(fiber.Config{
		BodyLimit:    cfg.Server.PostBodyLimit,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	// Global middlewares
	app.Use(cors.New())

	// 🔒 Rate Limiter (GLOBAL)
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

	eventProducer, err := provider.NewKafkaEventProducer(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create kafka event producer: %w", err)
	}
	defer func() {
		err := eventProducer.Close()
		if err != nil {
			log.Printf("failed to close event producer: %v", err)
		}
	}()

	db, err := postgres.OpenConnection(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("could not connect to postgres: %w", err)
	}

	motherService := usecase.NewMotherService(postgres.NewMotherServiceRepository(db))
	handler := handler.NewMotherServiceHandler(motherService)

	apiV1 := app.Group("/api/v1")

	// Register APIs
	apiV1.Options("/mother-service", handler.Create())
	apiV1.Options("/mother-service/:id", handler.GetByID())
	apiV1.Options("/mother-service/paginated", handler.GetPaginated())

	log.Printf("🚀 Fiber server started on :%d\n", cfg.Server.Port)

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
		log.Println("Shutting down server gracefully...")
		// Gracefully shutdown the server
		if err := app.Shutdown(); err != nil {
			return fmt.Errorf("fiber shutdown failed: %w", err)
		}
		log.Println("Server shut down successfully")

		return nil
	case err := <-errChan:
		return err
	}
}
