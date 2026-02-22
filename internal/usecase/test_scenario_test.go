package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	prvMock "control-panel-service/internal/provider/mocks"
	"control-panel-service/internal/repository/mocks"
	repoMocks "control-panel-service/internal/repository/mocks"
	"control-panel-service/internal/usecase/interfaces"
	svcMock "control-panel-service/internal/usecase/mocks"
	svcMocks "control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestScenarioUsecase_Init(t *testing.T) {
	mockRepo := new(mocks.MockTestScenario)
	mockTestCatRepo := new(mocks.MockTestCategory)
	mockTestServiceConfig := new(mocks.MockTestServiceConfig)
	mockMotherService := new(mocks.MockMotherService)
	mockExecutor := new(svcMock.MockScenarioExecutorBox)
	mockTestServiceRepo := new(mocks.MockTestServiceRepository)
	provisioningService := new(prvMock.MockProvisioningService)

	service := NewTestScenarioUsecase(
		getMockDB(t),
		mockRepo,
		mockTestCatRepo,
		mockTestServiceConfig,
		mockMotherService,
		mockExecutor,
		mockTestServiceRepo,
		provisioningService,
	)
	assert.NotNil(t, service)

	st, ok := service.(*testScenario)
	assert.True(t, ok)
	assert.NotNil(t, st.testScenarioRepository)
	assert.NotNil(t, st.scenarioExecutorBox)
	assert.NotNil(t, st.testScenarioRepository)
	assert.NotNil(t, st.provisioningService)
}

func TestTestScenarioUsecase_Create(t *testing.T) {
	t.Run("failed case - when create test service config", func(t *testing.T) {
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
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
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt64 := int64(-1)
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		dur := int64(-1)
		testSci := &entity.TestScenario{
			Name:              "load1",
			TestCategoryID:    uint64(2),
			MotherServiceID:   uint64(1),
			ExecutionDuration: &dur,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt64 := int64(-1)
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			provisioningService,
		)

		sampleInt64 := int64(-1)
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleInt64 := int64(-1)
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		testSci := &entity.TestScenario{
			Name:            "load1",
			TestCategoryID:  uint64(2),
			MotherServiceID: uint64(1),
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleInt64 := int64(1)
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt64,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutor := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
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
		mockExecutorBox := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutorBox,
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
		mockExecutorBox := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutorBox,
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
		mockExecutorBox := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutorBox,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusRunning,
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
		mockExecutorBox := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutorBox,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPending,
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
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockExecutorBox := new(svcMock.MockScenarioExecutorBox)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		provisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutorBox,
			mockTestServiceRepo,
			provisioningService,
		)
		sampleID := uint64(4)

		scenario := &entity.TestScenario{
			ID:     sampleID,
			Status: entity.ScenarioStatusPending,
		}
		mockRepo.On("GetByID", mock.Anything, sampleID).Return(scenario, nil)
		mockRepo.On("SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false).Return(nil)

		mockExecutorBox.On("Add", mock.Anything)

		err := service.Start(ctx, sampleID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertCalled(t, "SetStatus", mock.Anything, sampleID, entity.ScenarioStatusRunning, false)
		mockExecutorBox.AssertCalled(t, "Add", mock.Anything)

		mockRepo.AssertExpectations(t)
	})

}

func TestTestScenarioUsecase_RestOrphanedScenarios(t *testing.T) {
	t.Run("failed case - repo error", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockExecutor := new(svcMocks.MockScenarioExecutorBox)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		mockRepo.On("GetByStatus", mock.Anything, entity.ScenarioStatusRunning).Return(nil, errors.New("db failure"))

		err := service.ResetOrphanedScenarios(context.Background())

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("skip existing executor", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockExecutor := new(svcMocks.MockScenarioExecutorBox)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		sc := &entity.TestScenario{ID: 2, Status: entity.ScenarioStatusRunning}

		mockRepo.On("GetByStatus", mock.Anything, entity.ScenarioStatusRunning).Return([]*entity.TestScenario{sc}, nil)
		mockExecutor.On("HasExecutor", sc.ID).Return(true)

		err := service.ResetOrphanedScenarios(context.Background())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockExecutor.AssertExpectations(t)
	})

	t.Run("success case - add missing executors", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockExecutor := new(svcMocks.MockScenarioExecutorBox)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		sc := &entity.TestScenario{ID: 1, Status: entity.ScenarioStatusRunning}

		mockRepo.On("GetByStatus", mock.Anything, entity.ScenarioStatusRunning).Return([]*entity.TestScenario{sc}, nil)
		mockExecutor.On("HasExecutor", sc.ID).Return(false)
		mockExecutor.On("Add", mock.MatchedBy(func(ex interfaces.ScenarioExecutor) bool {
			return ex.GetID() == sc.ID
		})).Return()

		err := service.ResetOrphanedScenarios(context.Background())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockExecutor.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_DeprovisionAllPods(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockExecutor := new(svcMocks.MockScenarioExecutorBox)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		ctx := context.Background()

		sc := &entity.TestScenario{ID: 1, Status: entity.ScenarioStatusRunning}

		mockRepo.On("GetByStatus", mock.Anything, entity.ScenarioStatusRunning).Return([]*entity.TestScenario{sc}, nil)
		mockProvisioningService.On("DeprovisionTestService", mock.Anything, sc, mock.Anything).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, sc.ID, entity.ScenarioStatusAborted, true).Return(nil)

		err := service.DeprovisionAllPods(ctx)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockExecutor.AssertExpectations(t)
	})

	t.Run("failed case - failed to get by status", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockExecutor := new(svcMocks.MockScenarioExecutorBox)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)
		mockProvisioningService := new(prvMock.MockProvisioningService)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockExecutor,
			mockTestServiceRepo,
			mockProvisioningService,
		)

		ctx := context.Background()

		sc := &entity.TestScenario{ID: 1, Status: entity.ScenarioStatusRunning}

		mockRepo.On("GetByStatus", mock.Anything, entity.ScenarioStatusRunning).Return([]*entity.TestScenario{sc}, errors.New("error happened"))

		err := service.DeprovisionAllPods(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenariosByStatus)
		mockRepo.AssertExpectations(t)
		mockExecutor.AssertExpectations(t)
	})
}
