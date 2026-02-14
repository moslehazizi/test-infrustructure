package provider

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg"
	"fmt"
	"strconv"
	"time"

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
	Secret      = "-secrets"

	// Config map.
	ServiceId                       = "SERVICE_ID"
	ServiceName                     = "SERVICE_NAME"
	HTTPPort                        = "HTTP_PORT"
	SwaggerHost                     = "SWAGGER_HOST"
	SwaggerScheme                   = "SWAGGER_SCHEME"
	SwaggerDocJson                  = "SWAGGER_DOC_JSON"
	HTTPPostBodyLimit               = "HTTP_POST_BODY_LIMIT"
	HTTPReadTimeout                 = "HTTP_READ_TIMEOUT"
	HTTPWriteTimeout                = "HTTP_WRITE_TIMEOUT" // #nosec G101 -- env key name
	HTTPRateLimitMaxRequest         = "HTTP_RATE_LIMIT_MAX_REQUEST"
	HTTPRateLimitExpirationduration = "HTTP_RATE_LIMIT_EXPIRATION_DURATION"
	HTTPShutdownTimeout             = "HTTP_SHUTDOWN_TIMEOUT"
	LogLevel                        = "LOG_LEVEL"
	LogFormat                       = "LOG_FORMAT"
	LogOutput                       = "LOG_OUTPUT"
	OTLPGrpcPort                    = "OTLP_GRPC_PORT"
	OTLPGrpcHost                    = "OTLP_GRPC_HOST"
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
	// Mother service.
	HTTPErrorInjectionRate        = "HTTP_ERROR_INJECTION_RATE"
	HTTPDelayInjectionRate        = "HTTP_DELAY_INJECTION_RATE"
	HTTPDelayInjectionDurationMin = "HTTP_DELAY_INJECTION_DURATION_MIN"
	HTTPDelayInjectionDurationMax = "HTTP_DELAY_INJECTION_DURATION_MAX"
	// Test service.
	HttpMotherServiceBaseUrl = "HTTP_MOTHER_SERVICE_BASE_URL"
	HttpMotherServiceId      = "HTTP_MOTHER_SERVICE_ID"
	HttpMaxTxsCount          = "HTTP_MAX_TXS_COUNT"
	HttpMaxTxsDuration       = "HTTP_MAX_TXS_DURATION"
	HttpMinDelayBetweenTxs   = "HTTP_MIN_DELAY_BETWEEN_TXS"
	HttpMaxDelayBetweenTxs   = "HTTP_MAX_DELAY_BETWEEN_TXS"
	HttpMinInputNum          = "HTTP_MIN_INPUT_NUM"
	HttpMaxInputNum          = "HTTP_MAX_INPUT_NUM"
	HttpTotalErr             = "HTTP_TOTAL_ERR"
	HttpRealNumErr           = "HTTP_REAL_NUM_ERR"
	HttpNegativeNumErr       = "HTTP_NEGATIVE_NUM_ERR"
	HttpZeroNumErr           = "HTTP_ZERO_NUM_ERR"
	HttpShortStrErr          = "HTTP_SHORT_STR_ERR"
	HttpLongStrErr           = "HTTP_LONG_STR_ERR"
	HttpNilErr               = "HTTP_NIL_ERR"

	// Secret map.
	KafkaUsername    = "KAFKA_USERNAME"
	KafkaPassword    = "KAFKA_PASSWORD"
	PostgresUser     = "POSTGRES_USER"
	PostgresPassword = "POSTGRES_PASSWORD" // #nosec G101 -- env key name

	DelayBetweenProvisioning = 5 * time.Second
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
	// Test scenario service serve deployment logic here
	// check test scenario
	if testScenario == nil {
		return fmt.Errorf("%w", pkg.ErrTestScenarioServiceIsNil)
	}

	// check test service config
	if testScenario.TestServiceConfig == nil {
		return fmt.Errorf("%w", pkg.ErrTestServiceConfigIsNil)
	}

	// Create config map
	randomRequestDelayMin, randomRequestDelayMax := number.CalcRange(testScenario.TestServiceConfig.RequestDelayDuration, testScenario.TestServiceConfig.RandomRequestDelayMin, testScenario.TestServiceConfig.RandomRequestDelayMax)
	randomTestNumberMin, randomTestNumberMax := number.CalcRange(testScenario.TestServiceConfig.FixedTestNumber, testScenario.TestServiceConfig.RandomTestNumberMin, testScenario.TestServiceConfig.RandomTestNumberMax)

	configMap := map[string]string{
		ServiceId:                strconv.FormatUint(testScenario.ID, 10),
		ServiceName:              testScenario.Name,
		HttpMotherServiceBaseUrl: fmt.Sprintf("http://%s-%v:%v", ps.cfg.Kubernetese.MotherServiceAPPServe, testScenario.MotherServiceID, ps.cfg.Server.Port), // http://mother-service-serve-9:8080  // mother-service-serv-9.default.svc.cluster.local:8080
		HttpMotherServiceId:      strconv.FormatUint(testScenario.MotherServiceID, 10),
		HttpMaxTxsCount:          strconv.Itoa(testScenario.TestServiceConfig.MaxRequests),
		HttpMaxTxsDuration:       fmt.Sprintf("%v%s", testScenario.TestServiceConfig.MaxDuration, "ms"),
		HttpMinDelayBetweenTxs:   fmt.Sprintf("%v%s", randomRequestDelayMin, "ms"),
		HttpMaxDelayBetweenTxs:   fmt.Sprintf("%v%s", randomRequestDelayMax, "ms"),
		HttpMinInputNum:          strconv.Itoa(randomTestNumberMin),
		HttpMaxInputNum:          strconv.Itoa(randomTestNumberMax),
		HttpTotalErr:             strconv.Itoa(testScenario.TestServiceConfig.BadValueRate),
		HttpRealNumErr:           strconv.Itoa(testScenario.TestServiceConfig.RealValueRate),
		HttpNegativeNumErr:       strconv.Itoa(testScenario.TestServiceConfig.NegativeValueRate),
		HttpZeroNumErr:           strconv.Itoa(testScenario.TestServiceConfig.ZeroValueRate),
		HttpShortStrErr:          strconv.Itoa(testScenario.TestServiceConfig.StringValueRate),
		HttpLongStrErr:           strconv.Itoa(testScenario.TestServiceConfig.LongStringValueRate),
		HttpNilErr:               strconv.Itoa(testScenario.TestServiceConfig.NullValueRate),
		PostgresDatabase:         testScenario.TestServiceConfig.DatabaseName,
		PostgresTable:            testScenario.TestServiceConfig.DatabaseTableName,

		KafkaHost:          ps.cfg.Kubernetese.TestServiceKafkaHost,
		PostgresHost:       ps.cfg.Kubernetese.TestServicePostgresHost,
		KafkaLiveFeedTopic: ps.cfg.Kubernetese.TestServiceLiveFeedTopic,
		KafkaDatabaseTopic: fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.TestServiceKafkaDatabaseTopic, testScenario.ID),
		KafkaConsumerGroup: fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.TestServiceKafkaCounsumerGroup, testScenario.ID),

		HTTPPort:                        strconv.Itoa(ps.cfg.Server.Port),
		SwaggerHost:                     ps.cfg.Server.SwaggerHost,
		SwaggerScheme:                   ps.cfg.Server.SwaggerScheme[0],
		SwaggerDocJson:                  ps.cfg.Server.SwaggerDocJSON,
		HTTPPostBodyLimit:               strconv.Itoa(ps.cfg.Server.PostBodyLimit),
		HTTPReadTimeout:                 ps.cfg.Server.ReadTimeout.String(),
		HTTPWriteTimeout:                ps.cfg.Server.WriteTimeout.String(),
		HTTPRateLimitMaxRequest:         strconv.Itoa(ps.cfg.Server.RateLimitMaxRequest),
		HTTPRateLimitExpirationduration: ps.cfg.Server.RateLimitExpirationDuration.String(),
		HTTPShutdownTimeout:             ps.cfg.Server.ShutdownTimeout.String(),
		KafkaPort:                       strconv.Itoa(ps.cfg.Kafka.Port),
		KafkaDialerTimeout:              ps.cfg.Kafka.DialerTimeout.String(),
		KafkaMaxBytes:                   strconv.Itoa(ps.cfg.Kafka.MaxBytes),
		KafkaBatchTimeout:               ps.cfg.Kafka.BatchTimeout.String(),
		KafkaBatchSize:                  strconv.Itoa(ps.cfg.Kafka.BatchSize),
		KafkaBatchBytes:                 strconv.Itoa(ps.cfg.Kafka.BatchBytes),
		PostgresPort:                    strconv.Itoa(ps.cfg.Postgres.Port),
		PostgresSSLMode:                 ps.cfg.Postgres.SSLMode,
		PostgresMaxOpenConnection:       strconv.Itoa(ps.cfg.Postgres.MaxOpenConnections),
		PostgresMaxIdleConnection:       strconv.Itoa(ps.cfg.Postgres.MaxIdleConnections),
		PostgresConnMaxLifetime:         ps.cfg.Postgres.ConnMaxLifetime.String(),
		PostgresConnMaxIdleTime:         ps.cfg.Postgres.ConnMaxIdleTime.String(),
		LogLevel:                        ps.cfg.Logger.Level,
		LogFormat:                       ps.cfg.Logger.Format,
		LogOutput:                       ps.cfg.Logger.Output,
		OTLPGrpcPort:                    strconv.Itoa(ps.cfg.Otlp.GRPCPort),
		OTLPGrpcHost:                    ps.cfg.Otlp.GRPCHost,
	}

	// Create secret
	secretMap := map[string]string{
		KafkaUsername:    ps.cfg.Kafka.Username,
		KafkaPassword:    ps.cfg.Kafka.Password,
		PostgresUser:     ps.cfg.Postgres.User,
		PostgresPassword: ps.cfg.Postgres.Password,
	}

	// jobs :
	// Create deploy spec
	jobsDepSpec := testJobsDepSpec(ps.cfg, Replication, testScenario.ID)

	// Call ApplyDeployment from kubernetese interface
	err := ps.kubernetes.ApplyDeployment(ctx, jobsDepSpec, configMap, secretMap)
	if err != nil {
		zap.L().Error("apply deployment fail",
			zap.String("apllication", jobsDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	// Wait for deployment to be ready
	jobsSvcName := fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.TestServiceAPPJobs, testScenario.ID)
	err = ps.kubernetes.WaitForDeployment(ctx, jobsSvcName, ps.cfg.Kubernetese.TestServiceAPPJobsWaitReady)
	if err != nil {
		zap.L().Error("create pod fail",
			zap.String("apllication", jobsDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	time.Sleep(DelayBetweenProvisioning)

	// Serve :
	// Create deploy spec
	serveDepSpec := testServDepSpec(ps.cfg, replica, testScenario.ID)

	// Create service spec
	serveSvcSpec := testServeSvcSpec(ps.cfg, testScenario.ID)

	// Call ApplyDeployment from kubernetese interface
	err = ps.kubernetes.ApplyDeployment(ctx, serveDepSpec, configMap, secretMap)
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
			zap.String("apllication", serveSvcSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	// Wait for deployment to be ready
	serveSvcName := fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.TestServiceAPPServe, testScenario.ID)
	err = ps.kubernetes.WaitForDeployment(ctx, serveSvcName, ps.cfg.Kubernetese.TestServiceAPPServeWaitReady)
	if err != nil {
		zap.L().Error("create pod fail",
			zap.String("apllication", serveDepSpec.Name),
			zap.String("error", err.Error()),
		)

		return err
	}

	return nil
}

func (ps *provisioningService) DeprovisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error {
	if testScenario == nil {
		return fmt.Errorf("%w", pkg.ErrTestScenarioServiceIsNil)
	}

	// check test service config
	if testScenario.TestServiceConfig == nil {
		return fmt.Errorf("%w", pkg.ErrTestServiceConfigIsNil)
	}

	serveName := fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.TestServiceAPPServe, testScenario.ID)
	currentReplicas, err := ps.kubernetes.GetDeploymentReplicas(ctx, serveName)
	if err != nil {
		return err
	}

	newReplicas := currentReplicas - replica
	if newReplicas < 0 {
		newReplicas = 0
	}

	if newReplicas > 0 {
		return ps.kubernetes.ScaleDeployment(ctx, serveName, newReplicas)
	}

	// remove all: delete deployment (pods cascade), service, configmap, secret
	configName := serveName + Config
	secretName := serveName + Secret

	if err := ps.kubernetes.DeleteDeployment(ctx, serveName); err != nil {
		return err
	}
	if err := ps.kubernetes.DeleteService(ctx, serveName); err != nil {
		return err
	}
	if err := ps.kubernetes.DeleteConfigMap(ctx, configName); err != nil {
		return err
	}
	if err := ps.kubernetes.DeleteSecret(ctx, secretName); err != nil {
		return err
	}

	return nil
}

func (ps *provisioningService) ProvisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	// Mother service serve deployment logic here
	// check mother service
	if motherService == nil {
		return fmt.Errorf("%w", pkg.ErrMotherServiceIsNil)
	}

	// Create config map
	randomResponseDelayMin, randomResponseDelayMax := number.CalcRange(motherService.ResponseDelayDuration, motherService.RandomResponseDelayMin, motherService.RandomResponseDelayMax)
	configMap := map[string]string{
		ServiceId:                     strconv.FormatUint(motherService.ID, 10),
		ServiceName:                   motherService.Name,
		HTTPErrorInjectionRate:        strconv.Itoa(motherService.ExceptionRate),
		HTTPDelayInjectionRate:        strconv.Itoa(motherService.ResponseDelayRate),
		HTTPDelayInjectionDurationMin: strconv.Itoa(randomResponseDelayMin),
		HTTPDelayInjectionDurationMax: strconv.Itoa(randomResponseDelayMax),
		PostgresDatabase:              motherService.DatabaseName,
		PostgresTable:                 motherService.DatabaseTableName,

		KafkaHost:          ps.cfg.Kubernetese.MotherServiceKafkaHost,
		PostgresHost:       ps.cfg.Kubernetese.MotherServicePostgresHost,
		KafkaLiveFeedTopic: ps.cfg.Kubernetese.MotherServiceLiveFeedTopic,
		KafkaDatabaseTopic: fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.MotherServiceKafkaDbTopic, motherService.ID),
		KafkaConsumerGroup: fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.MotherServiceKafkaDbGroup, motherService.ID),

		HTTPPort:                        strconv.Itoa(ps.cfg.Server.Port),
		SwaggerHost:                     ps.cfg.Server.SwaggerHost,
		SwaggerScheme:                   ps.cfg.Server.SwaggerScheme[0],
		SwaggerDocJson:                  ps.cfg.Server.SwaggerDocJSON,
		HTTPPostBodyLimit:               strconv.Itoa(ps.cfg.Server.PostBodyLimit),
		HTTPReadTimeout:                 ps.cfg.Server.ReadTimeout.String(),
		HTTPWriteTimeout:                ps.cfg.Server.WriteTimeout.String(),
		HTTPRateLimitMaxRequest:         strconv.Itoa(ps.cfg.Server.RateLimitMaxRequest),
		HTTPRateLimitExpirationduration: ps.cfg.Server.RateLimitExpirationDuration.String(),
		HTTPShutdownTimeout:             ps.cfg.Server.ShutdownTimeout.String(),
		KafkaPort:                       strconv.Itoa(ps.cfg.Kafka.Port),
		KafkaDialerTimeout:              ps.cfg.Kafka.DialerTimeout.String(),
		KafkaMaxBytes:                   strconv.Itoa(ps.cfg.Kafka.MaxBytes),
		KafkaBatchTimeout:               ps.cfg.Kafka.BatchTimeout.String(),
		KafkaBatchSize:                  strconv.Itoa(ps.cfg.Kafka.BatchSize),
		KafkaBatchBytes:                 strconv.Itoa(ps.cfg.Kafka.BatchBytes),
		PostgresPort:                    strconv.Itoa(ps.cfg.Postgres.Port),
		PostgresSSLMode:                 ps.cfg.Postgres.SSLMode,
		PostgresMaxOpenConnection:       strconv.Itoa(ps.cfg.Postgres.MaxOpenConnections),
		PostgresMaxIdleConnection:       strconv.Itoa(ps.cfg.Postgres.MaxIdleConnections),
		PostgresConnMaxLifetime:         ps.cfg.Postgres.ConnMaxLifetime.String(),
		PostgresConnMaxIdleTime:         ps.cfg.Postgres.ConnMaxIdleTime.String(),
		LogLevel:                        ps.cfg.Logger.Level,
		LogFormat:                       ps.cfg.Logger.Format,
		LogOutput:                       ps.cfg.Logger.Output,
		OTLPGrpcPort:                    strconv.Itoa(ps.cfg.Otlp.GRPCPort),
		OTLPGrpcHost:                    ps.cfg.Otlp.GRPCHost,
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
			zap.String("apllication", serveSvcSpec.Name),
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

	time.Sleep(DelayBetweenProvisioning)

	// jobs :
	// Create deploy spec
	jobsDepSpec := motherJobsDepSpec(ps.cfg, Replication, motherService.ID)

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
	jobsSvcName := fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.MotherServiceAPPJobs, motherService.ID)
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

func (ps *provisioningService) DeprovisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	if motherService == nil {
		return fmt.Errorf("%w", pkg.ErrMotherServiceIsNil)
	}

	serveName := fmt.Sprintf("%s-%v", ps.cfg.Kubernetese.MotherServiceAPPServe, motherService.ID)
	configName := serveName + Config
	secretName := serveName + Secret

	if err := ps.kubernetes.DeleteDeployment(ctx, serveName); err != nil {
		return err
	}
	if err := ps.kubernetes.DeleteService(ctx, serveName); err != nil {
		return err
	}
	if err := ps.kubernetes.DeleteConfigMap(ctx, configName); err != nil {
		return err
	}
	if err := ps.kubernetes.DeleteSecret(ctx, secretName); err != nil {
		return err
	}

	return nil
}

