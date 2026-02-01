/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/server"
	"control-panel-service/pkg/telemetry"
	"fmt"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
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
			log.Fatal("could not load config:", err)

			return
		}

		// Create a context that can be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		conn, err := initConn()
		if err != nil {
			log.Fatal(err)
		}

		serviceName := semconv.ServiceNameKey.String(cfg.ServiceName)

		res, err := resource.New(ctx,
			resource.WithAttributes(
				serviceName,
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		shutdownTracerProvider, err := telemetry.InitTracerProvider(ctx, res, conn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := shutdownTracerProvider(ctx); err != nil {
				log.Fatalf("failed to shutdown TracerProvider: %s", err)
			}
		}()

		shutdownMeterProvider, err := telemetry.InitMeterProvider(ctx, res, conn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := shutdownMeterProvider(ctx); err != nil {
				log.Fatalf("failed to shutdown MeterProvider: %s", err)
			}
		}()

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Run the server in a goroutine
		serverDone := make(chan struct{})
		go func() {
			log.Println("starting server serve")
			if err := server.Serve(ctx, &cfg); err != nil {
				log.Println("could not serve:", err)
			}
			close(serverDone)
		}()

		// Wait for shutdown signal or server completion
		select {
		case sig := <-sigChan:
			log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)
			cancel() // Cancel the context to stop the server

			// Give the server some time to finish gracefully
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
			defer shutdownCancel()

			select {
			case <-serverDone:
				log.Println("Server completed gracefully")
			case <-shutdownCtx.Done():
				log.Println("Shutdown timeout reached, forcing exit")
			}
		case <-serverDone:
			log.Println("Server completed")
		}
	},
}

func initConn() (*grpc.ClientConn, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("could not load config:", err)

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
