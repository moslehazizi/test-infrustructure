package config

import (
	"control-panel-service/pkg"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("success fetch otlp config", func(t *testing.T) {
		expectedPort := 8080
		exportedHost := "localhost"
		expectedServiceName := "control-panel"

		os.Setenv("OTLP_GRPC_PORT", strconv.Itoa(expectedPort))
		os.Setenv("OTLP_GRPC_HOST", exportedHost)
		os.Setenv("SERVICE_NAME", expectedServiceName)

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, cfg.Otlp.GRPCPort, expectedPort)
		assert.Equal(t, cfg.Otlp.GRPCHost, exportedHost)
		assert.Equal(t, cfg.ServiceName, expectedServiceName)
	})
	t.Run("success fetch http config", func(t *testing.T) {
		expectedPort := 8080
		expectedHost := "localhost"
		expectedSwaggerScheme := []string{"http", "https"}
		expectedSwaggerDocJSON := "doc.json"
		expectedPostBodyLimit := 4 * 1024
		expectedReadTimeout := time.Second * 10
		expectedWriteTimeout := time.Second * 20
		expectedRateLimitMaxRequest := int(50_000)
		expectedRateLimitExpirationDuration := time.Minute
		expectedShutdownTimeout := 30 * time.Second

		os.Setenv("HTTP_PORT", strconv.Itoa(expectedPort))
		os.Setenv("HTTP_POST_BODY_LIMIT", strconv.Itoa(expectedPostBodyLimit))
		os.Setenv("HTTP_READ_TIMEOUT", fmt.Sprintf("%v", expectedReadTimeout))
		os.Setenv("HTTP_WRITE_TIMEOUT", fmt.Sprintf("%v", expectedWriteTimeout))
		os.Setenv("HTTP_RATE_LIMIT_MAX_REQUEST", strconv.Itoa(expectedRateLimitMaxRequest))
		os.Setenv("HTTP_RATE_LIMIT_EXPIRATION_DURATION", fmt.Sprintf("%v", expectedRateLimitExpirationDuration))
		os.Setenv("HTTP_SHUTDOWN_TIMEOUT", fmt.Sprintf("%v", expectedShutdownTimeout))
		os.Setenv("SWAGGER_HOST", expectedHost)
		os.Setenv("SWAGGER_SCHEME", "http,https")
		os.Setenv("SWAGGER_DOC_JSON", "doc.json")

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, cfg.Server.Port, expectedPort)
		assert.Equal(t, cfg.Server.PostBodyLimit, expectedPostBodyLimit)
		assert.Equal(t, cfg.Server.ReadTimeout, expectedReadTimeout)
		assert.Equal(t, cfg.Server.WriteTimeout, expectedWriteTimeout)
		assert.Equal(t, cfg.Server.RateLimitMaxRequest, expectedRateLimitMaxRequest)
		assert.Equal(t, cfg.Server.SwaggerHost, expectedHost)
		assert.Equal(t, []string(cfg.Server.SwaggerScheme), expectedSwaggerScheme)
		assert.Equal(t, cfg.Server.SwaggerDocJSON, expectedSwaggerDocJSON)
		assert.Equal(t, cfg.Server.RateLimitExpirationDuration, expectedRateLimitExpirationDuration)
		assert.Equal(t, cfg.Server.ShutdownTimeout, expectedShutdownTimeout)
	})
	t.Run("success fetch kafka config", func(t *testing.T) {
		expectedKafkaHost := "127.0.0.1"
		expectedKafkaPort := 9092
		expectedKafkaUsername := "kafka_user"
		expectedKafkaPassword := "kafka_password_123"
		expectedKafkaProvisioningTopic := "factorial"
		expectedKafkaDialerTimeout := 15 * time.Second
		expectedKafkaMaxBytes := 20 * 1024 * 1024 // 20MB
		expectedKafkaConsumerGroup := "factorial-consumer-group"
		expectedKafkaBatchTimeout := 10 * time.Millisecond
		expectedKafkaBatchSize := 2000
		expectedKafkaBatchBytes := 2000000 // 2MB

		os.Setenv("KAFKA_HOST", expectedKafkaHost)
		os.Setenv("KAFKA_PORT", strconv.Itoa(expectedKafkaPort))
		os.Setenv("KAFKA_USERNAME", expectedKafkaUsername)
		os.Setenv("KAFKA_PASSWORD", expectedKafkaPassword)
		os.Setenv("KAFKA_PROVISIONING_TOPIC", expectedKafkaProvisioningTopic)
		os.Setenv("KAFKA_DIALER_TIMEOUT", fmt.Sprintf("%v", expectedKafkaDialerTimeout))
		os.Setenv("KAFKA_MAX_BYTES", strconv.Itoa(expectedKafkaMaxBytes))
		os.Setenv("KAFKA_CONSUMER_GROUP", expectedKafkaConsumerGroup)
		os.Setenv("KAFKA_BATCH_TIMEOUT", fmt.Sprintf("%v", expectedKafkaBatchTimeout))
		os.Setenv("KAFKA_BATCH_SIZE", strconv.Itoa(expectedKafkaBatchSize))
		os.Setenv("KAFKA_BATCH_BYTES", strconv.Itoa(expectedKafkaBatchBytes))

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, cfg.Kafka.Host, expectedKafkaHost)
		assert.Equal(t, cfg.Kafka.Port, expectedKafkaPort)
		assert.Equal(t, cfg.Kafka.Username, expectedKafkaUsername)
		assert.Equal(t, cfg.Kafka.Password, expectedKafkaPassword)
		assert.Equal(t, cfg.Kafka.ProvisioningTopic, expectedKafkaProvisioningTopic)
		assert.Equal(t, cfg.Kafka.DialerTimeout, expectedKafkaDialerTimeout)
		assert.Equal(t, cfg.Kafka.MaxBytes, expectedKafkaMaxBytes)
		assert.Equal(t, cfg.Kafka.ConsumerGroup, expectedKafkaConsumerGroup)
		assert.Equal(t, cfg.Kafka.BatchTimeout, expectedKafkaBatchTimeout)
		assert.Equal(t, cfg.Kafka.BatchSize, expectedKafkaBatchSize)
		assert.Equal(t, cfg.Kafka.BatchBytes, expectedKafkaBatchBytes)

	})
	t.Run("success fetch kafka config without username and password", func(t *testing.T) {
		expectedKafkaHost := "127.0.0.1"
		expectedKafkaPort := 9092
		expectedKafkaProvisioningTopic := "factorial"

		os.Setenv("KAFKA_HOST", expectedKafkaHost)
		os.Setenv("KAFKA_PORT", strconv.Itoa(expectedKafkaPort))
		os.Setenv("KAFKA_PROVISIONING_TOPIC", expectedKafkaProvisioningTopic)
		// Explicitly unset username and password to test backward compatibility
		os.Unsetenv("KAFKA_USERNAME")
		os.Unsetenv("KAFKA_PASSWORD")

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, cfg.Kafka.Host, expectedKafkaHost)
		assert.Equal(t, cfg.Kafka.Port, expectedKafkaPort)
		assert.Equal(t, cfg.Kafka.Username, "")
		assert.Equal(t, cfg.Kafka.Password, "")
		assert.Equal(t, cfg.Kafka.ProvisioningTopic, expectedKafkaProvisioningTopic)

	})
	t.Run("success fetch postgres config", func(t *testing.T) {
		expectedPostgresHost := "127.0.0.1"
		expectedPostgresPort := 5432
		expectedPostgresUser := "postgres"
		expectedPostgresPassword := "1234$*&^3249M"
		expectedPostgresDatabase := "control-panel"
		expectedPostgresSSLMode := "off"
		expectedPostgresMaxOpenConnection := 100
		expectedPostgresMaxIdleConnections := 10
		expectedPostgresConnMaxLifetime := time.Hour
		expectedPostgresConnMaxIdleTime := 30 * time.Minute

		os.Setenv("POSTGRES_HOST", expectedPostgresHost)
		os.Setenv("POSTGRES_PORT", strconv.Itoa(expectedPostgresPort))
		os.Setenv("POSTGRES_USER", expectedPostgresUser)
		os.Setenv("POSTGRES_PASSWORD", expectedPostgresPassword)
		os.Setenv("POSTGRES_DATABASE", expectedPostgresDatabase)
		os.Setenv("POSTGRES_SSL_MODE", expectedPostgresSSLMode)
		os.Setenv("POSTGRES_MAX_OPEN_CONNECTIONS", strconv.Itoa(expectedPostgresMaxOpenConnection))
		os.Setenv("POSTGRES_MAX_IDLE_CONNECTIONS", strconv.Itoa(expectedPostgresMaxIdleConnections))
		os.Setenv("POSTGRES_CONN_MAX_LIFETIME", fmt.Sprintf("%v", expectedPostgresConnMaxLifetime))
		os.Setenv("POSTGRES_CONN_MAX_IDLE_TIME", fmt.Sprintf("%v", expectedPostgresConnMaxIdleTime))

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, cfg.Postgres.Host, expectedPostgresHost)
		assert.Equal(t, cfg.Postgres.Port, expectedPostgresPort)
		assert.Equal(t, cfg.Postgres.User, expectedPostgresUser)
		assert.Equal(t, cfg.Postgres.Password, expectedPostgresPassword)
		assert.Equal(t, cfg.Postgres.Database, expectedPostgresDatabase)
		assert.Equal(t, cfg.Postgres.SSLMode, expectedPostgresSSLMode)
		assert.Equal(t, cfg.Postgres.MaxOpenConnections, expectedPostgresMaxOpenConnection)
		assert.Equal(t, cfg.Postgres.MaxIdleConnections, expectedPostgresMaxIdleConnections)
		assert.Equal(t, cfg.Postgres.ConnMaxLifetime, expectedPostgresConnMaxLifetime)
	})

	t.Run("success fetch logger config", func(t *testing.T) {
		expectedLogLevel := "debug"
		expectedLogFormat := "console"
		expectedLogOutput := "stderr"

		os.Setenv("LOG_LEVEL", expectedLogLevel)
		os.Setenv("LOG_FORMAT", expectedLogFormat)
		os.Setenv("LOG_OUTPUT", expectedLogOutput)

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, cfg.Logger.Level, expectedLogLevel)
		assert.Equal(t, cfg.Logger.Format, expectedLogFormat)
		assert.Equal(t, cfg.Logger.Output, expectedLogOutput)
	})

	t.Run("check default values config", func(t *testing.T) {
		expectedDefaultPort := 8080
		expectedDefaultHTTPHost := "localhost"
		expectedDefaultSwaggerDocJSON := "doc.json"
		expectedDefaultPostBodyLimit := 4 * 1024
		expectedDefaultReadTimeout := time.Second * 10
		expectedDefaultWriteTimeout := time.Second * 20
		expectedDefaultRateLimitMaxRequest := int(50_000)
		expectedDefaultRateLimitExpirationDuration := time.Minute
		expectedDefaultShutdownTimeout := 30 * time.Second
		expectedDefaultKafkaDialerTimeout := 10 * time.Second
		expectedDefaultKafkaMaxBytes := 10485760 // 10MB
		expectedDefaultKafkaBatchTimeout := 5 * time.Millisecond
		expectedDefaultKafkaBatchSize := 1000
		expectedDefaultKafkaBatchBytes := 1000000 // 1MB
		expectedDefaultKafkaProvisioningTopic := "provisioning"
		expectedDefaultLogLevel := "info"
		expectedDefaultLogFormat := "json"
		expectedDefaultLogOutput := "stdout"

		// Unset Kafka environment variables to test defaults
		os.Unsetenv("KAFKA_DIALER_TIMEOUT")
		os.Unsetenv("KAFKA_MAX_BYTES")
		os.Unsetenv("KAFKA_BATCH_TIMEOUT")
		os.Unsetenv("KAFKA_BATCH_SIZE")
		os.Unsetenv("KAFKA_BATCH_BYTES")
		// Unset shutdown timeout to test default
		os.Unsetenv("HTTP_SHUTDOWN_TIMEOUT")
		os.Unsetenv("KAFKA_PROVISIONING_TOPIC")
		// Unset logger environment variables to test defaults
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_FORMAT")
		os.Unsetenv("LOG_OUTPUT")

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, cfg.Server.Port, expectedDefaultPort)
		assert.Equal(t, cfg.Server.SwaggerHost, expectedDefaultHTTPHost)
		assert.Equal(t, []string(cfg.Server.SwaggerScheme), []string{"http", "https"})
		assert.Equal(t, cfg.Server.SwaggerDocJSON, expectedDefaultSwaggerDocJSON)
		assert.Equal(t, cfg.Server.PostBodyLimit, expectedDefaultPostBodyLimit)
		assert.Equal(t, cfg.Server.ReadTimeout, expectedDefaultReadTimeout)
		assert.Equal(t, cfg.Server.WriteTimeout, expectedDefaultWriteTimeout)
		assert.Equal(t, cfg.Server.RateLimitMaxRequest, expectedDefaultRateLimitMaxRequest)
		assert.Equal(t, cfg.Server.RateLimitExpirationDuration, expectedDefaultRateLimitExpirationDuration)
		assert.Equal(t, cfg.Server.ShutdownTimeout, expectedDefaultShutdownTimeout)
		assert.Equal(t, cfg.Kafka.DialerTimeout, expectedDefaultKafkaDialerTimeout)
		assert.Equal(t, cfg.Kafka.MaxBytes, expectedDefaultKafkaMaxBytes)
		assert.Equal(t, cfg.Kafka.BatchTimeout, expectedDefaultKafkaBatchTimeout)
		assert.Equal(t, cfg.Kafka.BatchSize, expectedDefaultKafkaBatchSize)
		assert.Equal(t, cfg.Kafka.BatchBytes, expectedDefaultKafkaBatchBytes)
		assert.Equal(t, cfg.Kafka.ProvisioningTopic, expectedDefaultKafkaProvisioningTopic)
		assert.Equal(t, cfg.Logger.Level, expectedDefaultLogLevel)
		assert.Equal(t, cfg.Logger.Format, expectedDefaultLogFormat)
		assert.Equal(t, cfg.Logger.Output, expectedDefaultLogOutput)
	})

	t.Run("error - invalid environment variable value", func(t *testing.T) {
		// Set an invalid value for POSTGRES_PORT that can't be parsed as int
		os.Setenv("POSTGRES_PORT", "invalid_port")
		defer os.Unsetenv("POSTGRES_PORT")

		cfg, err := LoadConfig()

		assert.Error(t, err)
		assert.True(t, errors.Is(err, pkg.ErrFailedToLoadConfig), "error should wrap ErrFailedToLoadConfig")
		assert.Contains(t, err.Error(), "failed to load config from env")
		assert.Equal(t, Config{}, cfg)
	})
}