func testServDepSpec(cfg *config.Config, replica int32, appId uint64) inEntity.DeploymentSpec {
	replicas := replica
	appName := fmt.Sprintf("%s-%v", cfg.Kubernetese.TestServiceAPPServe, appId)
	configName := appName + Config
	secretName := cfg.Kubernetese.TestServiceAPPServe + Secret
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
								Name:  App,
								Image: cfg.Kubernetese.ContainerRegistryUrl + cfg.Kubernetese.TestServiceImage,
								// ImagePullPolicy: corev1.PullAlways,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Command:         []string{Main, Serve},
								Ports:           []corev1.ContainerPort{{ContainerPort: int32(cfg.Server.Port)}}, // #nosec G115 -- port from config
								EnvFrom: []corev1.EnvFromSource{
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: configName}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: secretName}}},
								},
							},
						},
					},
				},
			},
		},
	}
}

func testServeSvcSpec(cfg *config.Config, appId uint64) inEntity.ServiceSpec {
	appName := fmt.Sprintf("%s-%v", cfg.Kubernetese.TestServiceAPPServe, appId)

	return inEntity.ServiceSpec{
		Name: appName,
		Service: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: appName},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{App: appName},
				Type:     corev1.ServiceTypeClusterIP,
				Ports:    []corev1.ServicePort{{Port: int32(cfg.Server.Port), TargetPort: intstr.FromInt(cfg.Server.Port)}}, // #nosec G115 -- port from config
			},
		},
	}
}

