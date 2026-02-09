package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	prvMock "control-panel-service/internal/provider/mocks"
	repomock "control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewScenarioTypeRunnerGroupA(t *testing.T) {
	r := NewScenarioTypeRunnerGroupA(
		new(repomock.MockTestServiceRepository),
		new(prvMock.MockProvisioningService),
	)

	runner, ok := r.(*scenarioTypeRunnerGroupA)
	assert.True(t, ok)
	assert.NotNil(t, runner.provisioningService)
	assert.NotNil(t, runner.testServiceRepo)
}

func Test_scenarioTypeRunnerGroupA_Run(t *testing.T) {
	t.Run("failed case: scenario max service count is null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		ex := NewScenarioTypeRunnerGroupA(scRepo, provisioningService)

		err := ex.Run(context.Background(), &scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("failed case: error on getting running services count", func(t *testing.T) {
		serviceCnt := int64(10)
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		scRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(0), errors.New("something went wrong"))

		provisioningService := new(prvMock.MockProvisioningService)

		ex := NewScenarioTypeRunnerGroupA(scRepo, provisioningService)

		err := ex.Run(context.Background(), &scenario)
		assert.Error(t, err)
		scRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
	})

	t.Run("failed case: failed to provision remaining test services", func(t *testing.T) {
		serviceCnt := int64(20)
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		scRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(10), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("ProvisionTestService", mock.Anything, &scenario, int32(10)).Return(errors.New("something went wrong"))

		ex := NewScenarioTypeRunnerGroupA(scRepo, provisioningService)

		err := ex.Run(context.Background(), &scenario)

		assert.Error(t, err)
		scRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "ProvisionTestService", mock.Anything, &scenario, int32(10))
	})

	t.Run("failed case: unable to deprovision test services", func(t *testing.T) {
		serviceCnt := int64(3)
		dur := time.Millisecond * 20
		start := time.Now()
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           &start,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt)).Return(errors.New("something went wrong"))

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService)

		err := ex.Run(context.Background(), &scenario)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToDeprovisionTestServices)
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt))
		assert.GreaterOrEqual(t, time.Now(), start)
	})
	t.Run("success case: waiting for execution duration to be spent", func(t *testing.T) {
		serviceCnt := int64(3)
		dur := time.Millisecond * 2000
		start := time.Now()
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           &start,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt)).Return(nil)

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService)

		err := ex.Run(context.Background(), &scenario)
		assert.NoError(t, err)
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt))
		assert.GreaterOrEqual(t, time.Now(), start)
	})
	t.Run("success case: start time and exec duration is null", func(t *testing.T) {
		serviceCnt := int64(3)
		start := time.Now()
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   nil,
			StartedAt:           nil,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt)).Return(nil)

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService)

		err := ex.Run(context.Background(), &scenario)
		assert.NoError(t, err)
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt))
		assert.GreaterOrEqual(t, time.Now(), start)
	})
}
