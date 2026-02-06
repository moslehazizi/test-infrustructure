package usecase

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	kubernetese "control-panel-service/pkg/kubernetes"
	inEntity "control-panel-service/pkg/kubernetes/domain/entity"
	"control-panel-service/pkg/number"
	"errors"
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

const (
	App   = "app"
	Main  = "./main"
	Serve = "serve"
	Jobs  = "jobs"

	Replication = 1
	Config      = "-config"
	Secret      = "-secret"

	// config map
	ServiceId                       = "SERVICE_ID"
	HTTPPort                        = "HTTP_PORT"
	HTTPPostBodyLimit               = "HTTP_POST_BODY_LIMIT"
	HTTPReadTimeout                 = "HTTP_READ_TIMEOUT"
	HTTPWriteTimeout                = "HTTP_WRITE_TIMEOUT"
	HTTPRateLimitMaxRequest         = "HTTP_RATE_LIMIT_MAX_REQUEST"
	HTTPRateLimitExpirationduration = "HTTP_RATE_LIMIT_EXPIRATION_DURATION"
	HTTPShutdownTimeout             = "HTTP_SHUTDOWN_TIMEOUT"
	HTTPErrorInjectionRate          = "HTTP_ERROR_INJECTION_RATE"
	HTTPDelayInjectionRate          = "HTTP_DELAY_INJECTION_RATE"
	HTTPDelayInjectionDurationMin   = "HTTP_DELAY_INJECTION_DURATION_MIN"
	HTTPDelayInjectionDurationMax   = "HTTP_DELAY_INJECTION_DURATION_MAX"
	LogLevel                        = "LOG_LEVEL"
	LogFormat                       = "LOG_FORMAT"
	LogOutput                       = "LOG_OUTPUT"
	KafkaHost                       = "KAFKA_HOST"
	KafkaPort                       = "KAFKA_PORT"
	KafkaDialerTimeout              = "KAFKA_DIALER_TIMEOUT"
	KafkaMaxBytes                   = "KAFKA_MAX_BYTES"
	KafkaBatchTimeout               = "KAFKA_BATCH_TIMEOUT"
	KafkaBatchSize                  = "KAFKA_BATCH_SIZE"
	KafkaBatchBytes                 = "KAFKA_BATCH_BYTES"
	KafkaDatabaseTopic              = "KAFKA_DATABASE_TOPIC"
	KafkaConsumerGroup              = "KAFKA_CONSUMER_GROUP"
	KafkaLiveFeedTopic              = "KAFKA_LIVE_FEED_TOPIC"
	PostgresHost                    = "POSTGRES_HOST"
	PostgresPort                    = "POSTGRES_PORT"
	PostgresDatabase                = "POSTGRES_DATABASE"
	PostgresTable                   = "POSTGRES_TABLE"
	PostgresSSLMode                 = "POSTGRES_SSL_MODE"
	PostgresMaxOpenConnection       = "POSTGRES_MAX_OPEN_CONNECTION"
	PostgresMaxIdleConnection       = "POSTGRES_MAX_IDLE_CONNECTIONS"
	PostgresConnMaxLifetime         = "POSTGRES_CONN_MAX_LIFETIME"
	PostgresConnMaxIdleTime         = "POSTGRES_CONN_MAX_IDLE_TIME"

	// secret map
	KafkaUsername    = "KAFKA_USERNAME"
	KafkaPassword    = "KAFKA_PASSWORD"
	PostgresUser     = "POSTGRES_USER"
	PostgresPassword = "POSTGRES_PASSWORD"
)

type MotherService interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error)
	DeployMotherService(ctx context.Context, motherService *entity.MotherService) error
}

func NewMotherService(
	cfg *config.Config,
	db database.Database,
	motherServiceRepo repository.MotherServiceRepository,
	eventProducer provider.EventProducer,
	kubernetes kubernetese.Kubernetese,
) MotherService {
	return &motherService{
		cfg,
		db,
		motherServiceRepo,
		eventProducer,
		kubernetes,
	}
}

type motherService struct {
	cfg               *config.Config
	db                database.Database
	motherServiceRepo repository.MotherServiceRepository
	eventProducer     provider.EventProducer
	kubernetes        kubernetese.Kubernetese
}

func (service *motherService) Create(ctx context.Context, motherService *entity.MotherService) (e error) {
	err := motherService.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}

	motherService.Status = entity.MotherServiceStatusPending

	tx := service.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
		}
	}()

	err = service.motherServiceRepo.Create(dbCtx, motherService)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceAlreadyExist) {
			return pkg.ErrMotherServiceAlreadyExist
		}

		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateMotherService, err)
	}
	_ = tx.Commit()

	err = service.DeployMotherService(ctx, motherService)
	if err != nil {
		return fmt.Errorf("%w, %w", pkg.ErrFailedToDeployMotherService, err)
	}

	// TODO: update status to deployed

	return nil
}

