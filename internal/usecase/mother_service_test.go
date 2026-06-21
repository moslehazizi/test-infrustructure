package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"

	provisionProvider "control-panel-service/internal/provider/mocks"
	mockstem "control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	connmock "control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getMockDB(t *testing.T) database.Database {
	mockConn := new(connmock.Connection)
	db, _, err := mockConn.OpenConnection()
	require.NoError(t, err)

	return db
}

func TestNewMotherService(t *testing.T) {
	mockRepo := new(mocks.MockMotherService)
	mockTestRepo := new(mocks.MockTestScenario)
	mockSTEM := new(mockstem.MockExecutionManage)
	mockProvision := new(provisionProvider.MockProvisioningService)

	service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

	assert.NotNil(t, service)

	assert.NotNil(t, service.motherServiceRepo)
}

func TestMotherServiceUsecase_Create(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockSTEM := new(mockstem.MockExecutionManage)
		mockProvision := new(provisionProvider.MockProvisioningService)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(1), nil)
		mockProvision.On("ProvisionMotherService", mock.Anything, mock.Anything).Return(nil)

		err := service.Create(ctx, sampleMS)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
		mockProvision.AssertExpectations(t)
	})
	t.Run("failed_DeployMotherService_returns_ErrFailedToDeployMotherService", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockSTEM := new(mockstem.MockExecutionManage)
		mockProvision := new(provisionProvider.MockProvisioningService)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(1), nil)
		mockProvision.On("ProvisionMotherService", mock.Anything, mock.Anything).Return(errors.New("apply deployment failed"))

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToDeployMotherService)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
		mockProvision.AssertExpectations(t)
	})

	t.Run("failed_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(0), pkg.ErrFailedToCreateMotherService)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateMotherService)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
	})

	t.Run("failed_case_duplicate", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(0), pkg.ErrMotherServiceAlreadyExist)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceAlreadyExist)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
	})

	t.Run("failed_case_validation_error_service_name_is_missing", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		sampleMS := &entity.MotherService{
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMotherServiceName)
	})

	t.Run("failed_case_validation_error_response_delay_rete_not_be_negative", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		sampleMS := &entity.MotherService{
			Name:              "mother",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
			ExceptionRate:     10,
			ResponseDelayRate: -1,
		}

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidResponseDelayRate)
	})

	t.Run("failed_case_validation_error_exception_rate_is_negative", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		sampleMS := &entity.MotherService{
			Name:              "mother",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
			ExceptionRate:     -10,
		}

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidExceptionRate)
	})

	t.Run("failed_case_validation_error_fixed_delay_is_set_but_rate_is_0", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)
		duration := 100

		sampleMS := &entity.MotherService{
			Name:                  "mother",
			DatabaseName:          "db1",
			DatabaseTableName:     "factorial",
			ExceptionRate:         10,
			ResponseDelayRate:     0,
			ResponseDelayDuration: &duration,
		}

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDelayConfiguration)
	})
}

func TestMotherServiceUsecase_GetByID(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		inputID := uint64(1)
		expectedResult := &entity.MotherService{
			ID:                uint64(1),
			Name:              "mother1",
			Status:            entity.MotherServiceStatusRunning,
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("GetByID", mock.Anything, inputID).Return(expectedResult, nil)

		result, err := service.GetByID(ctx, inputID)

		assert.Nil(t, err)
		assert.Equal(t, result.ID, expectedResult.ID)
		assert.Equal(t, result.Name, expectedResult.Name)
		assert.Equal(t, result.Status, expectedResult.Status)
		assert.Equal(t, result.DatabaseName, expectedResult.DatabaseName)
		assert.Equal(t, result.DatabaseTableName, expectedResult.DatabaseTableName)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, inputID)
	})

	t.Run("failed_case_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		inputID := uint64(1)

		mockRepo.On("GetByID", mock.Anything, inputID).Return(nil, pkg.ErrMotherServiceNotFound)

		result, err := service.GetByID(ctx, inputID)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, inputID)
	})

	t.Run("failed_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		inputID := uint64(1)

		mockRepo.On("GetByID", mock.Anything, inputID).Return(nil, errors.New("failed to get mother service"))

		result, err := service.GetByID(ctx, inputID)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherService)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, inputID)
	})
}

func TestMotherServiceUsecase_GetPaginated(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		serviceAddress1 := "http://service1.example.com"
		serviceAddress2 := "http://service2.example.com"
		count := int64(2)

		expectedMotherServices := []*entity.MotherService{
			{
				ID:                       uint64(5),
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
				Name:                     "mother5",
				ExceptionRate:            0.0,
				ResponseDelayRate:        0.0,
				Status:                   entity.MotherServiceStatusReady,
				ServiceDeploymentAddress: &serviceAddress2,
				DatabaseName:             "test_db5",
				DatabaseTableName:        "test_table5",
			},
			{
				ID:                       uint64(4),
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
				Name:                     "mother4",
				ExceptionRate:            10,
				ResponseDelayRate:        20,
				Status:                   entity.MotherServiceStatusRunning,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db4",
				DatabaseTableName:        "test_table4",
			},
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(expectedMotherServices, count, nil)

		result, total, err := service.GetPaginated(ctx, paginationRequest)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, count)
		assert.Equal(t, expectedMotherServices, result)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})

	t.Run("failed_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(nil, int64(0), errors.New("failed to get mother services"))

		result, count, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})

	t.Run("failed_case_negative_page", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		paginationRequest := entity.PaginationRequest{
			Page:    -1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(nil, int64(0), errors.New("failed to get mother services"))

		result, count, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})
}

