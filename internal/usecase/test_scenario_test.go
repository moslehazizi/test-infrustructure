package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
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
	service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)
	assert.NotNil(t, service)

	st, ok := service.(*testScenario)
	assert.True(t, ok)
	assert.NotNil(t, st.testScenarioRepository)
}

func TestTestScenarioUsecase_Create(t *testing.T) {
	t.Run("failed case - when create test service config", func(t *testing.T) {
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}
		expectedID := uint64(1)

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("Create", mock.Anything, testSci).Return(expectedID, nil)

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:       expectedID,
			MaxRequests:          1,
			MaxDuration:          1,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockTestServiceConfig.On("Create", mock.Anything, testServiceConfig).Return(errors.New("error happened"))

		err := service.Create(context.Background(), testSci, testServiceConfig)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateTestScenario)
		mockRepo.AssertCalled(t, "Create", mock.Anything, testSci)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
	})
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}
		expectedID := uint64(1)

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("Create", mock.Anything, testSci).Return(expectedID, nil)

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:       expectedID,
			MaxRequests:          1,
			MaxDuration:          1,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockTestServiceConfig.On("Create", mock.Anything, testServiceConfig).Return(nil)

		err := service.Create(ctx, testSci, testServiceConfig)

		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "Create", mock.Anything, testSci)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
	})
	t.Run("failed case - when create test scenario", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}
		testServiceCfg := &entity.TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          1,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockRepo.On("Create", mock.Anything, testSci).Return(uint64(0), errors.New("error happened"))
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci, testServiceCfg)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateTestScenario)
		mockRepo.AssertCalled(t, "Create", mock.Anything, testSci)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
	})

	t.Run("failed case - validation error - max test service count less than one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci, nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})

	t.Run("failed case - validation error -  execution duration less than one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:              "load1",
			TestCategoryID:    uint64(2),
			MotherServiceID:   uint64(1),
			ExecutionDuration: &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci, nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrExecutionDurationLessThanOne)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})

	t.Run("failed case - validation error -  auto step change less than one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci, nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrAutoStepChangeRateLessThanOne)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})

	t.Run("failed case - error get test category - not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt,
		}

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(nil, pkg.ErrTestCategoryNotFound)

		err := service.Create(ctx, testSci, nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestCategoryNotFound)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})

	t.Run("failed case - error get test category - unknown", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt,
		}

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(nil, errors.New("error happened"))

		err := service.Create(ctx, testSci, nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestCategory)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})

	t.Run("failed case - validation error - max service count not set", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

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
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci, nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})
	t.Run("failed case - validation error - no need to auto step change rate", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  false,
		}
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci, nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNoNeedAutoStepChange)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})
	t.Run("failed case - validation error - test service config", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:               "load1",
			TestCategoryID:     uint64(2),
			MotherServiceID:    uint64(1),
			AutoStepChangeRate: &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  true,
		}
		testServiceConfig := &entity.TestServiceConfig{
			MaxRequests:  1,
			MaxDuration:  0,
			BadValueRate: -1,
		}
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci, testServiceConfig)

		assert.Error(t, err)
		assert.Equal(t, err, fmt.Errorf("%w: %w", pkg.ErrFailedToValidateTestSvcCfg, pkg.ErrInvalidBadValueRate))
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := 1
		sampleID := uint64(4)
		expectedTestScenario := &entity.TestScenario{
			ID:              sampleID,
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
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
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
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

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
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

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
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

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
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		sampleInt := 2
		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}
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
				MaxTestServiceCount: &sampleInt,
				ExecutionDuration:   &sampleInt,
				AutoStepChangeRate:  &sampleInt,
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
				MaxTestServiceCount: &sampleInt,
				ExecutionDuration:   &sampleInt,
				AutoStepChangeRate:  &sampleInt,
			},
		}

		mockRepo.On("GetPaginated", mock.Anything, pagReq).Return(expectedResult, nil)

		result, err := service.GetPaginated(ctx, pagReq)

		assert.Nil(t, err)
		assert.NotNil(t, result)
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
		service := NewTestScenarioUsecase(getMockDB(t), mockRepo, mockTestCatRepo, mockTestServiceConfig)

		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, pagReq).Return(nil, errors.New("error happened"))

		result, err := service.GetPaginated(ctx, pagReq)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenarios)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, pagReq)
		mockRepo.AssertExpectations(t)
	})
}
