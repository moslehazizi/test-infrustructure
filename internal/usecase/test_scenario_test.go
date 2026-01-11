package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestScenarioUsecase_Init(t *testing.T) {
	mockRepo := new(mocks.MockTestScenario)
	mockTestCatRepo := new(mocks.MockTestCategory)
	service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)
	assert.NotNil(t, service)

	st, ok := service.(*testScenario)
	assert.True(t, ok)
	assert.NotNil(t, st.testScenarioRepository)
}

func TestTestScenarioUsecase_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:                 "load1",
			TestCategoryID:       2,
			MotherServiceID:      uint64(1),
			MaxTestServiceCount:  &sampleInt,
			ExecutionDuration:    &sampleInt,
			AutoStepIncreaseRate: &sampleInt,
		}
		testCat := &entity.TestCategory{
			ID:                      1,
			Name:                    "load",
			Label:                   "my load",
			HasMaxTestServiceCount:  true,
			HasExecutionDuration:    true,
			HasAutoStepIncreaseRate: true,
		}

		mockRepo.On("Create", ctx, testSci).Return(nil)
		mockTestCatRepo.On("GetByID", ctx, testSci.TestCategoryID).Return(testCat)

		err := service.Create(ctx, testSci)

		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "Create", ctx, testSci)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:                 "load1",
			TestCategoryID:       uint64(2),
			MotherServiceID:      uint64(1),
			MaxTestServiceCount:  &sampleInt,
			ExecutionDuration:    &sampleInt,
			AutoStepIncreaseRate: &sampleInt,
		}

		mockRepo.On("Create", ctx, testSci).Return(errors.New("error happened"))

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateTestScenario)
		mockRepo.AssertCalled(t, "Create", ctx, testSci)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed case - validation error - max test service count less than one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt,
		}

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
	})

	t.Run("failed case - validation error -  execution duration less than one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:              "load1",
			TestCategoryID:    uint64(2),
			MotherServiceID:   uint64(1),
			ExecutionDuration: &sampleInt,
		}

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrExecutionDurationLessThanOne)
	})

	t.Run("failed case - validation error -  auto step increase less than one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleInt := -1
		testSci := &entity.TestScenario{
			Name:                 "load1",
			TestCategoryID:       uint64(2),
			MotherServiceID:      uint64(1),
			AutoStepIncreaseRate: &sampleInt,
		}

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrAutoStepIncreaseRateLessThanOne)
	})
}

func TestTestScenarioUsecase_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleInt := 1
		sampleID := uint64(4)
		expectedTestScenario := &entity.TestScenario{
			ID:                   sampleID,
			Name:                 "load1",
			TestCategoryID:       uint64(1),
			MotherServiceID:      uint64(2),
			Status:               entity.ScenarioStatusSucceed,
			MaxTestServiceCount:  &sampleInt,
			ExecutionDuration:    &sampleInt,
			AutoStepIncreaseRate: &sampleInt,
		}

		mockRepo.On("GetByID", ctx, sampleID).Return(expectedTestScenario, nil)

		result, err := service.GetByID(ctx, sampleID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario.ID, result.ID)
		mockRepo.AssertCalled(t, "GetByID", ctx, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed case - not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleID := uint64(4)

		mockRepo.On("GetByID", ctx, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", ctx, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed case - repository unknown error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleID := uint64(4)

		mockRepo.On("GetByID", ctx, sampleID).Return(nil, errors.New("error happened"))

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", ctx, sampleID)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_GetPaginated(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		sampleInt := 2
		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}
		expectedResult := []*entity.TestScenario{
			{
				ID:                   uint64(2),
				Name:                 "load1",
				TestCategoryID:       uint64(1),
				MotherServiceID:      uint64(2),
				Status:               entity.ScenarioStatusSucceed,
				MaxTestServiceCount:  &sampleInt,
				ExecutionDuration:    &sampleInt,
				AutoStepIncreaseRate: &sampleInt,
			},
			{
				ID:                   uint64(1),
				Name:                 "smoke2",
				TestCategoryID:       uint64(1),
				MotherServiceID:      uint64(2),
				Status:               entity.ScenarioStatusSucceed,
				MaxTestServiceCount:  &sampleInt,
				ExecutionDuration:    &sampleInt,
				AutoStepIncreaseRate: &sampleInt,
			},
		}

		mockRepo.On("GetPaginated", ctx, pagReq).Return(expectedResult, nil)

		result, err := service.GetPaginated(ctx, pagReq)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(result), 2)

		assert.Equal(t, expectedResult[0].ID, result[0].ID)
		assert.Equal(t, expectedResult[0].Name, result[0].Name)
		assert.Equal(t, expectedResult[0].MotherServiceID, result[0].MotherServiceID)

		assert.Equal(t, expectedResult[1].ID, result[1].ID)
		assert.Equal(t, expectedResult[1].Name, result[1].Name)
		assert.Equal(t, expectedResult[1].MotherServiceID, result[1].MotherServiceID)

		mockRepo.AssertCalled(t, "GetPaginated", ctx, pagReq)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		service := NewTestScenarioUsecase(mockRepo, mockTestCatRepo)

		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", ctx, pagReq).Return(nil, errors.New("error happened"))

		result, err := service.GetPaginated(ctx, pagReq)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenarios)
		mockRepo.AssertCalled(t, "GetPaginated", ctx, pagReq)
		mockRepo.AssertExpectations(t)
	})
}
