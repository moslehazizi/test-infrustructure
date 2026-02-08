package usecase

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
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

	mockKubernetes := new(kubermock.KuberneteseMock)
	service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

	assert.NotNil(t, service)

	s, ok := service.(*motherService)
	assert.True(t, ok)
	assert.NotNil(t, s.motherServiceRepo)
}

func TestMotherServiceUsecase_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{
			Kubernetese: config.Kubernetese{
				MotherServiceAPPServe:          "mother-service-serve",
				MotherServiceAPPJobs:           "mother-service-jobs",
				MotherServiceAPPServeWaitReady: 10 * time.Second,
				MotherServiceAPPJobsWaitReady:  10 * time.Second,
			},
		}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(1), nil)
		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-serve-1", 10*time.Second).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-jobs", 10*time.Second).Return(nil).Once()

		err := service.Create(ctx, sampleMS)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
		mockKubernetes.AssertExpectations(t)
	})
	t.Run("failed DeployMotherService => returns ErrFailedToDeployMotherService", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{
			Kubernetese: config.Kubernetese{
				MotherServiceAPPServe:          "mother-service-serve",
				MotherServiceAPPJobs:           "mother-service-jobs",
				MotherServiceAPPServeWaitReady: 10 * time.Second,
				MotherServiceAPPJobsWaitReady:  10 * time.Second,
			},
		}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(1), nil)
		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply deployment failed")).Once()

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToDeployMotherService)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(0), pkg.ErrFailedToCreateMotherService)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateMotherService)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
	})

	t.Run("failed case - duplicate", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		sampleMS := &entity.MotherService{
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockRepo.On("Create", mock.Anything, sampleMS).Return(uint64(0), pkg.ErrMotherServiceAlreadyExist)

		err := service.Create(ctx, sampleMS)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceAlreadyExist)
		mockRepo.AssertCalled(t, "Create", mock.Anything, sampleMS)
	})

	t.Run("failed case - validation error service name is missing", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

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
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

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
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

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
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)
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
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

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
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

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
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

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
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		serviceAddress1 := "http://service1.example.com"
		serviceAddress2 := "http://service2.example.com"
		count := int64(2)

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

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(expectedMotherServices, count, nil)

		result, total, err := service.GetPaginated(ctx, paginationRequest)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, count)
		assert.Equal(t, expectedMotherServices, result)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})

	t.Run("failed case", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(nil, int64(0), errors.New("failed to get mother services"))

		result, count, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})

	t.Run("failed case - negative page", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		paginationRequest := entity.PaginationRequest{
			Page:    -1,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, paginationRequest).Return(nil, int64(0), errors.New("failed to get mother services"))

		result, count, err := service.GetPaginated(ctx, paginationRequest)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))
		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherServices)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, paginationRequest)
	})
}

func deployMotherServiceConfig() *config.Config {
	return &config.Config{
		Server: config.Server{Port: 8080},
		Kafka:  config.Kafka{Port: 9092},
		Postgres: config.Postgres{
			Port:               5432,
			SSLMode:            "disable",
			MaxOpenConnections: 10,
			MaxIdleConnections: 5,
		},
		Kubernetese: config.Kubernetese{
			MotherServiceAPPServe:          "mother-service-serve",
			MotherServiceAPPJobs:           "mother-service-jobs",
			MotherServiceAPPServeWaitReady: 10 * time.Second,
			MotherServiceAPPJobsWaitReady:  10 * time.Second,
		},
	}
}

func TestDeployMotherService(t *testing.T) {
	t.Run("mother_service_is_nil", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		var motherService *entity.MotherService

		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		err := service.DeployMotherService(ctx, motherService)
		assert.Error(t, err)
	})

	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-serve-1", 10*time.Second).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-jobs", 10*time.Second).Return(nil).Once()

		err := service.DeployMotherService(ctx, motherService)

		assert.NoError(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed ApplyDeployment serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply deployment failed")).Once()

		err := service.DeployMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed ApplyService serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply service failed")).Once()

		err := service.DeployMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed WaitForDeployment serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-serve-1", 10*time.Second).Return(errors.New("timeout waiting for deployment")).Once()

		err := service.DeployMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed ApplyDeployment jobs", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // serve
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-serve-1", 10*time.Second).Return(nil).Once()
		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply jobs deployment failed")).Once() // jobs

		err := service.DeployMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed WaitForDeployment jobs", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockRepo := new(mocks.MockMotherService)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewMotherService(cfg, getMockDB(t), mockRepo, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-serve-1", 10*time.Second).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-jobs", 10*time.Second).Return(errors.New("timeout waiting for jobs")).Once()

		err := service.DeployMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})
}
