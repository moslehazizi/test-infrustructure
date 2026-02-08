package logger

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// ErrInvalidLogLevel is returned when an invalid log level is provided.
	ErrInvalidLogLevel = errors.New("invalid log level")
)

const (
	// DefaultServiceName is the default service name added to all logs.
	DefaultServiceName    = "challenge-control-panel-service"
	DefaultJobServiceName = "challenge-control-panel-job-service"
)

// Config holds logger configuration.
type Config struct {
	Level  string // debug, info, warn, error, fatal
	Format string // json, console
	Output string // stdout, stderr, or file path
}

// New creates a new zap logger based on the provided configuration.
// If config is nil, it uses production defaults (JSON format, Info level).
func New(config *Config, serviceName string) (*zap.Logger, error) {
	if config == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create production logger: %w", err)
		}

		return logger, nil
	}

	var zapConfig zap.Config
	var encoderConfig zapcore.EncoderConfig

	// Set encoder config based on format
	if strings.ToLower(config.Format) == "console" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig = zap.NewProductionConfig()
		// Apply the encoder config to ensure ISO8601 timestamps
		zapConfig.EncoderConfig = encoderConfig
	}

	// Set log level
	level, err := parseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Set output
	switch {
	case config.Output != "" && config.Output != "stdout" && config.Output != "stderr":
		// File output
		zapConfig.OutputPaths = []string{config.Output}
		zapConfig.ErrorOutputPaths = []string{config.Output}
	case config.Output == "stderr":
		zapConfig.OutputPaths = []string{"stderr"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}
	default:
		// Default to stdout
		zapConfig.OutputPaths = []string{"stdout"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}
	}

	// Build logger
	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	// Add service field to all logs automatically
	logger := zapLogger.With(zap.String(FieldService, serviceName))

	return logger, nil
}

// parseLevel converts a string level to zapcore.Level.
func parseLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("%w: %s", ErrInvalidLogLevel, level)
	}
}
