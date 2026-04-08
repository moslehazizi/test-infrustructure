package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	prvMock "control-panel-service/internal/provider/mocks"
	"control-panel-service/internal/repository/mocks"
	repoMocks "control-panel-service/internal/repository/mocks"
	"control-panel-service/internal/server/dto/request"
	svcMock "control-panel-service/internal/usecase/mocks"
	svcMocks "control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTestScenarioUsecase_Init(t *testing.T) {
	mockRepo := new(mocks.MockTestScenario)
	mockTestCatRepo := new(mocks.MockTestCategory)
	mockTestServiceConfig := new(mocks.MockTestServiceConfig)
	mockMotherService := new(mocks.MockMotherService)
	mockStressTestExecutor := new(svcMock.MockExecutionManage)
	mockTestServiceRepo := new(mocks.MockTestServiceRepository)
	provisioningService := new(prvMock.MockProvisioningService)

	service := NewTestScenarioUsecase(
		getMockDB(t),
		mockRepo,
		mockTestCatRepo,
		mockTestServiceConfig,
		mockMotherService,
		mockStressTestExecutor,
		mockTestServiceRepo,
		provisioningService,
	)
	assert.NotNil(t, service)

	st, ok := service.(*testScenario)
	assert.True(t, ok)
	assert.NotNil(t, st.testScenarioRepository)
	assert.NotNil(t, st.stressTestExecutionManager)
	assert.NotNil(t, st.testScenarioRepository)
	assert.NotNil(t, st.provisioningService)
}

