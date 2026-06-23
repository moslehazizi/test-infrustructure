package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
	svcMock "control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestScenarioOperationUsecase_Init(t *testing.T) {
	mockRepo := new(mocks.MockTestScenario)
	mockStressTestExecutor := new(svcMock.MockExecutionManage)

	service := NewTestScenarioOperationUsecase(
		mockRepo,
		mockStressTestExecutor,
	)
	assert.NotNil(t, service)

	assert.NotNil(t, service.testScenarioRepository)
	assert.NotNil(t, service.stressTestExecutionManager)
}

func TestTestScenarioOperationUsecase_Start(t *testing.T) {
	t.Run("failed_case_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_scenario_status_is_not_ready", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyReadyScenariosCanBeStarted)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_repository_error_on_marking_as_running", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusReady,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("RunScenario", mock.Anything, scenario).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(errors.New("something went wrong"))

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockStressTestExecutor.AssertCalled(t, "RunScenario", mock.Anything, scenario)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockRepo.AssertExpectations(t)
	})
	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusReady,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(nil)
		mockStressTestExecutor.On("RunScenario", mock.Anything, scenario).Return(nil)

		err := service.Start(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockStressTestExecutor.AssertCalled(t, "RunScenario", mock.Anything, scenario)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error_case_stress_test_adding_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusReady,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("RunScenario", mock.Anything, scenario).Return(errors.New("something went wrong"))

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockStressTestExecutor.AssertCalled(t, "RunScenario", mock.Anything, scenario)
		mockRepo.AssertExpectations(t)
	})
	t.Run("success_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusReady,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: "notsupported",
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrStartingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

}

func TestTestScenarioOperationUsecase_Pause(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_not_running", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusReady,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyRunningScenariosCanBePaused)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_error_on_marking_as_pause", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				Name: entity.STRESS,
			},
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("PauseScenario", mock.Anything, scenario, mock.Anything).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false).Return(errors.New("something went wrong"))

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockStressTestExecutor.AssertCalled(t, "PauseScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_pause_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusRunning,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("PauseScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockStressTestExecutor.AssertCalled(t, "PauseScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusRunning,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: "notsupported",
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrPausingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusRunning,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false).Return(nil)
		mockStressTestExecutor.On("PauseScenario", mock.Anything, scenario, mock.Anything).Return(nil)

		err := service.Pause(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false)
		mockStressTestExecutor.AssertCalled(t, "PauseScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioOperationUsecase_Resume(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_not_paused", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyPausedScenariosCanBeResume)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_error_on_marking_as_pause", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusPaused,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				Name: entity.STRESS,
			},
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("ResumeScenario", mock.Anything, scenario, mock.Anything).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(errors.New("something went wrong"))

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockStressTestExecutor.AssertCalled(t, "ResumeScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_resume_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPaused,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("ResumeScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockStressTestExecutor.AssertCalled(t, "ResumeScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPaused,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: "notsupported",
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrResumingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPaused,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(nil)
		mockStressTestExecutor.On("ResumeScenario", mock.Anything, scenario, mock.Anything).Return(nil)

		err := service.Resume(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockStressTestExecutor.AssertCalled(t, "ResumeScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioOperationUsecase_Stop(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_ready", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusReady,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyRunAndPauseScenariosCanBeStop)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_stop", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusStopped,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyRunAndPauseScenariosCanBeStop)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_succeed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusSucceed,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyRunAndPauseScenariosCanBeStop)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_delete", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusDeleted,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyRunAndPauseScenariosCanBeStop)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_error_on_marking_as_ready", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				Name: entity.STRESS,
			},
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("StopScenario", mock.Anything, scenario, mock.Anything).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusReady, false).Return(errors.New("something went wrong"))

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusReady, false)
		mockStressTestExecutor.AssertCalled(t, "StopScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_stop_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPaused,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockStressTestExecutor.On("StopScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockStressTestExecutor.AssertCalled(t, "StopScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPaused,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: "notsupported",
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrStoppingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPaused,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusReady, false).Return(nil)
		mockStressTestExecutor.On("StopScenario", mock.Anything, scenario, mock.Anything).Return(nil)

		err := service.Stop(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusReady, false)
		mockStressTestExecutor.AssertCalled(t, "StopScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioOperationUsecase_Delete(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		err := service.Delete(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Delete(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_not_ready", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusDeleted,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Delete(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrScenariosCanNotBeDelete)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_error_on_marking_as_deleteed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusReady,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false).Return(errors.New("something went wrong"))

		err := service.Delete(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_delete_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusReady,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false).Return(nil)
		mockStressTestExecutor.On("DeleteScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Delete(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false)
		mockStressTestExecutor.AssertCalled(t, "DeleteScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusReady,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: "notsupported",
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false).Return(nil)

		err := service.Delete(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrDeleteTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)

		service := NewTestScenarioOperationUsecase(
			mockRepo,
			mockStressTestExecutor,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusReady,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false).Return(nil)
		mockStressTestExecutor.On("DeleteScenario", mock.Anything, scenario, mock.Anything).Return(nil)

		err := service.Delete(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusDeleted, false)
		mockStressTestExecutor.AssertCalled(t, "DeleteScenario", mock.Anything, scenario, mock.Anything)
		mockRepo.AssertExpectations(t)
	})
}