func (service *motherService) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	result, err := service.motherServiceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceNotFound) {
			return nil, pkg.ErrMotherServiceNotFound
		}

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetMotherService, err)
	}

	return result, nil
}

func (service *motherService) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error) {
	result, err := service.motherServiceRepo.GetPaginated(ctx, paginationRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherServices, err)
	}

	return result, nil
}

func (s *motherService) DeployMotherService(ctx context.Context, motherService *entity.MotherService) error {
	// Mother service serve deployment logic here
	// chekc mother service
	if motherService == nil {
		return fmt.Errorf("motherService struct is empty")
	}

	// Create config map
	randomResponseDelayMin, randomResponseDelayMax := number.CalcRange(motherService.ResponseDelayDuration, motherService.RandomResponseDelayMin, motherService.RandomResponseDelayMax)
	configMap := map[string]string{
		ServiceId:                       fmt.Sprintf("%v", motherService.ID),
		HTTPPort:                        strconv.Itoa(s.cfg.Server.Port),
		HTTPPostBodyLimit:               strconv.Itoa(s.cfg.Server.PostBodyLimit),
		HTTPReadTimeout:                 s.cfg.Server.ReadTimeout.String(),
		HTTPWriteTimeout:                s.cfg.Server.WriteTimeout.String(),
		HTTPRateLimitMaxRequest:         strconv.Itoa(s.cfg.Server.RateLimitMaxRequest),
		HTTPRateLimitExpirationduration: s.cfg.Server.RateLimitExpirationDuration.String(),
		HTTPShutdownTimeout:             s.cfg.Server.ShutdownTimeout.String(),
		HTTPErrorInjectionRate:          strconv.Itoa(motherService.ExceptionRate),
		HTTPDelayInjectionRate:          strconv.Itoa(motherService.ResponseDelayRate),
		HTTPDelayInjectionDurationMin:   strconv.Itoa(randomResponseDelayMin),
		HTTPDelayInjectionDurationMax:   strconv.Itoa(randomResponseDelayMax),
		LogLevel:                        s.cfg.Kubernetese.MotherServiceLogLevel,
		LogFormat:                       s.cfg.Kubernetese.MotherServiceLogFormat,
		LogOutput:                       s.cfg.Kubernetese.MotherServicelogOutput,
		KafkaPort:                       s.cfg.Kafka.Host,
		KafkaDialerTimeout:              s.cfg.Kafka.DialerTimeout.String(),
		KafkaMaxBytes:                   strconv.Itoa(s.cfg.Kafka.MaxBytes),
		KafkaBatchTimeout:               s.cfg.Kafka.BatchTimeout.String(),
		KafkaBatchSize:                  strconv.Itoa(s.cfg.Kafka.BatchSize),
		KafkaBatchBytes:                 strconv.Itoa(s.cfg.Kafka.BatchBytes),
		KafkaDatabaseTopic:              fmt.Sprintf("%s-%v", s.cfg.Kubernetese.MotherServiceKafkaDbTopic, motherService.ID),
		KafkaConsumerGroup:              fmt.Sprintf("%s-%v", s.cfg.Kubernetese.MotherServiceKafkaDbGroup, motherService.ID),
		KafkaLiveFeedTopic:              s.cfg.Kubernetese.MotherServiceLiveFeedTopic,
		PostgresHost:                    s.cfg.Postgres.Host,
		PostgresPort:                    strconv.Itoa(s.cfg.Postgres.Port),
		PostgresDatabase:                motherService.DatabaseName,
		PostgresTable:                   motherService.DatabaseTableName,
		PostgresSSLMode:                 s.cfg.Postgres.SSLMode,
		PostgresMaxOpenConnection:       strconv.Itoa(s.cfg.Postgres.MaxOpenConnections),
		PostgresMaxIdleConnection:       strconv.Itoa(s.cfg.Postgres.MaxIdleConnections),
		PostgresConnMaxLifetime:         s.cfg.Postgres.ConnMaxLifetime.String(),
		PostgresConnMaxIdleTime:         s.cfg.Postgres.ConnMaxIdleTime.String(),
	}

	// Create secret
	secretMap := map[string]string{
		KafkaUsername:    s.cfg.Kafka.Username,
		KafkaPassword:    s.cfg.Kafka.Password,
		PostgresUser:     s.cfg.Postgres.User,
		PostgresPassword: s.cfg.Postgres.Password,
	}

	// Serve :
	// Create deploy spec
	serveDepSpec := motherServDepSpec(s.cfg, Replication)

	// Create service spec
	serveSvcSpec := motherServeSvcSpec(s.cfg)

	// Call ApplyDeployment from kubernetese interface
	err := s.kubernetes.ApplyDeployment(ctx, serveDepSpec, configMap, secretMap)
	if err != nil {
		// TODO: log error
		return err
	}

	// Call ApplyService from kubernetese interface
	err = s.kubernetes.ApplyService(ctx, serveSvcSpec, configMap, secretMap)
	if err != nil {
		// TODO: log error
		return err
	}

	// Wait for deployment to be ready
	err = s.kubernetes.WaitForDeployment(ctx, s.cfg.Kubernetese.MotherServiceAPPServe, s.cfg.Kubernetese.MotherServiceAPPServeWaitReady)
	if err != nil {
		// TODO: log error
		return err
	}

	// jobs :
	// Create deploy spec
	jobsDepSpec := motherJobsDepSpec(s.cfg, Replication)

	// Call ApplyDeployment from kubernetese interface
	err = s.kubernetes.ApplyDeployment(ctx, jobsDepSpec, configMap, secretMap)
	if err != nil {
		// TODO: log error
		return err
	}

	// Wait for deployment to be ready
	err = s.kubernetes.WaitForDeployment(ctx, s.cfg.Kubernetese.MotherServiceAPPJobs, s.cfg.Kubernetese.MotherServiceAPPJobsWaitReady)
	if err != nil {
		// TODO: log error
		return err
	}

	return nil
}