func TestTestScenarioUsecase_Create(t *testing.T) {
	t.Run("failed_case_when_create_test_service_config", func(t *testing.T) {
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt := 1
		sampleInt64 := int64(1)
		expectedID := uint64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:         expectedID,
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}

		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("Create", mock.Anything, testSci).Return(expectedID, nil)

		mockTestServiceConfig.On("Create", mock.Anything, testServiceConfig).Return(errors.New("error happened"))

		err := service.Create(context.Background(), testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateTestScenario)
		mockRepo.AssertCalled(t, "Create", mock.Anything, testSci)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockMotherService.AssertCalled(t, "GetByID", mock.Anything, testSci.MotherServiceID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_test_service_config_could_not_be_null", func(t *testing.T) {
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt64 := int64(1)

		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   nil,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(context.Background(), testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestServiceConfigIsRequired)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleInt := 1
		sampleInt64 := int64(1)
		expectedID := uint64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:         expectedID,
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testCat := &entity.TestCategory{
			ID:                     1,
			Name:                   entity.STRESS,
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		testSci := &entity.TestScenario{
			ID:                  uint64(4),
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			TestCategory:        testCat,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("Create", mock.Anything, testSci).Return(expectedID, nil)
		mockTestServiceConfig.On("Create", mock.Anything, testServiceConfig).Return(nil)

		err := service.Create(ctx, testSci)

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)

	})

	t.Run("failed_case_when_create_test_scenario", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)
		databaseName := "db1"
		databaseTableName := "factorial"

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleInt := 1
		sampleInt64 := int64(1)
		testServiceCfg := &entity.TestServiceConfig{
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceCfg,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockRepo.On("Create", mock.Anything, testSci).Return(uint64(0), errors.New("error happened"))
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateTestScenario)
		mockRepo.AssertCalled(t, "Create", mock.Anything, testSci)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_validation_error_max_test_service_count_less_than_one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt64 := int64(-1)
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_error_get_test_category_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(nil, pkg.ErrTestCategoryNotFound)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestCategoryNotFound)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_error_get_test_category_unknown", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(nil, errors.New("error happened"))

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestCategory)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_validation_error_max_service_count_not_set", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_validation_error_test_service_config", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		databaseName := "db1"
		databaseTableName := "factorial"
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:            1,
				MaxDuration:            0,
				BadValueRate:           -1,
				DatabaseName:           databaseName,
				DatabaseTableName:      databaseTableName,
				IncreaseFixedInput:     1,
				ExecNumMultiFixedInput: 1,
			},
			NumSteps: 2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.Equal(t, err, fmt.Errorf("%w: %w", pkg.ErrFailedToValidateTestSvcCfg, pkg.ErrInvalidBadValueRate))
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_mother_service_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt := 1
		sampleInt64 := int64(1)
		expectedID := uint64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:         expectedID,
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     0,
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(nil, pkg.ErrMotherServiceNotFound)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceNotFound)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
	})
	t.Run("success_case_when_test_category_has_num_steps_is_false", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleInt := 1
		sampleInt64 := int64(1)
		expectedID := uint64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:         expectedID,
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testCat := &entity.TestCategory{
			ID:                     1,
			Name:                   entity.STRESS,
			HasMaxTestServiceCount: true,
			HasNumSteps:            false,
		}
		testSci := &entity.TestScenario{
			ID:                  uint64(4),
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			TestCategory:        testCat,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("Create", mock.Anything, testSci).Return(expectedID, nil)
		mockTestServiceConfig.On("Create", mock.Anything, testServiceConfig).Return(nil)

		err := service.Create(ctx, testSci)

		assert.Nil(t, err)
		assert.Equal(t, testSci.NumSteps, int64(1))
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_GetByID(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		someTime := time.Date(2026, 01, 13, 14, 10, 0, 0, time.Now().Location())
		num := 10
		sampleInt64 := int64(1)
		sampleID := uint64(4)

		expectedTestScenario := &entity.TestScenario{
			ID:              sampleID,
			Name:            "load1",
			TestCategoryID:  uint64(1),
			MotherServiceID: uint64(2),
			NumSteps:        2,
			TestCategory: &entity.TestCategory{
				ID:                     1,
				Name:                   "load",
				Label:                  "Load Test",
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				HasMaxTestServiceCount: true,
				HasNumSteps:            true,
			},
			MotherService: &entity.MotherService{
				ID:   2,
				Name: "m2",
			},
			Status:              entity.ScenarioStatusSucceed,
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                    100,
				TestScenarioID:        1,
				CreatedAt:             someTime,
				UpdatedAt:             someTime,
				MaxRequests:           100,
				MaxDuration:           0,
				RequestDelayDuration:  nil,
				RandomRequestDelayMin: nil,
				RandomRequestDelayMax: nil,
				FixedTestNumber:       &num,
				RandomTestNumberMin:   nil,
				RandomTestNumberMax:   nil,
				BadValueRate:          0,
				NegativeValueRate:     0,
				RealValueRate:         0,
				ZeroValueRate:         0,
				StringValueRate:       0,
				LongStringValueRate:   0,
				NullValueRate:         0,
			},
			Editable: true,
		}

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(expectedTestScenario, nil)

		result, err := service.GetByID(ctx, sampleID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario, result)
		assert.Equal(t, expectedTestScenario.ID, result.ID)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_GetPaginated(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt64 := int64(2)
		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}
		count := int64(2)
		expectedResult := []*entity.TestScenario{
			{
				ID:              uint64(2),
				Name:            "load1",
				TestCategoryID:  uint64(1),
				MotherServiceID: uint64(2),
				NumSteps:        2,
				TestCategory: &entity.TestCategory{
					ID:                     1,
					Name:                   "load",
					Label:                  "Load Test",
					CreatedAt:              time.Now(),
					UpdatedAt:              time.Now(),
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   2,
					Name: "m2",
				},
				Status:              entity.ScenarioStatusSucceed,
				MaxTestServiceCount: &sampleInt64,
				Editable:            true,
			},
			{
				ID:              uint64(1),
				Name:            "smoke2",
				TestCategoryID:  uint64(1),
				MotherServiceID: uint64(2),
				TestCategory: &entity.TestCategory{
					ID:                     1,
					Name:                   "load",
					Label:                  "Load Test",
					CreatedAt:              time.Now(),
					UpdatedAt:              time.Now(),
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   2,
					Name: "m2",
				},
				Status:              entity.ScenarioStatusSucceed,
				MaxTestServiceCount: &sampleInt64,
				Editable:            true,
			},
		}

		mockRepo.On("GetPaginated", mock.Anything, pagReq).Return(expectedResult, count, nil)

		result, total, err := service.GetPaginated(ctx, pagReq)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, count)
		assert.Equal(t, len(result), 2)

		assert.Equal(t, expectedResult, result)

		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, pagReq)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, pagReq).Return(nil, int64(0), errors.New("error happened"))

		result, count, err := service.GetPaginated(ctx, pagReq)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenarios)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, pagReq)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_Start(t *testing.T) {
	t.Run("failed_case_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_scenario_status_is_not_pending", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		assert.ErrorIs(t, err, pkg.ErrOnlyPendingScenariosCanBeStarted)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_repository_error_on_marking_as_running", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusPending,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(errors.New("something went wrong"))

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockRepo.AssertExpectations(t)
	})
	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPending,
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
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockStressTestExecutor.AssertCalled(t, "RunScenario", mock.Anything, scenario)

		mockRepo.AssertExpectations(t)
	})
	t.Run("error_case_stress_test_adding_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPending,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: entity.STRESS,
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(nil)

		mockStressTestExecutor.On("RunScenario", mock.Anything, scenario).Return(errors.New("something went wrong"))

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockStressTestExecutor.AssertCalled(t, "RunScenario", mock.Anything, scenario)

		mockRepo.AssertExpectations(t)
	})
	t.Run("success_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPending,
			TestCategory: &entity.TestCategory{
				ID:   1,
				Name: "notsupported",
			},
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(nil)

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrStartingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)

		mockRepo.AssertExpectations(t)
	})

}

