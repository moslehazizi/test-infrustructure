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
	ServiceName string `envconfig:"SERVICE_NAME"`
	Server      Server
	Kafka       Kafka
	Postgres    Postgres
	Kubernetese Kubernetese
	Otlp        Otlp
	Logger      Logger
}

type Otlp struct {
	GRPCPort int    `envconfig:"OTLP_GRPC_PORT"`
	GRPCHost string `envconfig:"OTLP_GRPC_HOST"`
}

type Server struct {
	Port                        int           `envconfig:"HTTP_PORT" default:"8080"`
	SwaggerHost                 string        `envconfig:"SWAGGER_HOST" default:"localhost"`
	SwaggerScheme               StringSlice   `envconfig:"SWAGGER_SCHEME" default:"http,https"`
	SwaggerDocJSON              string        `envconfig:"SWAGGER_DOC_JSON" default:"doc.json"`
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
	Password          string        `envconfig:"KAFKA_PASSWORD" json:"-"`
	DialerTimeout     time.Duration `envconfig:"KAFKA_DIALER_TIMEOUT" default:"10s"`
	MaxBytes          int           `envconfig:"KAFKA_MAX_BYTES" default:"10485760"` // 10MB = 10 * 1024 * 1024
	ConsumerGroup     string        `envconfig:"KAFKA_CONSUMER_GROUP" default:"factorial-consumer-group"`
	BatchTimeout      time.Duration `envconfig:"KAFKA_BATCH_TIMEOUT" default:"5ms"`
	BatchSize         int           `envconfig:"KAFKA_BATCH_SIZE" default:"1000"`
	BatchBytes        int           `envconfig:"KAFKA_BATCH_BYTES" default:"1000000"` // 1MB = 1e6
	MaxAttempts       int           `envconfig:"KAFKA_MAX_ATTEMPTS" default:"10"`
	WriteTimeOut      time.Duration `envconfig:"KAFKA_WRITE_TIME_OUT" default:"1ms"`
	AttemptsSleepTime time.Duration `envconfig:"KAFKA_ATTEMPTS_SLEEP_TIME" default:"100ms"`
}

type Postgres struct {
	Host               string        `envconfig:"POSTGRES_HOST"`
	Port               int           `envconfig:"POSTGRES_PORT"`
	User               string        `envconfig:"POSTGRES_USER"`
	Password           string        `envconfig:"POSTGRES_PASSWORD" json:"-"`
	Database           string        `envconfig:"POSTGRES_DATABASE"`
	SSLMode            string        `envconfig:"POSTGRES_SSL_MODE"`
	MaxOpenConnections int           `envconfig:"POSTGRES_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int           `envconfig:"POSTGRES_MAX_IDLE_CONNECTIONS"`
	ConnMaxLifetime    time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME"`
	ConnMaxIdleTime    time.Duration `envconfig:"POSTGRES_CONN_MAX_IDLE_TIME"`
}

type Logger struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"` // json or console
	Output string `envconfig:"LOG_OUTPUT" default:"stdout"`
}

type Kubernetese struct {
	NameSpace                       string        `envconfig:"KUBERNETES_NAMESPACE" default:"default"`
	ContainerRegistryUrl            string        `envconfig:"CONTAINER_REGISTRY_URL" default:"chalenge.azurecr.io/"`
	ImagePullPolicy                 string        `envconfig:"IMAGE_PULL_POLICY" default:"Always"`
	IngressClassName                string        `envconfig:"INGRESS_CLASS_NAME" default:"traefik"`
	IngressHost                     string        `envconfig:"INGRESS_HOST" default:"127.0.0.1"`
	IngressPort                     int           `envconfig:"INGRESS_PORT" default:"8081"`
	MotherServiceImage              string        `envconfig:"MOTHER_SERVICE_IMAGE" default:"challenge-mother-service:0.1"`
	MotherServiceAPPServe           string        `envconfig:"MOTHER_SERVICE_APP_SERVE" default:"mother-service-serve"`
	MotherServiceAPPServeWaitReady  time.Duration `envconfig:"MOTHER_SERVICE_APP_SERVE_WAIT_READY" default:"10s"`
	MotherServiceAPPJobs            string        `envconfig:"MOTHER_SERVICE_APP_JOBS" default:"mother-service-jobs"`
	MotherServiceAPPJobsWaitReady   time.Duration `envconfig:"MOTHER_SERVICE_APP_JOBS_WAIT_READY" default:"10s"`
	MotherServiceKafkaDatabaseTopic string        `envconfig:"MOTHER_SERVICE_KAFKA_DATABASE_TOPIC" default:"mother-db"`
	MotherServiceLiveFeedTopic      string        `envconfig:"MOTHER_SERVICE_KAFKA_LIVE_FEED_TOPIC" default:"mother-live-feed"`
	MotherServiceKafkaHost          string        `envconfig:"MOTHER_SERVICE_KAFKA_HOST" default:"kafka"`
	MotherServiceKafkaPort          int           `envconfig:"MOTHER_SERVICE_KAFKA_PORT" default:"9092"`
	MotherServicePostgresHost       string        `envconfig:"MOTHER_SERVICE_POSTGRES_HOST" default:"postgres"`
	TestServiceImage                string        `envconfig:"TEST_SERVICE_IMAGE" default:"challenge-test-service:0.1"`
	TestServiceAPPServe             string        `envconfig:"TEST_SERVICE_APP_SERVE" default:"test-service-serve"`
	TestServiceAPPServeWaitReady    time.Duration `envconfig:"TEST_SERVICE_APP_SERVE_WAIT_READY" default:"10s"`
	TestServiceAPPJobs              string        `envconfig:"TEST_SERVICE_APP_JOBS" default:"test-service-jobs"`
	TestServiceAPPJobsWaitReady     time.Duration `envconfig:"TEST_SERVICE_APP_JOBS_WAIT_READY" default:"10s"`
	TestServicePostgresHost         string        `envconfig:"TEST_SERVICE_POSTGRES_HOST" default:"localhost"`
	TestServiceKafkaHost            string        `envconfig:"TEST_SERVICE_KAFKA_HOST" default:"kafka"`
	TestServiceKafkaPort            int           `envconfig:"TEST_SERVICE_KAFKA_PORT" default:"9092"`
	TestServiceKafkaDatabaseTopic   string        `envconfig:"TEST_SERVICE_KAFKA_DATABASE_TOPIC" default:"test-db"`
	TestServiceLiveFeedTopic        string        `envconfig:"TEST_SERVICE_KAFKA_LIVE_FEED_TOPIC" default:"test-live-feed"`
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
