package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/pkg"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	provisioningRetries      = 3
	provisioningRetriesSleep = time.Second
	SSLIP                    = "sslip.io"
	HOST                     = "Host"
	CONTENT_TYPE             = "Content-Type"
)

func NewTestAgentController(
	provisioningService provider.ProvisioningService,
	testServiceSDK provider.SDKTestService,
	scenario *entity.TestScenario,
	serviceHost string,
	ingressHost string,
	ingressPort int,
) *testAgentController {
	return &testAgentController{
		provisioningService:      provisioningService,
		testServiceSDK:           testServiceSDK,
		scenario:                 scenario,
		provisioningRetries:      provisioningRetries,
		provisioningRetriesSleep: provisioningRetriesSleep,
		serviceHost:              serviceHost,
		ingressHost:              ingressHost,
		ingressPort:              ingressPort,
	}
}

type testAgentController struct {
	provisioningService      provider.ProvisioningService
	testServiceSDK           provider.SDKTestService
	scenario                 *entity.TestScenario
	uniqueID                 uuid.UUID
	provisioningRetries      int
	provisioningRetriesSleep time.Duration
	serviceHost              string
	ingressHost              string
	ingressPort              int
}

func (c *testAgentController) Run() error {
	ctx := context.Background()

	// provision related test service.
	tracer := otel.Tracer("testAgentController")
	_, span := tracer.Start(ctx, "Run")
	defer span.End()

	c.uniqueID = uuid.New()

	err := c.provisionTestService(ctx, c.scenario, c.uniqueID)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToRunAgentControllerDueToProvisioningFailure, err)
	}

	return nil
}

func (c *testAgentController) provisionTestService(ctx context.Context, scenario *entity.TestScenario, uniqueID uuid.UUID) error {
	tracer := otel.Tracer("testAgentController")
	_, span := tracer.Start(ctx, "provisionTestService")
	defer span.End()

	var provisioningErr error
	for i := 1; i <= c.provisioningRetries; i++ {
		err := c.provisioningService.ProvisionTestServiceByName(ctx, scenario, uniqueID)
		if err == nil {
			provisioningErr = nil

			break
		}

		provisioningErr = fmt.Errorf("%w: %w", pkg.ErrFailedToProvisionTestService, err)

		zap.L().Error(pkg.ErrFailedToProvisionTestService.Error(),
			zap.Uint64("scenarioID", scenario.ID),
			zap.String("uniqueID", uniqueID.String()),
			zap.Int("try", i),
			zap.Int("maxTry", c.provisioningRetries),
			zap.Error(err),
		)

		if i != c.provisioningRetries {
			time.Sleep(c.provisioningRetriesSleep)
		}
	}

	if provisioningErr != nil {
		zap.L().Error("provisioning test service failed after all tries",
			zap.Uint64("scenarioID", scenario.ID),
			zap.String("uniqueID", uniqueID.String()),
			zap.Error(provisioningErr),
		)

		return provisioningErr
	}

	zap.L().Debug("test services provisioned",
		zap.Uint64("scenarioID", scenario.ID),
		zap.String("uniqueID", uniqueID.String()),
	)

	return nil
}

func (c *testAgentController) Healthy() bool {
	ctx := context.Background()
	url := c.urlGenerator()
	health, err := c.testServiceSDK.Health(ctx, url)
	if err != nil {
		zap.L().Error("health error",
			zap.String("base_url", url.BaseURL),
			zap.Any("header", url.Header),
			zap.Error(err),
		)
		return false
	}
	return health.OK
}

func (c *testAgentController) ReadyForTesting() bool {
	ctx := context.Background()
	url := c.urlGenerator()
	ready, err := c.testServiceSDK.ReadyForTest(ctx, url)
	if err != nil {
		zap.L().Error("ready for test error",
			zap.String("base_url", url.BaseURL),
			zap.Any("header", url.Header),
			zap.Error(err),
		)

		return false
	}

	return ready.OK
}

func (c *testAgentController) StartTesting(ctx context.Context, req request.RunRequest) error {
	url := c.urlGenerator()
	_, err := c.testServiceSDK.RunExecute(ctx, url, req)
	if err != nil {
		zap.L().Error("run execute error",
			zap.String("base_url", url.BaseURL),
			zap.Any("header", url.Header),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (c *testAgentController) DeleteTesting(ctx context.Context) error {
	err := c.provisioningService.DeprovisionTestServiceByName(ctx, c.scenario, c.uniqueID)
	if err != nil {
		zap.L().Error("delete testing error",
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (c *testAgentController) PauseTesting(ctx context.Context) error {
	url := c.urlGenerator()
	_, err := c.testServiceSDK.Pause(ctx, url)
	if err != nil {
		zap.L().Error("pause test service error",
			zap.String("base_url", url.BaseURL),
			zap.Any("header", url.Header),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (c *testAgentController) ResumeTesting(ctx context.Context) error {
	url := c.urlGenerator()
	_, err := c.testServiceSDK.Resume(ctx, url)
	if err != nil {
		zap.L().Error("resume test service error",
			zap.String("base_url", url.BaseURL),
			zap.Any("header", url.Header),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (c *testAgentController) StopTesting(ctx context.Context) error {
	url := c.urlGenerator()
	_, err := c.testServiceSDK.Stop(ctx, url)
	if err != nil {
		zap.L().Error("stop test service error",
			zap.String("base_url", url.BaseURL),
			zap.Any("header", url.Header),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (c *testAgentController) urlGenerator() request.URL {
	return request.URL{
		BaseURL: fmt.Sprintf("%s://%s:%d", "http", c.ingressHost, c.ingressPort),
		Header: map[string]string{
			CONTENT_TYPE: "application/json",
			HOST:         fmt.Sprintf("%s-%d-%s.%s.%s", c.serviceHost, c.scenario.ID, c.uniqueID, c.ingressHost, SSLIP), // chalenge-tese-srvice-serve-{sid}-{uuid}.127.0.0.1.sslip.io
		},
	}
}
