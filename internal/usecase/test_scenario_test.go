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
	t.Run("failed case - when create test service config", func(t *testing.T) {
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
		dur := int64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:       expectedID,
			MaxRequests:          1,
			MaxDuration:          1,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			DatabaseName:         databaseName,
			DatabaseTableName:    databaseTableName,
		}

		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			NumSteps:            2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
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
	t.Run("failed case - test service config could not be null", func(t *testing.T) {
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
		dur := int64(1)

		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig:   nil,
			NumSteps:            2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
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
	t.Run("success case", func(t *testing.T) {
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
		dur := int64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:       expectedID,
			MaxRequests:          1,
			MaxDuration:          1,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			DatabaseName:         databaseName,
			DatabaseTableName:    databaseTableName,
		}
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			NumSteps:            2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
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
		mockRepo.AssertCalled(t, "Create", mock.Anything, testSci)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed case - when create test scenario", func(t *testing.T) {
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
		dur := int64(1)
		testServiceCfg := &entity.TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          1,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			DatabaseName:         databaseName,
			DatabaseTableName:    databaseTableName,
		}
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig:   testServiceCfg,
			NumSteps:            2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
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

	t.Run("failed case - validation error - max test service count less than one", func(t *testing.T) {
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
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
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

	t.Run("failed case - validation error -  execution duration less than one", func(t *testing.T) {
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

		dur := int64(-1)
		testSci := &entity.TestScenario{
			Name:              "load1",
			TestCategoryID:    uint64(2),
			MotherServiceID:   uint64(1),
			ExecutionDuration: &dur,
			NumSteps:          2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
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
		assert.ErrorIs(t, err, pkg.ErrExecutionDurationLessThanOne)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("failed case - validation error -  auto step change less than one", func(t *testing.T) {
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
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
			NumSteps:           2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
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
		assert.ErrorIs(t, err, pkg.ErrAutoStepChangeRateLessThanOne)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("failed case - error get test category - not found", func(t *testing.T) {
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
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
			NumSteps:           2,
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

	t.Run("failed case - error get test category - unknown", func(t *testing.T) {
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
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
			NumSteps:           2,
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

	t.Run("failed case - validation error - max service count not set", func(t *testing.T) {
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
			Name:            "load1",
			TestCategoryID:  uint64(2),
			MotherServiceID: uint64(1),
			NumSteps:        2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  false,
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
	t.Run("failed case - validation error - no need to auto step change rate", func(t *testing.T) {
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
		sampleInt64 := int64(1)
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
			NumSteps:           2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  false,
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
		assert.ErrorIs(t, err, pkg.ErrNoNeedAutoStepChange)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed case - validation error - test service config", func(t *testing.T) {
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

		sampleInt64 := int64(1)
		databaseName := "db1"
		databaseTableName := "factorial"
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:       1,
				MaxDuration:       0,
				BadValueRate:      -1,
				DatabaseName:      databaseName,
				DatabaseTableName: databaseTableName,
			},
			NumSteps: 2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  true,
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

	t.Run("failed case - mother service not found", func(t *testing.T) {
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
		dur := int64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:       expectedID,
			MaxRequests:          1,
			MaxDuration:          1,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			DatabaseName:         databaseName,
			DatabaseTableName:    databaseTableName,
		}
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     0,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			NumSteps:            2,
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(nil, pkg.ErrMotherServiceNotFound)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceNotFound)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
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
		dur := int64(1)

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
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
			MotherService: &entity.MotherService{
				ID:   2,
				Name: "m2",
			},
			Status:              entity.ScenarioStatusSucceed,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
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

	t.Run("failed case - not found", func(t *testing.T) {
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

	t.Run("failed case - repository unknown error", func(t *testing.T) {
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

	t.Run("failed case - repository unknown error", func(t *testing.T) {
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
	t.Run("success case", func(t *testing.T) {
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
		dur := int64(2)
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
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   2,
					Name: "m2",
				},
				Status:              entity.ScenarioStatusSucceed,
				MaxTestServiceCount: &sampleInt64,
				ExecutionDuration:   &dur,
				AutoStepChangeRate:  &sampleInt64,
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
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   2,
					Name: "m2",
				},
				Status:              entity.ScenarioStatusSucceed,
				MaxTestServiceCount: &sampleInt64,
				ExecutionDuration:   &dur,
				AutoStepChangeRate:  &sampleInt64,
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

	t.Run("failed case", func(t *testing.T) {
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
	t.Run("failed case - not found", func(t *testing.T) {
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

	t.Run("failed case - repository unknown error", func(t *testing.T) {
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

	t.Run("failed case - scenario status is not pending", func(t *testing.T) {
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
	t.Run("failed case - repository error on marking as running", func(t *testing.T) {
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
	t.Run("success case - stress test", func(t *testing.T) {
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

		mockStressTestExecutor.On("AddScenario", mock.Anything, scenario, mock.Anything).Return(nil)

		err := service.Start(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockStressTestExecutor.AssertCalled(t, "AddScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})
	t.Run("error case - stress test adding failed", func(t *testing.T) {
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

		mockStressTestExecutor.On("AddScenario", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		err := service.Start(ctx, sampleID)

		assert.Error(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockStressTestExecutor.AssertCalled(t, "AddScenario", mock.Anything, scenario, mock.Anything)

		mockRepo.AssertExpectations(t)
	})

	t.Run("success case - not implemented", func(t *testing.T) {
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

func TestTestScenarioUsecase_DeprovisionAllPods(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
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

	t.Run("failed case - failed to get by status", func(t *testing.T) {
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
		HasExecutionDuration:   false,
		HasAutoStepChangeRate:  false,
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
			ID:              1,
			Name:            "Updated Load Test Scenario",
			MotherServiceID: 12,
			NumSteps:        2,
			Config: &request.TestServiceConfigRequest{
				MaxRequests:          1000,
				MaxDuration:          300,
				RequestDelayDuration: &requestDelay,
				FixedTestNumber:      &fixedTestNumber,
				BadValueRate:         0,
				NegativeValueRate:    0,
				RealValueRate:        0,
				ZeroValueRate:        0,
				StringValueRate:      0,
				LongStringValueRate:  0,
				NullValueRate:        0,
				DatabaseName:         "test_db",
				DatabaseTableName:    "transactions",
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

	t.Run("success - category requires all optional fields", func(t *testing.T) {
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
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}
		maxCount := int64(10)
		execDuration := int64(3600)
		stepRate := int64(5)

		existing := baseExistingScenario()
		existing.TestCategory = catWithAll

		req := baseRequest()
		req.MaxTestServiceCount = &maxCount
		req.ExecutionDuration = &execDuration
		req.AutoStepChangeRate = &stepRate

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

	t.Run("success - random delay and random test number config", func(t *testing.T) {
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

	t.Run("success - bad value rate enabled with valid sub-rates summing to 100", func(t *testing.T) {
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

	t.Run("success - optional fields nil are preserved from existing", func(t *testing.T) {
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
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  false,
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

	t.Run("error - nil config returns ErrTestServiceConfigIsRequired immediately", func(t *testing.T) {
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

	t.Run("error - scenario not found", func(t *testing.T) {
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

	t.Run("error - mother service not found", func(t *testing.T) {
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

	t.Run("error - scenario validation: MaxTestServiceCount < 1", func(t *testing.T) {
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

		cat := &entity.TestCategory{HasMaxTestServiceCount: true}
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

	t.Run("error - scenario validation: optional field set but category does not require it", func(t *testing.T) {
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

		execDuration := int64(3600)
		req := baseRequest()
		req.ExecutionDuration = &execDuration

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		assert.ErrorIs(t, err, pkg.ErrNoNeedExecutionDuration)
		mockTestServiceConfig.AssertNotCalled(t, "GetByID")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error - scenario validation: required field not set for category", func(t *testing.T) {
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

		cat := &entity.TestCategory{HasExecutionDuration: true}
		existing := baseExistingScenario()
		existing.TestCategory = cat

		req := baseRequest()
		req.ExecutionDuration = nil

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		assert.ErrorIs(t, err, pkg.ErrExecutionDurationNotSet)
		mockTestServiceConfig.AssertNotCalled(t, "GetByID")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error - test service config not preloaded (ID = 0)", func(t *testing.T) {
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

	t.Run("error - failed to fetch test service config", func(t *testing.T) {
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

	t.Run("error - config validation: no delay config set", func(t *testing.T) {
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

	t.Run("error - config validation: fixed delay and random delay both set", func(t *testing.T) {
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

	t.Run("error - config validation: random delay min >= max", func(t *testing.T) {
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

	t.Run("error - config validation: no test number config set", func(t *testing.T) {
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

	t.Run("error - config validation: fixed test number <= 0", func(t *testing.T) {
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

	t.Run("error - config validation: fixed and random test number both set", func(t *testing.T) {
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

	t.Run("error - config validation: random test number min >= max", func(t *testing.T) {
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

	t.Run("error - config validation: bad value rate > 0 but sub-rates do not sum to 100", func(t *testing.T) {
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

	t.Run("error - config validation: bad value rate == 0 but sub-rates are non-zero", func(t *testing.T) {
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

	t.Run("error - config validation: empty database name", func(t *testing.T) {
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

	t.Run("error - config validation: empty database table name", func(t *testing.T) {
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

	t.Run("error - scenario repository update fails", func(t *testing.T) {
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

	t.Run("error - test service config update fails", func(t *testing.T) {
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
