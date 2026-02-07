package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	prvMock "control-panel-service/internal/provider/mocks"
	repomock "control-panel-service/internal/repository/mocks"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewScenarioExecutor(t *testing.T) {
	ex := NewScenarioExecutor(
		entity.TestScenario{ID: 1},
		new(repomock.MockTestServiceRepository),
		new(prvMock.MockProvisioningService),
	)
	assert.NotNil(t, ex)

	e, ok := ex.(*scenarioExecutor)
	assert.True(t, ok)
	assert.NotNil(t, e.scenario)
	assert.NotNil(t, e.testServiceRepo)
	assert.NotNil(t, e.provisioningService)
}

func Test_scenarioExecutor_GetID(t *testing.T) {
	ex := NewScenarioExecutor(
		entity.TestScenario{ID: 1},
		new(repomock.MockTestServiceRepository),
		new(prvMock.MockProvisioningService),
	)
	assert.NotNil(t, ex)

	assert.Equal(t, uint64(1), ex.GetID())
}

func Test_scenarioExecutor_GetScenario(t *testing.T) {
	sc := entity.TestScenario{ID: 1}
	ex := NewScenarioExecutor(
		sc,
		new(repomock.MockTestServiceRepository),
		new(prvMock.MockProvisioningService),
	)
	assert.NotNil(t, ex)

	assert.NotNil(t, ex.GetScenario())
	assert.Equal(t, sc.ID, ex.GetScenario().ID)
}

func TestNewInMemoryScenarioExecutorEngine(t *testing.T) {
	eng := NewInMemoryScenarioExecutorBox()
	assert.NotNil(t, eng)
}

func Test_inMemoryScenarioExecutorBox_Add(t *testing.T) {
	t.Parallel()

	scenarioID := uint64(1)
	mockExe := new(mocks.MockScenarioExecutorWithWG)
	mockExe.WG.Add(1)

	mockExe.On("GetID").Return(scenarioID)
	mockExe.On("GetScenario").Return(&entity.TestScenario{ID: scenarioID})
	mockExe.On("ResumeOrStart", mock.Anything)

	eng := NewInMemoryScenarioExecutorBox()
	eng.Add(mockExe)

	mockExe.WG.Wait()

	assert.True(t, eng.HasExecutor(scenarioID))
	mockExe.AssertCalled(t, "GetID")
	mockExe.AssertCalled(t, "GetScenario")
	mockExe.AssertCalled(t, "ResumeOrStart", mock.Anything)
}

func Test_inMemoryScenarioExecutorBox_HasExecutor(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		mockExe := new(mocks.MockScenarioExecutor)
		mockExe.On("GetID").Return(uint64(1))
		mockExe.On("GetScenario").Return(&entity.TestScenario{ID: 1})
		mockExe.On("ResumeOrStart", mock.Anything)

		eng := NewInMemoryScenarioExecutorBox()
		eng.Add(mockExe)

		assert.True(t, eng.HasExecutor(1))
	})
	t.Run("false", func(t *testing.T) {
		eng := NewInMemoryScenarioExecutorBox()

		assert.False(t, eng.HasExecutor(1))
	})
}

func Test_scenarioExecutor_runGroupA(t *testing.T) {
	t.Run("failed case: scenario max service count is null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		ex := NewScenarioExecutor(scenario, scRepo, provisioningService).(*scenarioExecutor)

		err := ex.runGroupA(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("failed case: error on getting running services count", func(t *testing.T) {
		serviceCnt := 10
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		scRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int(0), errors.New("something went wrong"))

		provisioningService := new(prvMock.MockProvisioningService)

		ex := NewScenarioExecutor(scenario, scRepo, provisioningService).(*scenarioExecutor)

		err := ex.runGroupA(context.Background())
		assert.Error(t, err)
		scRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
	})

	t.Run("failed case: failed to provision remaining test services", func(t *testing.T) {
		serviceCnt := 20
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		scRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int(10), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("ProvisionTestService", mock.Anything, mock.Anything, 10).Return(errors.New("something went wrong"))

		ex := NewScenarioExecutor(scenario, scRepo, provisioningService).(*scenarioExecutor)

		err := ex.runGroupA(context.Background())

		assert.Error(t, err)
		scRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "ProvisionTestService", mock.Anything, mock.Anything, 10)
	})

	t.Run("success case", func(t *testing.T) {
		serviceCnt := 3
		scenario := entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
		}

		scRepo := new(repomock.MockTestServiceRepository)
		scRepo.On("GetCountAllRunningByScenario", mock.Anything, scenario.ID).Return(int(3), nil)

		provisioningService := new(prvMock.MockProvisioningService)
		provisioningService.On("DeprovisionTestService", mock.Anything, []uint64{1, 2, 3}).Return(nil)

		ex := NewScenarioExecutor(scenario, scRepo, provisioningService).(*scenarioExecutor)

		err := ex.runGroupA(context.Background())
		assert.NoError(t, err)
		scRepo.AssertCalled(t, "GetCountAllRunningByScenario", mock.Anything, scenario.ID)
		provisioningService.AssertCalled(t, "DeprovisionTestService", mock.Anything, []uint64{1, 2, 3})
	})
}