func TestTestScenarioUsecase_Pause(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusPending,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyRunningScenariosCanBePaused)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_error_on_marking_as_puase", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false).Return(errors.New("something went wrong"))

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_pause_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockStressTestExecutor.On("PauseScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false)
		mockStressTestExecutor.AssertCalled(t, "PauseScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false).Return(nil)

		err := service.Pause(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrPauseingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPaused, false)

		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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

func TestTestScenarioUsecase_Resune(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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

	t.Run("failed_case_repository_error_on_marking_as_puase", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusPaused,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(errors.New("something went wrong"))

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_resume_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockStressTestExecutor.On("ResumeScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockStressTestExecutor.AssertCalled(t, "ResumeScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(nil)

		err := service.Resume(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrResumeingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)

		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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

func TestTestScenarioUsecase_Stop(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_pending", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusPending,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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

	t.Run("failed_case_scenario_status_is_abort", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusAborted,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrOnlyRunAndPauseScenariosCanBeStop)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_error_on_marking_as_pending", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false).Return(errors.New("something went wrong"))

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_stop_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false).Return(nil)
		mockStressTestExecutor.On("StopScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false)
		mockStressTestExecutor.AssertCalled(t, "StopScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false).Return(nil)

		err := service.Stop(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrStopingTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false)

		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false).Return(nil)
		mockStressTestExecutor.On("StopScenario", mock.Anything, scenario, mock.Anything).Return(nil)

		err := service.Stop(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusPending, false)
		mockStressTestExecutor.AssertCalled(t, "StopScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_Abort(t *testing.T) {
	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		err := service.Abort(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		err := service.Abort(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_scenario_status_is_abort", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusAborted,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)

		err := service.Abort(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrAbortedScenariosCanBeAbort)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_error_on_marking_as_aborted", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:       sampleID,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false).Return(errors.New("something went wrong"))

		err := service.Abort(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_stress_test_abort_scenario_failed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false).Return(nil)
		mockStressTestExecutor.On("AbortScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Abort(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false)
		mockStressTestExecutor.AssertCalled(t, "AbortScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_implemented", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false).Return(nil)

		err := service.Abort(ctx, sampleID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrAbortTestNotImplemented)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false)

		mockRepo.AssertExpectations(t)
	})

	t.Run("success_case_stress_test", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockStressTestExecutor := new(svcMock.MockExecutionManage)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			provisioningService,
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
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false).Return(nil)
		mockStressTestExecutor.On("AbortScenario", mock.Anything, scenario, mock.Anything).Return(nil)

		err := service.Abort(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusAborted, false)
		mockStressTestExecutor.AssertCalled(t, "AbortScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_DeprovisionAllPods(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		ctx := context.Background()

		sc := &entity.TestScenario{
			ID:       1,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
		}

		mockRepo.On("GetByStatus", mock.Anything, entity.ScenarioStatusRunning).Return([]*entity.TestScenario{sc}, nil)
		mockProvisioningService.On("DeprovisionTestService", mock.Anything, sc, mock.Anything).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, sc.ID, entity.ScenarioStatusAborted, false).Return(nil)

		err := service.DeprovisionAllPods(ctx)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_failed_to_get_by_status", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		ctx := context.Background()

		sc := &entity.TestScenario{
			ID:       1,
			Status:   entity.ScenarioStatusRunning,
			NumSteps: 2,
		}

		mockRepo.On("GetByStatus", mock.Anything, entity.ScenarioStatusRunning).Return([]*entity.TestScenario{sc}, errors.New("error happened"))

		err := service.DeprovisionAllPods(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenariosByStatus)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_Update(t *testing.T) {
	sampleCategory := &entity.TestCategory{
		ID:                     1,
		Name:                   "stress",
		HasMaxTestServiceCount: false,
		HasNumSteps:            true,
	}

	baseExistingScenario := func() *entity.TestScenario {
		return &entity.TestScenario{
			ID:              1,
			Name:            "Old Name",
			MotherServiceID: 10,
			MotherService:   &entity.MotherService{ID: 10, Name: "Old Mother"},
			TestCategory:    sampleCategory,
			TestServiceConfig: &entity.TestServiceConfig{
				ID: 100,
			},
			NumSteps: 2,
		}
	}

	baseMotherService := &entity.MotherService{ID: 12}

	baseConfig := func() *entity.TestServiceConfig {
		return &entity.TestServiceConfig{ID: 100}
	}

	fixedTestNumber := 5
	requestDelay := 180

	baseRequest := func() *request.TestScenarioUpdateRequest {
		return &request.TestScenarioUpdateRequest{
			ID:                  1,
			Name:                "Updated Load Test Scenario",
			MotherServiceID:     12,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
			Config: &request.TestServiceConfigRequest{
				MaxRequests:            1000,
				MaxDuration:            300,
				RequestDelayDuration:   &requestDelay,
				FixedTestNumber:        &fixedTestNumber,
				BadValueRate:           0,
				NegativeValueRate:      0,
				RealValueRate:          0,
				ZeroValueRate:          0,
				StringValueRate:        0,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "test_db",
				DatabaseTableName:      "transactions",
				IncreaseFixedInput:     1,
				ExecNumMultiFixedInput: 1,
			},
		}
	}
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)
		mockTestServiceConfig.On("UpdateByScenarioID",
			mock.Anything, uint64(1), mock.AnythingOfType("*entity.TestServiceConfig"),
		).Return(nil)

		err := svc.Update(context.Background(), baseRequest())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("success_category_requires_all_optional_fields", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		catWithAll := &entity.TestCategory{
			ID:                     2,
			Name:                   "advanced",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		maxCount := int64(10)

		existing := baseExistingScenario()
		existing.TestCategory = catWithAll

		req := baseRequest()
		req.MaxTestServiceCount = &maxCount

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)
		mockTestServiceConfig.On("UpdateByScenarioID",
			mock.Anything, uint64(1), mock.AnythingOfType("*entity.TestServiceConfig"),
		).Return(nil)

		err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("success_random_delay_and_random_test_number_config", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		minDelay := 100
		maxDelay := 500
		minTest := 1
		maxTest := 10

		req := baseRequest()
		req.Config.RequestDelayDuration = nil
		req.Config.RandomRequestDelayMin = &minDelay
		req.Config.RandomRequestDelayMax = &maxDelay
		req.Config.FixedTestNumber = nil
		req.Config.RandomTestNumberMin = &minTest
		req.Config.RandomTestNumberMax = &maxTest

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)
		mockTestServiceConfig.On("UpdateByScenarioID",
			mock.Anything, uint64(1), mock.AnythingOfType("*entity.TestServiceConfig"),
		).Return(nil)

		err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("success_bad_value_rate_enabled_with_valid_sub_rates_summing_to_100", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config.BadValueRate = 30
		req.Config.NegativeValueRate = 20
		req.Config.RealValueRate = 20
		req.Config.ZeroValueRate = 20
		req.Config.StringValueRate = 10
		req.Config.LongStringValueRate = 10
		req.Config.NullValueRate = 20

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)
		mockTestServiceConfig.On("UpdateByScenarioID",
			mock.Anything, uint64(1), mock.AnythingOfType("*entity.TestServiceConfig"),
		).Return(nil)

		err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("success_optional_fields_nil_are_preserved_from_existing", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		originalMax := int64(99)
		cat := &entity.TestCategory{
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		existing := baseExistingScenario()
		existing.TestCategory = cat
		existing.MaxTestServiceCount = &originalMax

		req := baseRequest()
		req.MaxTestServiceCount = nil

		var capturedScenario *entity.TestScenario

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).
			Run(func(args mock.Arguments) {
				capturedScenario = args.Get(1).(*entity.TestScenario)
			}).Return(nil)
		mockTestServiceConfig.On("UpdateByScenarioID",
			mock.Anything, uint64(1), mock.AnythingOfType("*entity.TestServiceConfig"),
		).Return(nil)

		err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		require.NotNil(t, capturedScenario)
		assert.Equal(t, originalMax, *capturedScenario.MaxTestServiceCount,
			"MaxTestServiceCount must not be overwritten when nil is passed")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_nil_config_returns_ErrTestServiceConfigIsRequired_immediately", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config = nil

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrTestServiceConfigIsRequired)
		mockRepo.AssertNotCalled(t, "GetByID")
		mockMotherService.AssertNotCalled(t, "GetByID")
		mockTestServiceConfig.AssertNotCalled(t, "GetByID")
	})

	t.Run("error_scenario_not_found", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(nil, errors.New("not found"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockMotherService.AssertNotCalled(t, "GetByID")
		mockTestServiceConfig.AssertNotCalled(t, "GetByID")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_mother_service_not_found", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(nil, errors.New("not found"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherService)
		mockTestServiceConfig.AssertNotCalled(t, "GetByID")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_scenario_validation_MaxTestServiceCount_<_1", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		cat := &entity.TestCategory{HasMaxTestServiceCount: true, HasNumSteps: true}
		existing := baseExistingScenario()
		existing.TestCategory = cat

		invalidCount := int64(0)
		req := baseRequest()
		req.MaxTestServiceCount = &invalidCount

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
		mockTestServiceConfig.AssertNotCalled(t, "GetByID")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_test_service_config_not_preloaded_(ID=0)", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		existing := baseExistingScenario()
		existing.TestServiceConfig = &entity.TestServiceConfig{ID: 0}

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestServiceConfig)
		mockTestServiceConfig.AssertNotCalled(t, "GetByID")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_failed_to_fetch_test_service_config", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(nil, errors.New("db error"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestServiceConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_no_delay_config_set", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config.RequestDelayDuration = nil
		req.Config.RandomRequestDelayMin = nil
		req.Config.RandomRequestDelayMax = nil

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidRequestDelayDurationConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_fixed_delay_and_random_delay_both_set", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		minDelay := 100
		maxDelay := 500

		req := baseRequest()
		// RequestDelayDuration already set — mixing with random is invalid
		req.Config.RandomRequestDelayMin = &minDelay
		req.Config.RandomRequestDelayMax = &maxDelay

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidRequestDelayDurationConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_random_delay_min_>=_max", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		minDelay := 500
		maxDelay := 100 // min > max

		req := baseRequest()
		req.Config.RequestDelayDuration = nil
		req.Config.RandomRequestDelayMin = &minDelay
		req.Config.RandomRequestDelayMax = &maxDelay

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrMinDelayDurationMoreThanMax)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_no_test_number_config_set", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config.FixedTestNumber = nil
		req.Config.RandomTestNumberMin = nil
		req.Config.RandomTestNumberMax = nil

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidTestNumberConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_fixed_test_number_<=_0", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		zero := 0
		req := baseRequest()
		req.Config.FixedTestNumber = &zero

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidFixedTestNumber)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_fixed_and_random_test_number_both_set", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		minTest := 1
		maxTest := 10

		req := baseRequest()
		// FixedTestNumber already set in baseRequest — mixing with random is invalid
		req.Config.RandomTestNumberMin = &minTest
		req.Config.RandomTestNumberMax = &maxTest

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidFixedTestNumberConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_random_test_number_min_>=_max", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		minTest := 10
		maxTest := 1 // min > max

		req := baseRequest()
		req.Config.FixedTestNumber = nil
		req.Config.RandomTestNumberMin = &minTest
		req.Config.RandomTestNumberMax = &maxTest

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrMinRandomTestNumberMoreThanMax)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_bad_value_rate_>_0_but_sub_rates_do_not_sum_to_100", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config.BadValueRate = 30
		req.Config.NegativeValueRate = 10 // sum = 10, not 100

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalid100SumOfBadValues)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_bad_value_rate_==_0_but_sub_rates_are_non_zero", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config.BadValueRate = 0
		req.Config.NegativeValueRate = 10 // must be 0 when BadValueRate == 0

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidZeroSumOfBadValues)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_empty_database_name", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config.DatabaseName = ""

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidDatabaseName)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_config_validation_empty_database_table_name", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		req := baseRequest()
		req.Config.DatabaseTableName = ""

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidDatabaseTableName)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_scenario_repository_update_fails", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).
			Return(errors.New("update failed"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		mockTestServiceConfig.AssertNotCalled(t, "UpdateByScenarioID")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})

	t.Run("error_test_service_config_update_fails", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockStressTestExecutor := new(svcMocks.MockExecutionManage)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockStressTestExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockTestServiceConfig.On("GetByID", mock.Anything, uint64(100)).Return(baseConfig(), nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)
		mockTestServiceConfig.On("UpdateByScenarioID",
			mock.Anything, uint64(1), mock.AnythingOfType("*entity.TestServiceConfig"),
		).Return(errors.New("config update failed"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestServiceConfig.AssertExpectations(t)
	})
}
