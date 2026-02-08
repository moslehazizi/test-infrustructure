package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"fmt"
	"math"
	"time"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func NewScenarioTypeRunnerGroupA(
	testServiceRepo repository.TestServiceRepository,
	provisioningService provider.ProvisioningService,
) interfaces.ScenarioTypeRunner {
	return &scenarioTypeRunnerGroupA{
		testServiceRepo,
		provisioningService,
	}
}

type scenarioTypeRunnerGroupA struct {
	testServiceRepo     repository.TestServiceRepository
	provisioningService provider.ProvisioningService
}

// runGroupA
// See: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md#group-a.
func (r *scenarioTypeRunnerGroupA) Run(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("scenarioTypeRunnerGroupA")
	_, span := tracer.Start(ctx, "Run")
	defer span.End()

	if scenario.MaxTestServiceCount == nil {
		zap.L().Error("max test service count value is null but required",
			zap.Uint64("scenarioID", scenario.ID),
		)

		return pkg.ErrMaxTestServiceCountNotSet
	}

	cnt, err := r.testServiceRepo.GetCountAllRunningByScenario(ctx, scenario.ID)
	if err != nil {
		zap.L().Error("failed to get running test services count by scenario",
			zap.Uint64("scenarioID", scenario.ID),
			zap.Error(err),
		)

		return fmt.Errorf("failed to get running test services count: %w", err)
	}

	remaining := *scenario.MaxTestServiceCount - cnt
	// no more test service to provision and we are done here.
	if remaining > 0 {
		if remaining > math.MaxInt32 || remaining < math.MinInt32 {
			return fmt.Errorf("cnt out of int32 range: %d", cnt)
		}
		err = r.provisioningService.ProvisionTestService(ctx, scenario, int32(remaining))
		if err != nil {
			zap.L().Error("failed to provision remaining test services",
				zap.Uint64("scenarioID", scenario.ID),
				zap.Error(err),
			)

			return fmt.Errorf("failed to provision remaining test services: %w", err)
		}
	}

	zap.L().Debug("all test services already provisioned for scenario",
		zap.Uint64("scenarioID", scenario.ID),
	)

	if scenario.StartedAt == nil {
		now := time.Now()
		scenario.StartedAt = &now
	}

	if scenario.ExecutionDuration == nil {
		defaultWait := time.Millisecond
		scenario.ExecutionDuration = &defaultWait
	}

	// wait until execution time is over.
	deadLine := scenario.StartedAt.Add(*scenario.ExecutionDuration)
	for deadLine.After(time.Now()) {
		time.Sleep(defaultExecutorWaitingTime)
	}

	items, err := r.testServiceRepo.GetRunningByScenario(ctx, scenario.ID, countUnlimited)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrGettingRunningTestServicesByScenario, err)
	}
	_ = items

	if cnt > math.MaxInt32 || cnt < math.MinInt32 {
		return fmt.Errorf("cnt out of int32 range: %d", cnt)
	}

	err = r.provisioningService.DeprovisionTestService(ctx, scenario, int32(cnt))
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToDeprovisionTestServices, err)
	}

	return nil
}
