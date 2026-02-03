/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/server"
	"control-panel-service/pkg/logger"
	"control-panel-service/pkg/telemetry"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load(envFile)

		cfg, err := config.LoadConfig()
		if err != nil {
			// Use standard log for config loading errors since logger isn't initialized yet
			_, _ = os.Stderr.WriteString("could not load config: " + err.Error() + "\n")
			os.Exit(1)
		}

		// Initialize logger early
		loggerConfig := &logger.Config{
			Level:  cfg.Logger.Level,
			Format: cfg.Logger.Format,
			Output: cfg.Logger.Output,
		}
		zapLogger, err := logger.New(loggerConfig, logger.DefaultServiceName)
		if err != nil {
			_, _ = os.Stderr.WriteString("could not initialize logger: " + err.Error() + "\n")
			os.Exit(1)
		}
		// Replace global logger so zap.L() can be used throughout the application
		zap.ReplaceGlobals(zapLogger)
		defer func() {
			if syncErr := zap.L().Sync(); syncErr != nil {
				// Ignore sync errors on stderr/stdout in some cases
				_ = syncErr
			}
		}()

		zap.L().Info("application starting",
			zap.String("log_level", cfg.Logger.Level),
			zap.String("log_format", cfg.Logger.Format),
		)

		// Create a context that can be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		conn, err := initConn()
		if err != nil {
			_, _ = os.Stderr.WriteString("could not initialize grpc conn to otel collector: " + err.Error() + "\n")
			os.Exit(1)
		}

		serviceName := semconv.ServiceNameKey.String(cfg.ServiceName)

		res, err := resource.New(ctx,
			resource.WithAttributes(
				serviceName,
			),
		)
		if err != nil {
			_, _ = os.Stderr.WriteString("could not create new resource to otel collector: " + err.Error() + "\n")
			os.Exit(1)
		}

		shutdownTracerProvider, err := telemetry.InitTracerProvider(ctx, res, conn)
		if err != nil {
			_, _ = os.Stderr.WriteString("could not init tracer provider: " + err.Error() + "\n")
			os.Exit(1)
		}
		defer func() {
			if err := shutdownTracerProvider(ctx); err != nil {
				_, _ = os.Stderr.WriteString("failed to shutdown TracerProvider: " + err.Error() + "\n")
				os.Exit(1)
			}
		}()

		shutdownMeterProvider, err := telemetry.InitMeterProvider(ctx, res, conn)
		if err != nil {
			_, _ = os.Stderr.WriteString("could not init metric provider: " + err.Error() + "\n")
			os.Exit(1)
		}
		defer func() {
			if err := shutdownMeterProvider(ctx); err != nil {
				_, _ = os.Stderr.WriteString("failed to shutdown MeterProvider: " + err.Error() + "\n")
				os.Exit(1)
			}
		}()

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Run the server in a goroutine
		serverDone := make(chan struct{})
		go func() {
			zap.L().Info("starting server")
			if err := server.Serve(ctx, &cfg); err != nil {
				zap.L().Error("server error", zap.Error(err))
			}
			close(serverDone)
		}()

		// Wait for shutdown signal or server completion
		select {
		case sig := <-sigChan:
			zap.L().Info("received shutdown signal, initiating graceful shutdown",
				zap.String("signal", sig.String()),
			)
			cancel() // Cancel the context to stop the server

			// Give the server some time to finish gracefully
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
			defer shutdownCancel()

			select {
			case <-serverDone:
				zap.L().Info("server completed gracefully")
			case <-shutdownCtx.Done():
				zap.L().Warn("shutdown timeout reached, forcing exit")
			}
		case <-serverDone:
			zap.L().Info("server completed")
		}
	},
}

func initConn() (*grpc.ClientConn, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Error("could not load config:", zap.Error(err))

		return nil, nil
	}

	url := fmt.Sprintf("%s:%d", cfg.Otlp.GRPCHost, cfg.Otlp.GRPCPort)

	conn, err := grpc.NewClient(url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var serviceName = semconv.ServiceNameKey.String("test-service")
