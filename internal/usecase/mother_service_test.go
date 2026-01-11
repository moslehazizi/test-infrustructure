package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	providerMock "control-panel-service/internal/provider/mocks"
	"control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMotherService(t *testing.T) {
	mockRepo := new(mocks.MockMotherService)
	mockEventProducer := new(providerMock.KafkaMock)
	service := NewMotherService(mockRepo, mockEventProducer)

	assert.NotNil(t, service)

	s, ok := service.(*motherService)
	assert.True(t, ok)
	assert.NotNil(t, s.motherServiceRepo)
	assert.NotNil(t, s.eventProducer)
}

func TestMotherServiceUsecase_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Begin").Return(mock.Anything)
		mockRepo.On("Create", ctx, sampleMS).Return(nil)
		mockRepo.On("Commit").Return(nil)

		mockEventProducer.On("SendEvent", mock.Anything, mock.Anything, "provisioning").Return(nil)

		err := service.Create(ctx, sampleMS)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", ctx, sampleMS)
		mockRepo.AssertCalled(t, "Begin")
		mockRepo.AssertCalled(t, "Commit")
		mockEventProducer.AssertCalled(t, "SendEvent", mock.Anything, mock.Anything, "provisioning")
	})
	t.Run("failed to send kafka event => database should be rolled back", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Begin").Return(mock.Anything)
		mockRepo.On("Create", ctx, sampleMS).Return(nil)
		mockRepo.On("Rollback").Return(nil)

		mockEventProducer.On("SendEvent", mock.Anything, mock.Anything, "provisioning").Return(errors.New("something went wrong"))

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSendProvisioningEvent)

		mockRepo.AssertCalled(t, "Create", ctx, sampleMS)
		mockRepo.AssertCalled(t, "Begin")
		mockRepo.AssertCalled(t, "Rollback")

		mockEventProducer.AssertCalled(t, "SendEvent", mock.Anything, mock.Anything, "provisioning")
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Begin").Return(mock.Anything)
		mockRepo.On("Create", ctx, sampleMS).Return(pkg.ErrFailedToCreateMotherService)
		mockRepo.On("Rollback").Return(nil)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateMotherService)
		mockRepo.AssertCalled(t, "Create", ctx, sampleMS)
		mockRepo.AssertCalled(t, "Begin")
		mockRepo.AssertCalled(t, "Rollback")
	})

	t.Run("failed case - duplicate", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Begin").Return(mock.Anything)
		mockRepo.On("Create", ctx, sampleMS).Return(pkg.ErrMotherServiceAlreadyExist)
		mockRepo.On("Rollback").Return(nil)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceAlreadyExist)
		mockRepo.AssertCalled(t, "Create", ctx, sampleMS)
		mockRepo.AssertCalled(t, "Begin")
		mockRepo.AssertCalled(t, "Rollback")
	})

	t.Run("failed case - validation error service name is missing", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

		sampleMS := &entity.MotherService{
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidName)
	})

	t.Run("failed case - validation error response delay rete not be negative", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

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
	t.Run("failed case - validation error exception rate is negative", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

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

	t.Run("failed case - validation error fixed delay is set but rate is 0", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)
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
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

		inputID := uint64(1)
		expectedResult := &entity.MotherService{
			ID:                 uint64(1),
			Name:               "mother1",
			ProvisioningStatus: entity.ProvisioningStatusFailed,
			DatabaseName:       "db1",
			DatabaseTableName:  "factorial",
		}

		mockRepo.On("GetByID", ctx, inputID).Return(expectedResult, nil)

		result, err := service.GetByID(ctx, inputID)

		assert.Nil(t, err)
		assert.Equal(t, result.ID, expectedResult.ID)
		assert.Equal(t, result.Name, expectedResult.Name)
		assert.Equal(t, result.ProvisioningStatus, expectedResult.ProvisioningStatus)
		assert.Equal(t, result.DatabaseName, expectedResult.DatabaseName)
		assert.Equal(t, result.DatabaseTableName, expectedResult.DatabaseTableName)
		mockRepo.AssertCalled(t, "GetByID", ctx, inputID)
	})

	t.Run("failed case - not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

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
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

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
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

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
			},
			{
				ID:                       uint64(4),
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
				Name:                     "mother4",
				ExceptionRate:            10,
				ResponseDelayRate:        20,
				ProvisioningStatus:       entity.ProvisioningStatusProvisioned,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db4",
				DatabaseTableName:        "test_table4",
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
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

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
		mockEventProducer := new(providerMock.KafkaMock)
		service := NewMotherService(mockRepo, mockEventProducer)

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
