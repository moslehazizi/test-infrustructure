package usecase

import (
	"context"
	"control-panel-service/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDatabaseMetadataUsecase_initialization(t *testing.T) {
	databaseMetadataRepo := new(mocks.MockDatabaseMetadata)
	srv := NewDatabaseMetadata(databaseMetadataRepo)

	s, ok := srv.(*databaseMetadata)
	assert.True(t, ok)
	assert.NotNil(t, s.repo)
}

func TestDatabaseMetadataUsecase_GetAll(t *testing.T) {
	t.Run("success case - with result", func(t *testing.T) {
		repo := new(mocks.MockDatabaseMetadata)
		repo.On("GetAll", mock.Anything).
			Return([]string{"postgres", "load_test_db"}, nil)

		uc := NewDatabaseMetadata(repo)
		result, err := uc.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, []string{"postgres", "load_test_db"}, result)
		repo.AssertExpectations(t)
	})

	t.Run("success case - with empty result", func(t *testing.T) {
		repo := new(mocks.MockDatabaseMetadata)
		repo.On("GetAll", mock.Anything).
			Return([]string(nil), nil)

		uc := NewDatabaseMetadata(repo)
		result, err := uc.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		repo.AssertExpectations(t)
	})

	t.Run("failure case - repository error", func(t *testing.T) {
		repo := new(mocks.MockDatabaseMetadata)
		repo.On("GetAll", mock.Anything).
			Return(nil, errors.New("db error"))

		uc := NewDatabaseMetadata(repo)
		result, err := uc.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestStorageUsecase_GetTablesByDBName(t *testing.T) {
	t.Run("success case - with result", func(t *testing.T) {
		repo := new(mocks.MockDatabaseMetadata)
		repo.On("GetTablesByDBName", mock.Anything, "load_test_db").
			Return([]string{"events", "logs"}, nil)

		uc := NewDatabaseMetadata(repo)
		result, err := uc.GetTablesByDBName(context.Background(), "load_test_db")

		assert.NoError(t, err)
		assert.Equal(t, []string{"events", "logs"}, result)
		repo.AssertExpectations(t)
	})

	t.Run("success case - with empty result", func(t *testing.T) {
		repo := new(mocks.MockDatabaseMetadata)
		repo.On("GetTablesByDBName", mock.Anything, "load_test_db").
			Return([]string(nil), nil)

		uc := NewDatabaseMetadata(repo)
		result, err := uc.GetTablesByDBName(context.Background(), "load_test_db")

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		repo.AssertExpectations(t)
	})

	t.Run("failure case - bad request", func(t *testing.T) {
		repo := new(mocks.MockDatabaseMetadata)

		uc := NewDatabaseMetadata(repo)
		result, err := uc.GetTablesByDBName(context.Background(), "")

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("failure case - repository error", func(t *testing.T) {
		repo := new(mocks.MockDatabaseMetadata)
		repo.On("GetTablesByDBName", mock.Anything, "load_test_db").
			Return(nil, errors.New("db error"))

		uc := NewDatabaseMetadata(repo)
		result, err := uc.GetTablesByDBName(context.Background(), "load_test_db")

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}
