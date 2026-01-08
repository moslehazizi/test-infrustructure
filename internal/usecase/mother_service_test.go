package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMotherService(t *testing.T) {
	mockRepo := new(mocks.MockMotherService)
	service := NewMotherService(mockRepo)

	assert.NotNil(t, service)
}

func TestMotherServiceUsecase_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		sampleMS := &entity.MotherService{
			Name:                "mother1",
			ProvisioningStatus:  entity.ProvisioningStatusFailed,
			DatabaseName:        "db1",
			DatabaseTableName:   "factorial",
			KafkaLiveFeedTopic:  "live_feed",
			KafkaFactorialTopic: "factorial",
		}

		mockRepo.On("Create", ctx, sampleMS).Return(nil)

		err := service.Create(ctx, sampleMS)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", ctx, sampleMS)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		sampleMS := &entity.MotherService{
			Name:                "mother1",
			ProvisioningStatus:  entity.ProvisioningStatusFailed,
			DatabaseName:        "db1",
			DatabaseTableName:   "factorial",
			KafkaLiveFeedTopic:  "live_feed",
			KafkaFactorialTopic: "factorial",
		}

		mockRepo.On("Create", ctx, sampleMS).Return(pkg.ErrFailedToCreateMotherService)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateMotherService)
		mockRepo.AssertCalled(t, "Create", ctx, sampleMS)
	})

	t.Run("failed case - duplicate", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		sampleMS := &entity.MotherService{
			Name:                "mother1",
			ProvisioningStatus:  entity.ProvisioningStatusFailed,
			DatabaseName:        "db1",
			DatabaseTableName:   "factorial",
			KafkaLiveFeedTopic:  "live_feed",
			KafkaFactorialTopic: "factorial",
		}

		mockRepo.On("Create", ctx, sampleMS).Return(pkg.ErrMotherServiceAlreadyExist)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceAlreadyExist)
		mockRepo.AssertCalled(t, "Create", ctx, sampleMS)
	})
}

func TestMotherServiceUsecase_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		inputID := uint64(1)
		expectedResult := &entity.MotherService{
			ID:                  uint64(1),
			Name:                "mother1",
			ProvisioningStatus:  entity.ProvisioningStatusFailed,
			DatabaseName:        "db1",
			DatabaseTableName:   "factorial",
			KafkaLiveFeedTopic:  "live_feed",
			KafkaFactorialTopic: "factorial",
		}

		mockRepo.On("GetByID", ctx, inputID).Return(expectedResult, nil)

		result, err := service.GetByID(ctx, inputID)

		assert.Nil(t, err)
		assert.Equal(t, result.ID, expectedResult.ID)
		assert.Equal(t, result.Name, expectedResult.Name)
		assert.Equal(t, result.ProvisioningStatus, expectedResult.ProvisioningStatus)
		assert.Equal(t, result.DatabaseName, expectedResult.DatabaseName)
		assert.Equal(t, result.DatabaseTableName, expectedResult.DatabaseTableName)
		assert.Equal(t, result.KafkaLiveFeedTopic, expectedResult.KafkaLiveFeedTopic)
		assert.Equal(t, result.KafkaFactorialTopic, expectedResult.KafkaFactorialTopic)
		mockRepo.AssertCalled(t, "GetByID", ctx, inputID)
	})

	t.Run("failed case - not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		inputID := uint64(1)

		mockRepo.On("GetByID", ctx, inputID).Return(nil, pkg.ErrMotherServiceNotFound)

		result, err := service.GetByID(ctx, inputID)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceNotFound)
		mockRepo.AssertCalled(t, "GetByID", ctx, inputID)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		inputID := uint64(1)

		mockRepo.On("GetByID", ctx, inputID).Return(nil, errors.New("failed to get mother service"))

		result, err := service.GetByID(ctx, inputID)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherService)
		mockRepo.AssertCalled(t, "GetByID", ctx, inputID)
	})
}

func TestMotherServiceUsecase_GetPaginated(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		serviceAddress1 := "http://service1.example.com"
		serviceAddress2 := "http://service2.example.com"

		expectedMotherServices := []*entity.MotherService{
			{
				ID:                       uint64(5),
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
				Name:                     "mother5",
				ExceptionRate:            0.0,
				ResponseDelayRate:        0.0,
				ProvisioningStatus:       entity.ProvisioningStatusPending,
				ServiceDeploymentAddress: &serviceAddress2,
				DatabaseName:             "test_db5",
				DatabaseTableName:        "test_table5",
				KafkaLiveFeedTopic:       "live_feed",
				KafkaFactorialTopic:      "factorial",
			},
			{
				ID:                       uint64(4),
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
				Name:                     "mother4",
				ExceptionRate:            0.1,
				ResponseDelayRate:        0.2,
				ProvisioningStatus:       entity.ProvisioningStatusProvisioned,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db4",
				DatabaseTableName:        "test_table4",
				KafkaLiveFeedTopic:       "live_feed",
				KafkaFactorialTopic:      "factorial",
			},
		}

		mockRepo.On("GetPaginated", ctx, paginationRequest).Return(expectedMotherServices, nil)

		result, err := service.GetPaginated(ctx, paginationRequest)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMotherServices, result)
		mockRepo.AssertCalled(t, "GetPaginated", ctx, paginationRequest)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", ctx, paginationRequest).Return(nil, errors.New("failed to get mother services"))

		result, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", ctx, paginationRequest)
	})

	t.Run("failed case - negative page", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		paginationRequest := entity.PaginationRequest{
			Page:    -1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", ctx, paginationRequest).Return(nil, errors.New("failed to get mother services"))

		result, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", ctx, paginationRequest)
	})
}
