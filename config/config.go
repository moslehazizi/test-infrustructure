package config

import (
	"control-panel-service/pkg"
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type StringSlice []string

func (ss *StringSlice) Set(value string) error {
	if value == "" {
		*ss = []string{"http", "https"}
	} else {
		parts := strings.Split(value, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		*ss = parts
	}
	return nil
}

type Config struct {
	Server   Server
	Kafka    Kafka
	Postgres Postgres
}

type Server struct {
	Port                        int           `envconfig:"HTTP_PORT" default:"8080"`
	Host                        string        `envconfig:"HTTP_HOST" default:"localhost"`
	SwaggerScheme               StringSlice   `envconfig:"SWAGGER_SCHEME" default:"http,https"`
	PostBodyLimit               int           `envconfig:"HTTP_POST_BODY_LIMIT" default:"4096"` // 4096 = 4KB
	ReadTimeout                 time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout                time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"20s"`
	RateLimitMaxRequest         int           `envconfig:"HTTP_RATE_LIMIT_MAX_REQUEST" default:"50000"`
	RateLimitExpirationDuration time.Duration `envconfig:"HTTP_RATE_LIMIT_EXPIRATION_DURATION" default:"1m"`
	ShutdownTimeout             time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"30s"`
}

type Kafka struct {
	Host              string        `envconfig:"KAFKA_HOST"`
	Port              int           `envconfig:"KAFKA_PORT"`
	Username          string        `envconfig:"KAFKA_USERNAME"`
	Password          string        `envconfig:"KAFKA_PASSWORD"`
	ProvisioningTopic string        `envconfig:"KAFKA_PROVISIONING_TOPIC" default:"provisioning"`
	DialerTimeout     time.Duration `envconfig:"KAFKA_DIALER_TIMEOUT" default:"10s"`
	MaxBytes          int           `envconfig:"KAFKA_MAX_BYTES" default:"10485760"` // 10MB = 10 * 1024 * 1024
	ConsumerGroup     string        `envconfig:"KAFKA_CONSUMER_GROUP" default:"factorial-consumer-group"`
	BatchTimeout      time.Duration `envconfig:"KAFKA_BATCH_TIMEOUT" default:"5ms"`
	BatchSize         int           `envconfig:"KAFKA_BATCH_SIZE" default:"1000"`
	BatchBytes        int           `envconfig:"KAFKA_BATCH_BYTES" default:"1000000"` // 1MB = 1e6
}

type Postgres struct {
	Host               string        `envconfig:"POSTGRES_HOST"`
	Port               int           `envconfig:"POSTGRES_PORT"`
	User               string        `envconfig:"POSTGRES_USER"`
	Password           string        `envconfig:"POSTGRES_PASSWORD"`
	Database           string        `envconfig:"POSTGRES_DATABASE"`
	SSLMode            string        `envconfig:"POSTGRES_SSL_MODE"`
	MaxOpenConnections int           `envconfig:"POSTGRES_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int           `envconfig:"POSTGRES_MAX_IDLE_CONNECTIONS"`
	ConnMaxLifetime    time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME"`
	ConnMaxIdleTime    time.Duration `envconfig:"POSTGRES_CONN_MAX_IDLE_TIME"`
}

// var GlobalConfigInstance *Config

func LoadConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)

	if err != nil {
		return Config{}, fmt.Errorf("%w: %w", pkg.ErrFailedToLoadConfig, err)
	}

	return cfg, nil
}
