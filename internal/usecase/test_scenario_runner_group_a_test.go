package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	prvMock "control-panel-service/internal/provider/mocks"
	repomock "control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewScenarioTypeRunnerGroupA(t *testing.T) {
	r := NewScenarioTypeRunnerGroupA(
		new(repomock.MockTestServiceRepository),
		new(prvMock.MockProvisioningService),
		new(repomock.MockTestScenario),
	)

	runner, ok := r.(*scenarioTypeRunnerGroupA)
	assert.True(t, ok)
	assert.NotNil(t, runner.provisioningService)
	assert.NotNil(t, runner.testServiceRepo)
	assert.NotNil(t, runner.testScenarioRepo)
}

func Test_scenarioTypeRunnerGroupA_Run(t *testing.T) {
	t.Run("failed case: scenario max service count is null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)
		testScenarioRepo := new(repomock.MockTestScenario)

		ex := NewScenarioTypeRunnerGroupA(scRepo, provisioningService, testScenarioRepo)

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

		testScenarioRepo := new(repomock.MockTestScenario)

		ex := NewScenarioTypeRunnerGroupA(scRepo, provisioningService, testScenarioRepo)

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

		testScenarioRepo := new(repomock.MockTestScenario)

		ex := NewScenarioTypeRunnerGroupA(scRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)

		assert.Error(t, err)
		scRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "ProvisionTestService", mock.Anything, &scenario, int32(10))
	})

	t.Run("failed case: unable to deprovision test services", func(t *testing.T) {
		serviceCnt := int64(3)
		dur := int64(20)
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
		testScenarioRepo := new(repomock.MockTestScenario)

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToDeprovisionTestServices)
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt))
		assert.GreaterOrEqual(t, time.Now(), start)
	})
	t.Run("failed case: service count out of range", func(t *testing.T) {
		serviceCnt := int64(math.MaxInt64)
		dur := int64(20)
		start := time.Now()
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           &start,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)

		provisioningService := new(prvMock.MockProvisioningService)

		testScenarioRepo := new(repomock.MockTestScenario)

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInt32OutOfRange)
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		assert.GreaterOrEqual(t, time.Now(), start)
	})
	t.Run("failed case: service count out of range on deprovisioning", func(t *testing.T) {
		serviceCnt := int64(math.MaxInt32) + 1
		dur := int64(20)
		start := time.Now()
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           &start,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(serviceCnt, nil)

		provisioningService := new(prvMock.MockProvisioningService)
		testScenarioRepo := new(repomock.MockTestScenario)

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInt32OutOfRange)
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		assert.GreaterOrEqual(t, time.Now(), start)
	})
	t.Run("success case: waiting for execution duration to be spent", func(t *testing.T) {
		serviceCnt := int64(3)
		dur := int64(2000)
		start := time.Now()
		scenario := entity.TestScenario{
			ID:                  1,
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           &start,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt)).Return(nil)

		testScenarioRepo := new(repomock.MockTestScenario)
		testScenarioRepo.On("GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID).Return(int32(5), nil)
		testScenarioRepo.On("UpdateDeploymentNumber", mock.Anything, scenario.ID, int32(2)).Return(nil) // 5 - 3 = 2

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)
		assert.NoError(t, err)

		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt))
		testScenarioRepo.AssertCalled(t, "GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID)
		testScenarioRepo.AssertCalled(t, "UpdateDeploymentNumber", mock.Anything, scenario.ID, int32(2))
		assert.GreaterOrEqual(t, time.Now(), start)
	})

	t.Run("success case: start time and exec duration is null", func(t *testing.T) {
		serviceCnt := int64(3)
		start := time.Now()
		scenario := entity.TestScenario{
			ID:                  1,
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   nil,
			StartedAt:           nil,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt)).Return(nil)

		testScenarioRepo := new(repomock.MockTestScenario)
		testScenarioRepo.On("GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID).Return(int32(5), nil)
		testScenarioRepo.On("UpdateDeploymentNumber", mock.Anything, scenario.ID, int32(2)).Return(nil) // 5 - 3 = 2

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)
		assert.NoError(t, err)

		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt))
		testScenarioRepo.AssertCalled(t, "GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID)
		testScenarioRepo.AssertCalled(t, "UpdateDeploymentNumber", mock.Anything, scenario.ID, int32(2))
		assert.GreaterOrEqual(t, time.Now(), start)
	})

	t.Run("fail case: get deployment number error", func(t *testing.T) {
		serviceCnt := int64(3)
		scenario := entity.TestScenario{
			ID:                  1,
			MaxTestServiceCount: &serviceCnt,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)
		provisioningService.On("DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt)).Return(nil)

		testScenarioRepo := new(repomock.MockTestScenario)
		testScenarioRepo.On("GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID).
			Return(int32(0), errors.New("db is down")) // simulate error

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get deployment number")
		testScenarioRepo.AssertCalled(t, "GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID)
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
	})

	t.Run("fail case: update deployment number error", func(t *testing.T) {
		serviceCnt := int64(3)
		scenario := entity.TestScenario{
			ID:                  1,
			MaxTestServiceCount: &serviceCnt,
		}

		testServiceRepo := new(repomock.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)
		testServiceRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int64(3), nil)
		provisioningService.On("DeprovisionTestService", mock.Anything, &scenario, int32(serviceCnt)).Return(nil)

		testScenarioRepo := new(repomock.MockTestScenario)
		testScenarioRepo.On("GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID).Return(int32(5), nil)
		testScenarioRepo.On("UpdateDeploymentNumber", mock.Anything, scenario.ID, int32(2)).
			Return(errors.New("db is down")) // simulate error

		ex := NewScenarioTypeRunnerGroupA(testServiceRepo, provisioningService, testScenarioRepo)

		err := ex.Run(context.Background(), &scenario)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update deployment number")
		testServiceRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		testScenarioRepo.AssertCalled(t, "GetDeploymentNumberByScenarioID", mock.Anything, scenario.ID)
		testScenarioRepo.AssertCalled(t, "UpdateDeploymentNumber", mock.Anything, scenario.ID, int32(2))
	})

}
