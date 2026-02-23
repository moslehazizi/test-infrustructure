package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/usecase/interfaces"
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
)

func NewTestAgentController(
	provisioningService provider.ProvisioningService,
	testServiceSDK provider.SDKTestService,
	scenario *entity.TestScenario,
	testSvcServe string,
	testSvcPort int,
) interfaces.TestAgentController {
	return &testAgentController{
		provisioningService:      provisioningService,
		testServiceSDK:           testServiceSDK,
		scenario:                 scenario,
		provisioningRetries:      provisioningRetries,
		provisioningRetriesSleep: provisioningRetriesSleep,
		testSvcServe:             testSvcServe,
		testSvcPort:              testSvcPort,
		runChan:                  make(chan request.RunRequest),
	}
}

type testAgentController struct {
	provisioningService      provider.ProvisioningService
	testServiceSDK           provider.SDKTestService
	scenario                 *entity.TestScenario
	uniqueID                 uuid.UUID
	provisioningRetries      int
	provisioningRetriesSleep time.Duration
	testSvcServe             string
	testSvcPort              int
	runChan                  chan request.RunRequest
	// abortChan
	// healthChan
	// endStepChan
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

	// TODO: complete implementation
	// listen to channels for functions call
	// switch case for do action
	// all actions are api call
	// return nil
	for {
		select {
		case <-ctx.Done():
			return nil

		case req := <-c.runChan:
			c.startTesting(ctx, req)
		}
	}
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
	baseUrl := fmt.Sprintf("%s%s-%v-%s:%v", "http://", c.testSvcServe, c.scenario.ID, c.uniqueID, c.testSvcPort) // http://chalenge-tese-srvice-serve-{sid}-{uuid}:8080
	health, err := c.testServiceSDK.Health(ctx, baseUrl)
	if err != nil {
		zap.L().Error("health error",
			zap.String("base_url", baseUrl),
			zap.Error(err),
		)
		return false
	}
	return health.OK
}

func (c *testAgentController) StartTesting(ctx context.Context, req request.RunRequest) error {
	c.runChan <- req
	return nil
}

func (c *testAgentController) startTesting(ctx context.Context, req request.RunRequest) error {
	baseUrl := fmt.Sprintf("%s%s-%v-%s:%v", "http://", c.testSvcServe, c.scenario.ID, c.uniqueID, c.testSvcPort) // http://chalenge-tese-srvice-serve-{sid}-{uuid}:8080
	_, err := c.testServiceSDK.RunExecute(ctx, baseUrl, req)
	if err != nil {
		zap.L().Error("run execute error",
			zap.String("base_url", baseUrl),
			zap.Error(err),
		)
		return err
	}
	return nil
}
