package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
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
	scenario *entity.TestScenario,
) interfaces.TestAgentController {
	return &testAgentController{
		provisioningService:      provisioningService,
		scenario:                 scenario,
		provisioningRetries:      provisioningRetries,
		provisioningRetriesSleep: provisioningRetriesSleep,
	}
}

type testAgentController struct {
	provisioningService      provider.ProvisioningService
	scenario                 *entity.TestScenario
	uniqueID                 uuid.UUID
	provisioningRetries      int
	provisioningRetriesSleep time.Duration
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
	// TODO not implemented
	return true
}