func TestMotherServiceUsecase_Delete(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)
		motherServiceId := uint64(1)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		motherService := &entity.MotherService{
			ID:   motherServiceId,
			Name: "mother-service-1",
		}

		testScenarios := []*entity.TestScenario{
			{
				ID:              1,
				MotherServiceID: motherServiceId,
			},
			{
				ID:              2,
				MotherServiceID: motherServiceId,
			},
		}

		mockRepo.On("GetByID", mock.Anything, motherServiceId).Return(motherService, nil).Times(1)
		mockProvision.On("DeprovisionMotherService", mock.Anything, motherService).Return(nil).Times(1)
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusDeleted).Return(nil).Times(1)
		mockTestRepo.On("GetByMotherServiceId", mock.Anything, motherService.ID).Return(testScenarios, nil).Times(1)
		mockSTEM.On("DeleteScenario", mock.Anything, testScenarios[0]).Return(nil)
		mockSTEM.On("DeleteScenario", mock.Anything, testScenarios[1]).Return(nil)
		mockTestRepo.On("SetStatus", mock.Anything, testScenarios[0].ID, entity.ScenarioStatusPending, true).Return(nil)
		mockTestRepo.On("SetStatus", mock.Anything, testScenarios[1].ID, entity.ScenarioStatusPending, true).Return(nil)

		err := service.Delete(ctx, uint64(motherServiceId))
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
		mockTestRepo.AssertExpectations(t)
		mockSTEM.AssertExpectations(t)
	})

	t.Run("fail_case_failed_to_get_mother_service", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		mockSTEM := new(mockstem.MockExecutionManage)
		motherServiceId := uint64(1)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		motherService := &entity.MotherService{}

		mockRepo.On("GetByID", mock.Anything, motherServiceId).Return(motherService, errors.New("something went wrong")).Times(1)

		err := service.Delete(ctx, uint64(motherServiceId))
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("fail_case_failed_to_deprovision_mother_service", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		motherServiceId := uint64(1)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		motherService := &entity.MotherService{}

		mockRepo.On("GetByID", mock.Anything, motherServiceId).Return(motherService, nil).Times(1)
		mockProvision.On("DeprovisionMotherService", mock.Anything, motherService).Return(errors.New("something went wrong")).Times(1)

		err := service.Delete(ctx, uint64(motherServiceId))
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
	})

	t.Run("fail_case_failed_to_update_mother_service_status", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		motherServiceId := uint64(1)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		motherService := &entity.MotherService{}

		mockRepo.On("GetByID", mock.Anything, motherServiceId).Return(motherService, nil).Times(1)
		mockProvision.On("DeprovisionMotherService", mock.Anything, motherService).Return(nil).Times(1)
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusDeleted).Return(errors.New("something went wrong")).Times(1)

		err := service.Delete(ctx, uint64(motherServiceId))
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
	})

	t.Run("fail_case_failed_to_get_test_scenarios_by_mother_id", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		motherServiceId := uint64(1)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		motherService := &entity.MotherService{}
		testScenarios := []*entity.TestScenario{}

		mockRepo.On("GetByID", mock.Anything, motherServiceId).Return(motherService, nil).Times(1)
		mockProvision.On("DeprovisionMotherService", mock.Anything, motherService).Return(nil).Times(1)
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusDeleted).Return(nil).Times(1)
		mockTestRepo.On("GetByMotherServiceId", mock.Anything, motherService.ID).Return(testScenarios, errors.New("something went wrong")).Times(1)

		err := service.Delete(ctx, uint64(motherServiceId))
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
		mockTestRepo.AssertExpectations(t)
	})

	t.Run("fail_case_failed_to_deprovision_update_status", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockTestRepo := new(mocks.MockTestScenario)
		mockProvision := new(provisionProvider.MockProvisioningService)
		motherServiceId := uint64(1)
		mockSTEM := new(mockstem.MockExecutionManage)

		service := NewMotherService(getMockDB(t), mockRepo, mockTestRepo, mockProvision, mockSTEM)

		motherService := &entity.MotherService{
			ID:   motherServiceId,
			Name: "mother-service-1",
		}

		testScenarios := []*entity.TestScenario{
			{
				ID:              1,
				MotherServiceID: motherServiceId,
			},
			{
				ID:              2,
				MotherServiceID: motherServiceId,
			},
		}

		mockRepo.On("GetByID", mock.Anything, motherServiceId).Return(motherService, nil).Times(1)
		mockProvision.On("DeprovisionMotherService", mock.Anything, motherService).Return(nil).Times(1)
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusDeleted).Return(nil).Times(1)
		mockTestRepo.On("GetByMotherServiceId", mock.Anything, motherService.ID).Return(testScenarios, nil).Times(1)
		mockSTEM.On("DeleteScenario", mock.Anything, testScenarios[0]).Return(errors.New("something went wrong"))
		mockSTEM.On("DeleteScenario", mock.Anything, testScenarios[1]).Return(nil)
		mockTestRepo.On("SetStatus", mock.Anything, testScenarios[0].ID, entity.ScenarioStatusPending, true).Return(errors.New("something went wrong"))
		mockTestRepo.On("SetStatus", mock.Anything, testScenarios[1].ID, entity.ScenarioStatusPending, true).Return(nil)

		err := service.Delete(ctx, uint64(motherServiceId))
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
		mockTestRepo.AssertExpectations(t)
		mockSTEM.AssertExpectations(t)
	})

}
