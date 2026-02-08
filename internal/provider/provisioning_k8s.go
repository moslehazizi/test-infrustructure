package provider

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	kubernetese "control-panel-service/pkg/kubernetes"
	inEntity "control-panel-service/pkg/kubernetes/domain/entity"
	"control-panel-service/pkg/number"

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

func NewProvisioningService(cfg *config.Config, kubernetes kubernetese.Kubernetese) ProvisioningService {
	return &provisioningService{
		cfg:        cfg,
		kubernetes: kubernetes,
	}
}

type provisioningService struct {
	cfg        *config.Config
	kubernetes kubernetese.Kubernetese
}

func (ps *provisioningService) ProvisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error {
	return nil
}

func (ps *provisioningService) DeprovisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error {
	return nil
}

func (ps *provisioningService) ProvisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	// Mother service serve deployment logic here
	// chekc mother service
	if motherService == nil {
		return fmt.Errorf("%w", pkg.ErrMotherServiceIsNil)

	}

	// Create config map
	randomResponseDelayMin, randomResponseDelayMax := number.CalcRange(motherService.ResponseDelayDuration, motherService.RandomResponseDelayMin, motherService.RandomResponseDelayMax)
	configMap := map[string]string{
		ServiceId:                       strconv.FormatUint(motherService.ID, 10),
		HTTPPort:                        strconv.Itoa(ps.cfg.Server.Port),
		HTTPPostBodyLimit:               strconv.Itoa(ps.cfg.Server.PostBodyLimit),
		HTTPReadTimeout:                 ps.cfg.Server.ReadTimeout.String(),
		HTTPWriteTimeout:                ps.cfg.Server.WriteTimeout.String(),
		HTTPRateLimitMaxRequest:         strconv.Itoa(ps.cfg.Server.RateLimitMaxRequest),
		HTTPRateLimitExpirationduration: ps.cfg.Server.RateLimitExpirationDuration.String(),
		HTTPShutdownTimeout:             ps.cfg.Server.ShutdownTimeout.String(),
		HTTPErrorInjectionRate:          strconv.Itoa(motherService.ExceptionRate),
		HTTPDelayInjectionRate:          strconv.Itoa(motherService.ResponseDelayRate),
		HTTPDelayInjectionDurationMin:   strconv.Itoa(randomResponseDelayMin),
		HTTPDelayInjectionDurationMax:   strconv.Itoa(randomResponseDelayMax),
		LogLevel:                        ps.cfg.Kubernetese.MotherServiceLogLevel,
		LogFormat:                       ps.cfg.Kubernetese.MotherServiceLogFormat,
		LogOutput:                       ps.cfg.Kubernetese.MotherServicelogOutput,
		KafkaHost:                       ps.cfg.Kubernetese.MotherServiceKafkaHost,
		KafkaPort:                       strconv.Itoa(ps.cfg.Kafka.Port),
		KafkaDialerTimeout:              ps.cfg.Kafka.DialerTimeout.String(),
		KafkaMaxBytes:                   strconv.Itoa(ps.cfg.Kafka.MaxBytes),
		KafkaBatchTimeout:               ps.cfg.Kafka.BatchTimeout.String(),
		KafkaBatchSize:                  strconv.Itoa(ps.cfg.Kafka.BatchSize),
		KafkaBatchBytes:                 strconv.Itoa(ps.cfg.Kafka.BatchBytes),
		KafkaDatabaseTopic:              fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.MotherServiceKafkaDbTopic, motherService.ID),
		KafkaConsumerGroup:              fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.MotherServiceKafkaDbGroup, motherService.ID),
		KafkaLiveFeedTopic:              ps.cfg.Kubernetese.MotherServiceLiveFeedTopic,
		PostgresHost:                    ps.cfg.Kubernetese.MotherServicePostgresHost,
		PostgresPort:                    strconv.Itoa(ps.cfg.Postgres.Port),
		PostgresDatabase:                motherService.DatabaseName,
		PostgresTable:                   motherService.DatabaseTableName,
		PostgresSSLMode:                 ps.cfg.Postgres.SSLMode,
		PostgresMaxOpenConnection:       strconv.Itoa(ps.cfg.Postgres.MaxOpenConnections),
		PostgresMaxIdleConnection:       strconv.Itoa(ps.cfg.Postgres.MaxIdleConnections),
		PostgresConnMaxLifetime:         ps.cfg.Postgres.ConnMaxLifetime.String(),
		PostgresConnMaxIdleTime:         ps.cfg.Postgres.ConnMaxIdleTime.String(),
	}

	// Create secret
	secretMap := map[string]string{
		KafkaUsername:    ps.cfg.Kafka.Username,
		KafkaPassword:    ps.cfg.Kafka.Password,
		PostgresUser:     ps.cfg.Postgres.User,
		PostgresPassword: ps.cfg.Postgres.Password,
	}

	// Serve :
	// Create deploy spec
	serveDepSpec := motherServDepSpec(ps.cfg, Replication, motherService.ID)

	// Create service spec
	serveSvcSpec := motherServeSvcSpec(ps.cfg, motherService.ID)

	// Call ApplyDeployment from kubernetese interface
	err := ps.kubernetes.ApplyDeployment(ctx, serveDepSpec, configMap, secretMap)
	if err != nil {
		zap.L().Error("apply deployment fail",
			zap.String("apllication", serveDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	// Call ApplyService from kubernetese interface
	err = ps.kubernetes.ApplyService(ctx, serveSvcSpec, configMap, secretMap)
	if err != nil {
		zap.L().Error("apply serive fail",
			zap.String("apllication", serveDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	// Wait for deployment to be ready
	serveSvcName := fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.MotherServiceAPPServe, motherService.ID)
	err = ps.kubernetes.WaitForDeployment(ctx, serveSvcName, ps.cfg.Kubernetese.MotherServiceAPPServeWaitReady)
	if err != nil {
		zap.L().Error("create pod fail",
			zap.String("apllication", serveDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	// jobs :
	// Create deploy spec
	jobsDepSpec := motherJobsDepSpec(ps.cfg, Replication)

	// Call ApplyDeployment from kubernetese interface
	err = ps.kubernetes.ApplyDeployment(ctx, jobsDepSpec, configMap, secretMap)
	if err != nil {
		zap.L().Error("apply deployment fail",
			zap.String("apllication", jobsDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	// Wait for deployment to be ready
	jobsSvcName := ps.cfg.Kubernetese.MotherServiceAPPJobs
	err = ps.kubernetes.WaitForDeployment(ctx, jobsSvcName, ps.cfg.Kubernetese.MotherServiceAPPJobsWaitReady)
	if err != nil {
		zap.L().Error("create pod fail",
			zap.String("apllication", jobsDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	return nil
}

func motherServDepSpec(cfg *config.Config, replica int32, appId uint64) inEntity.DeploymentSpec {
	replicas := replica
	appName := fmt.Sprintf("%s-%v", cfg.Kubernetese.MotherServiceAPPServe, appId)
	labels := map[string]string{App: appName}

	return inEntity.DeploymentSpec{
		Name: appName,
		Deployment: &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: appName},
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
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: appName + Config}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: appName + Secret}}},
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

func motherServeSvcSpec(cfg *config.Config, appId uint64) inEntity.ServiceSpec {
	appName := fmt.Sprintf("%s-%v", cfg.Kubernetese.MotherServiceAPPServe, appId)

	return inEntity.ServiceSpec{
		Name: appName,
		Service: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: appName},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{App: appName},
				Type:     corev1.ServiceTypeClusterIP,
				Ports:    []corev1.ServicePort{{Port: int32(cfg.Server.Port), TargetPort: intstr.FromInt(cfg.Server.Port)}},
			},
		},
	}
}

func motherJobsDepSpec(cfg *config.Config, replica int32) inEntity.DeploymentSpec {
	replicas := replica
	appName := cfg.Kubernetese.MotherServiceAPPJobs
	labels := map[string]string{App: appName}

	return inEntity.DeploymentSpec{
		Name: appName,
		Deployment: &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: appName},
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
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: appName + Config}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: appName + Secret}}},
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

func (ps *provisioningService) DeprovisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	return nil
}
