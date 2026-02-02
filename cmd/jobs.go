/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/jobs"
	"control-panel-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var jobsCmd = &cobra.Command{
	Use:   "jobs",
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
			_, _ = os.Stderr.WriteString("could not load config: " + err.Error() + "\n")
			os.Exit(1)
		}

		// Initialize logger early
		loggerConfig := &logger.Config{
			Level:  cfg.Logger.Level,
			Format: cfg.Logger.Format,
			Output: cfg.Logger.Output,
		}
		zapLogger, err := logger.New(loggerConfig, logger.DefaultJobServiceName)
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

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Run the consumer job in a goroutine
		jobDone := make(chan struct{})
		go func() {
			zap.L().Info("starting job")
			if err := jobs.Serve(ctx, &cfg); err != nil {
				zap.L().Error("job error", zap.Error(err))
			}
			close(jobDone)
		}()

		// Wait for shutdown signal or job completion
		select {
		case sig := <-sigChan:
			zap.L().Info("received shutdown signal, initiating graceful shutdown",
				zap.String("signal", sig.String()),
			)
			cancel() // Cancel the context to stop the job

			// Give the job some time to finish gracefully
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
			defer shutdownCancel()

			select {
			case <-jobDone:
				zap.L().Info("job completed gracefully")
			case <-shutdownCtx.Done():
				zap.L().Warn("shutdown timeout reached, forcing exit")
			}
		case <-jobDone:
			zap.L().Info("job completed")
		}
	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jobsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
