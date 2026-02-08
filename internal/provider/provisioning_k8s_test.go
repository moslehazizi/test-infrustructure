package provider

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	kubermock "control-panel-service/pkg/kubernetes/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
		mockKubernetes := new(kubermock.KuberneteseMock)
		var motherService *entity.MotherService

		service := NewProvisioningService(cfg, mockKubernetes)

		err := service.ProvisionMotherService(ctx, motherService)
		assert.Error(t, err)
	})

	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

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

		err := service.ProvisionMotherService(ctx, motherService)

		assert.NoError(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed ApplyDeployment serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply deployment failed")).Once()

		err := service.ProvisionMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed ApplyService serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply service failed")).Once()

		err := service.ProvisionMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed WaitForDeployment serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{
			ID:                1,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "mother-service-serve-1", 10*time.Second).Return(errors.New("timeout waiting for deployment")).Once()

		err := service.ProvisionMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed ApplyDeployment jobs", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

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

		err := service.ProvisionMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed WaitForDeployment jobs", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

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

		err := service.ProvisionMotherService(ctx, motherService)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})
}

func deployTestScenarioServiceConfig() *config.Config {
	return &config.Config{
		Server: config.Server{Port: 8080},
		Kafka:  config.Kafka{Port: 9092},
		Postgres: config.Postgres{
			Port:               5432,
			SSLMode:            "disable",
			MaxOpenConnections: 10,
			MaxIdleConnections: 5,
		},
		// Kubernetese: config.Kubernetese{
		// 	MotherServiceAPPServe:          "mother-service-serve",
		// 	MotherServiceAPPJobs:           "mother-service-jobs",
		// 	MotherServiceAPPServeWaitReady: 10 * time.Second,
		// 	MotherServiceAPPJobsWaitReady:  10 * time.Second,
		// },
	}
}
