package usecase

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	inEntity "control-panel-service/internal/domain/entity"
	providerMock "control-panel-service/internal/provider/mocks"
	"control-panel-service/internal/repository/mocks"

	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	connmock "control-panel-service/pkg/database/postgres/mocks"
	kubermock "control-panel-service/pkg/kubernetes/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getMockDB(t *testing.T) database.Database {
	mockConn := new(connmock.Connection)
	db, _, err := mockConn.OpenConnection()
	require.NoError(t, err)

	return db
}

func TestNewMotherService(t *testing.T) {
	cfg := &config.Config{}
	mockRepo := new(mocks.MockMotherService)
	mockEventProducer := new(providerMock.KafkaMock)
	mockKubernetes := new(kubermock.KuberneteseMock)
	service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

	assert.NotNil(t, service)

	s, ok := service.(*motherService)
	assert.True(t, ok)
	assert.NotNil(t, s.motherServiceRepo)
	assert.NotNil(t, s.eventProducer)
}

func TestMotherServiceUsecase_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(nil)

		mockEventProducer.On("SendEvent", mock.Anything, mock.Anything, "provisioning").Return(nil)

		err := service.Create(ctx, sampleMS)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
		mockEventProducer.AssertCalled(t, "SendEvent", mock.Anything, mock.Anything, "provisioning")
	})
	t.Run("failed to send kafka event => database should be rolled back", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(nil)

		mockEventProducer.On("SendEvent", mock.Anything, mock.Anything, "provisioning").Return(errors.New("something went wrong"))

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToSendProvisioningEvent)

		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)

		mockEventProducer.AssertCalled(t, "SendEvent", mock.Anything, mock.Anything, "provisioning")
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(pkg.ErrFailedToCreateMotherService)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateMotherService)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
	})

	t.Run("failed case - duplicate", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(pkg.ErrMotherServiceAlreadyExist)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceAlreadyExist)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
	})

	t.Run("failed case - validation error service name is missing", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		sampleMS := &entity.MotherService{
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMotherServiceName)
	})

	t.Run("failed case - validation error response delay rete not be negative", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

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
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

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
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)
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
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		inputID := uint64(1)
		expectedResult := &entity.MotherService{
			ID:                uint64(1),
			Name:              "mother1",
			Status:            entity.MotherServiceStatusRunning,
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("GetByID", mock.Anything, inputID).Return(expectedResult, nil)

		result, err := service.GetByID(ctx, inputID)

		assert.Nil(t, err)
		assert.Equal(t, result.ID, expectedResult.ID)
		assert.Equal(t, result.Name, expectedResult.Name)
		assert.Equal(t, result.Status, expectedResult.Status)
		assert.Equal(t, result.DatabaseName, expectedResult.DatabaseName)
		assert.Equal(t, result.DatabaseTableName, expectedResult.DatabaseTableName)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, inputID)
	})

	t.Run("failed case - not found", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		inputID := uint64(1)

		mockRepo.On("GetByID", mock.Anything, inputID).Return(nil, pkg.ErrMotherServiceNotFound)

		result, err := service.GetByID(ctx, inputID)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, inputID)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		inputID := uint64(1)

		mockRepo.On("GetByID", mock.Anything, inputID).Return(nil, errors.New("failed to get mother service"))

		result, err := service.GetByID(ctx, inputID)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherService)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, inputID)
	})
}

func TestMotherServiceUsecase_GetPaginated(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

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
				Status:                   entity.MotherServiceStatusPending,
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
				Status:                   entity.MotherServiceStatusRunning,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db4",
				DatabaseTableName:        "test_table4",
			},
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(expectedMotherServices, nil)

		result, err := service.GetPaginated(ctx, paginationRequest)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMotherServices, result)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(nil, errors.New("failed to get mother services"))

		result, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})

	t.Run("failed case - negative page", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		paginationRequest := entity.PaginationRequest{
			Page:    -1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(nil, errors.New("failed to get mother services"))

		result, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})
}

func TestDeployMotherService(t *testing.T) {
	t.Run("mother_service_is_nil", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockEventProducer := new(providerMock.KafkaMock)
		mockKubernetes := new(kubermock.KuberneteseMock)
		var motherService *inEntity.MotherService

		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockEventProducer, mockKubernetes)

		err := service.DeployMotherService(ctx, motherService)
		assert.NotNil(t, err)
	})

}
