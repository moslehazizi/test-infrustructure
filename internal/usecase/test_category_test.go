package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTestCategoryService(t *testing.T) {
	testCategoryRepo := new(mocks.MockTestCategory)
	srv := NewTestCategoryService(testCategoryRepo)

	s, ok := srv.(*testCategoryService)
	assert.True(t, ok)
	assert.NotNil(t, s.testCategoryRepo)
}

func TestGetAll(t *testing.T) {
	t.Run("success case: get all empty list", func(t *testing.T) {
		testCategoryRepo := new(mocks.MockTestCategory)
		expected := []entity.TestCategory{
			{
				ID:                      1,
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
				Name:                    "load",
				Label:                   "Load Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
			},
			{
				ID:                      2,
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
				Name:                    "smoke",
				Label:                   "Smoke Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
			},
		}
		testCategoryRepo.On("GetAll", mock.Anything).Return(expected, nil)
		srv := NewTestCategoryService(testCategoryRepo)

		items, err := srv.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.True(t, reflect.DeepEqual(items, expected))
	})

	t.Run("failed case", func(t *testing.T) {
		testCategoryRepo := new(mocks.MockTestCategory)
		var expected []entity.TestCategory
		testCategoryRepo.On("GetAll", mock.Anything).Return(expected, errors.New("something went wrong"))
		srv := NewTestCategoryService(testCategoryRepo)

		items, err := srv.GetAll(context.Background())
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestCategoriesFromRepository)
		assert.Len(t, items, 0)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("error case - not found", func(t *testing.T) {
		testCategoryRepo := new(mocks.MockTestCategory)
		testCategoryRepo.On("GetByID", mock.Anything, uint64(1)).Return(nil, pkg.ErrTestCategoryNotFound)

		srv := NewTestCategoryService(testCategoryRepo)
		_, err := srv.GetByID(context.Background(), uint64(1))
		assert.ErrorIs(t, err, pkg.ErrTestCategoryNotFound)
	})
	t.Run("success case", func(t *testing.T) {
		testCategoryRepo := new(mocks.MockTestCategory)
		want := &entity.TestCategory{
			ID:                      1,
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
			Name:                    "load",
			Label:                   "Load Test",
			HasMaxTestServiceCount:  true,
			HasExecutionDuration:    true,
			HasAutoStepIncreaseRate: false,
		}
		testCategoryRepo.On("GetByID", mock.Anything, uint64(1)).Return(want, nil)

		srv := NewTestCategoryService(testCategoryRepo)
		got, err := srv.GetByID(context.Background(), uint64(1))
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(got, want))
	})
}
