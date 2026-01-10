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
	service := NewTestScenarioUsecase(mockRepo)
	assert.NotNil(t, service)

	st, ok := service.(*testScenario)
	assert.True(t, ok)
	assert.NotNil(t, st.testScenarioRepository)
}

func TestTestScenarioUsecase_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		service := NewTestScenarioUsecase(mockRepo)

		sampleInt := 1
		testSci := &entity.TestScenario{
			Name:                 "load1",
			TestCategoryID:       uint64(2),
			MotherServiceID:      uint64(1),
			MaxTestServiceCount:  &sampleInt,
			ExecutionDuration:    &sampleInt,
			AutoStepIncreaseRate: &sampleInt,
		}

		mockRepo.On("Create", ctx, testSci).Return(nil)

		err := service.Create(ctx, testSci)

		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "Create", ctx, testSci)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		service := NewTestScenarioUsecase(mockRepo)

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
}