func motherServDepSpec(cfg *config.Config, replica int) inEntity.DeploymentSpec {
	replicas := int32(replica)
	labels := map[string]string{App: cfg.Kubernetese.MotherServiceAPPServe}

	return inEntity.DeploymentSpec{
		Name: cfg.Kubernetese.MotherServiceAPPServe,
		Deployment: &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: cfg.Kubernetese.MotherServiceAPPServe},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{MatchLabels: labels},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: labels},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:            App,
								Image:           cfg.Kubernetese.MotherServiceImage,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Command:         []string{Main, Serve},
								Ports:           []corev1.ContainerPort{{ContainerPort: int32(cfg.Server.Port)}},
								EnvFrom: []corev1.EnvFromSource{
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: cfg.Kubernetese.MotherServiceAPPServe + Config}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: cfg.Kubernetese.MotherServiceAPPJobs + Secret}}},
								},
								Env: []corev1.EnvVar{
									{Name: KafkaUsername, Value: cfg.Kafka.Username},
									{Name: KafkaPassword, Value: cfg.Kafka.Password},
								},
							},
						},
					},
				},
			},
		},
	}
}

func motherServeSvcSpec(cfg *config.Config) inEntity.ServiceSpec {
	return inEntity.ServiceSpec{
		Name: cfg.Kubernetese.MotherServiceAPPServe,
		Service: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: cfg.Kubernetese.MotherServiceAPPServe},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{App: cfg.Kubernetese.MotherServiceAPPServe},
				Type:     corev1.ServiceTypeNodePort,
				Ports:    []corev1.ServicePort{{Port: int32(cfg.Server.Port), TargetPort: intstr.FromInt(cfg.Server.Port), NodePort: int32(cfg.Kubernetese.MotherServiceAPPServeNodePort)}},
			},
		},
	}
}

func motherJobsDepSpec(cfg *config.Config, replica int) inEntity.DeploymentSpec {
	replicas := int32(replica)
	labels := map[string]string{App: cfg.Kubernetese.MotherServiceAPPJobs}

	return inEntity.DeploymentSpec{
		Name: cfg.Kubernetese.MotherServiceAPPJobs,
		Deployment: &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: cfg.Kubernetese.MotherServiceAPPJobs},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{MatchLabels: labels},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: labels},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:            App,
								Image:           cfg.Kubernetese.MotherServiceImage,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Command:         []string{Main, Jobs},
								Ports:           []corev1.ContainerPort{{ContainerPort: int32(cfg.Server.Port)}},
								EnvFrom: []corev1.EnvFromSource{
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: cfg.Kubernetese.MotherServiceAPPJobs + Config}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: cfg.Kubernetese.MotherServiceAPPJobs + Secret}}},
								},
								Env: []corev1.EnvVar{
									{Name: KafkaUsername, Value: cfg.Kafka.Username},
									{Name: KafkaPassword, Value: cfg.Kafka.Password},
								},
							},
						},
					},
				},
			},
		},
	}
}