func testJobsDepSpec(cfg *config.Config, replica int32, appId uint64) inEntity.DeploymentSpec {
	replicas := replica
	appName := fmt.Sprintf("%s-%v", cfg.Kubernetese.TestServiceAPPJobs, appId)
	configName := appName + Config
	secretName := cfg.Kubernetese.TestServiceAPPServe + Secret
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
								Name:  App,
								Image: cfg.Kubernetese.ContainerRegistryUrl + cfg.Kubernetese.TestServiceImage,
								// ImagePullPolicy: corev1.PullAlways,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Command:         []string{Main, Jobs},
								Ports:           []corev1.ContainerPort{{ContainerPort: int32(cfg.Server.Port)}}, // #nosec G115 -- port from config
								EnvFrom: []corev1.EnvFromSource{
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: configName}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: secretName}}},
								},
							},
						},
					},
				},
			},
		},
	}
}

func motherServDepSpec(cfg *config.Config, replica int32, appId uint64) inEntity.DeploymentSpec {
	replicas := replica
	appName := fmt.Sprintf("%s-%v", cfg.Kubernetese.MotherServiceAPPServe, appId)
	configName := appName + Config
	secretName := cfg.Kubernetese.MotherServiceAPPServe + Secret
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
								Name:  App,
								Image: cfg.Kubernetese.ContainerRegistryUrl + cfg.Kubernetese.MotherServiceImage,
								// ImagePullPolicy: corev1.PullAlways,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Command:         []string{Main, Serve},
								Ports:           []corev1.ContainerPort{{ContainerPort: int32(cfg.Server.Port)}}, // #nosec G115 -- port from config
								EnvFrom: []corev1.EnvFromSource{
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: configName}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: secretName}}},
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
				Ports:    []corev1.ServicePort{{Port: int32(cfg.Server.Port), TargetPort: intstr.FromInt(cfg.Server.Port)}}, // #nosec G115 -- port from config
			},
		},
	}
}

func motherJobsDepSpec(cfg *config.Config, replica int32, appId uint64) inEntity.DeploymentSpec {
	replicas := replica
	appName := fmt.Sprintf("%s-%v", cfg.Kubernetese.MotherServiceAPPJobs, appId)
	configName := appName + Config
	secretName := cfg.Kubernetese.MotherServiceAPPServe + Secret
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
								Name:  App,
								Image: cfg.Kubernetese.ContainerRegistryUrl + cfg.Kubernetese.MotherServiceImage,
								// ImagePullPolicy: corev1.PullAlways,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Command:         []string{Main, Jobs},
								Ports:           []corev1.ContainerPort{{ContainerPort: int32(cfg.Server.Port)}}, // #nosec G115 -- port from config
								EnvFrom: []corev1.EnvFromSource{
									{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: configName}}},
									{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: secretName}}},
								},
							},
						},
					},
				},
			},
		},
	}
}
