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
		Server: config.Server{
			Port:          8080,
			SwaggerScheme: []string{"https"},
		},
		Kafka: config.Kafka{Port: 9092},
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
		Server: config.Server{
			Port:          8080,
			SwaggerScheme: []string{"https"},
		},
		Kafka: config.Kafka{Port: 9092},
		Postgres: config.Postgres{
			Port:               5432,
			SSLMode:            "disable",
			MaxOpenConnections: 10,
			MaxIdleConnections: 5,
		},
		Kubernetese: config.Kubernetese{
			TestServiceAPPServe:          "test-service-serve",
			TestServiceAPPJobs:           "test-service-jobs",
			TestServiceAPPServeWaitReady: 10 * time.Second,
			TestServiceAPPJobsWaitReady:  10 * time.Second,
		},
	}
}

func TestDeployTestScenarioService(t *testing.T) {
	t.Run("test_scenarios_service_is_nil", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		var testScenarioService *entity.TestScenario
		replica := int32(10)

		mockKubernetes := new(kubermock.KuberneteseMock)

		service := NewProvisioningService(cfg, mockKubernetes)

		err := service.ProvisionTestService(ctx, testScenarioService, replica)
		assert.Error(t, err)
	})

	t.Run("test_service_config_is_nil", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{}
		replica := int32(10)
		mockKubernetes := new(kubermock.KuberneteseMock)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario",
		}

		service := NewProvisioningService(cfg, mockKubernetes)

		err := service.ProvisionTestService(ctx, testScenarioService, replica)
		assert.Error(t, err)
	})

	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployTestScenarioServiceConfig()
		replica := int32(10)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "test-service-serve-1", 10*time.Second).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "test-service-jobs", 10*time.Second).Return(nil).Once()

		err := service.ProvisionTestService(ctx, testScenarioService, replica)

		assert.NoError(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_ApplyDeployment_serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployTestScenarioServiceConfig()
		replica := int32(10)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply deployment failed")).Once()

		err := service.ProvisionTestService(ctx, testScenarioService, replica)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_ApplyService_serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployTestScenarioServiceConfig()
		replica := int32(10)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply service failed")).Once()

		err := service.ProvisionTestService(ctx, testScenarioService, replica)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_WaitForDeployment_serve", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployTestScenarioServiceConfig()
		replica := int32(10)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "test-service-serve-1", 10*time.Second).Return(errors.New("timeout waiting for deployment")).Once()

		err := service.ProvisionTestService(ctx, testScenarioService, replica)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_ApplyDeployment_jobs", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployTestScenarioServiceConfig()
		replica := int32(10)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // serve
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "test-service-serve-1", 10*time.Second).Return(nil).Once()
		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("apply jobs deployment failed")).Once() // jobs

		err := service.ProvisionTestService(ctx, testScenarioService, replica)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_WaitForDeployment_jobs", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployTestScenarioServiceConfig()
		replica := int32(10)
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		mockKubernetes.On("ApplyDeployment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
		mockKubernetes.On("ApplyService", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "test-service-serve-1", 10*time.Second).Return(nil).Once()
		mockKubernetes.On("WaitForDeployment", mock.Anything, "test-service-jobs", 10*time.Second).Return(errors.New("timeout waiting for jobs")).Once()

		err := service.ProvisionTestService(ctx, testScenarioService, replica)

		assert.Error(t, err)
		mockKubernetes.AssertExpectations(t)
	})

}

func TestDeprovisionMotherService(t *testing.T) {
	t.Run("mother_service_is_nil", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		err := service.DeprovisionMotherService(ctx, nil)
		assert.Error(t, err)
		mockKubernetes.AssertNotCalled(t, "DeleteDeployment")
	})

	t.Run("success_deletes_serve_deployment_service_config_secret", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{ID: 1, Name: "mother1"}
		serveName := "mother-service-serve-1"
		configName := serveName + "-config"
		secretName := serveName + "-secret"

		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteConfigMap", mock.Anything, configName).Return(nil).Once()
		mockKubernetes.On("DeleteSecret", mock.Anything, secretName).Return(nil).Once()

		err := service.DeprovisionMotherService(ctx, motherService)

		assert.NoError(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_DeleteDeployment", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{ID: 1, Name: "mother1"}
		serveName := "mother-service-serve-1"

		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(errors.New("delete deployment failed")).Once()

		err := service.DeprovisionMotherService(ctx, motherService)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_DeleteService", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{ID: 1, Name: "mother1"}
		serveName := "mother-service-serve-1"

		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(errors.New("delete service failed")).Once()

		err := service.DeprovisionMotherService(ctx, motherService)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_DeleteConfigMap", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{ID: 1, Name: "mother1"}
		serveName := "mother-service-serve-1"
		configName := serveName + "-config"

		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteConfigMap", mock.Anything, configName).Return(errors.New("delete service failed")).Once()

		err := service.DeprovisionMotherService(ctx, motherService)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("failed_DeleteSecret", func(t *testing.T) {
		ctx := context.Background()
		cfg := deployMotherServiceConfig()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		motherService := &entity.MotherService{ID: 1, Name: "mother1"}
		serveName := "mother-service-serve-1"
		configName := serveName + "-config"
		secretName := serveName + "-secret"

		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteConfigMap", mock.Anything, configName).Return(nil).Once()
		mockKubernetes.On("DeleteSecret", mock.Anything, secretName).Return(errors.New("delete service failed")).Once()

		err := service.DeprovisionMotherService(ctx, motherService)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})
}

func TestDeprovisionTestService(t *testing.T) {
	cfg := deployTestScenarioServiceConfig()
	cfg.Kubernetese = config.Kubernetese{
		TestServiceAPPServe: "test-service-serve",
	}

	t.Run("test_scenario_is_nil", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		err := service.DeprovisionTestService(ctx, nil, 2)
		assert.Error(t, err)
	})

	t.Run("test_service_config_is_nil", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenarioService := &entity.TestScenario{
			ID:   1,
			Name: "testScenario",
		}

		err := service.DeprovisionTestService(ctx, testScenarioService, 2)
		assert.Error(t, err)
		mockKubernetes.AssertNotCalled(t, "GetDeploymentReplicas")
	})

	t.Run("scale_down_partial_replicas_keeps_deployment", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              5,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(2)
		serveName := "test-service-serve-5"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(5, nil).Once()
		mockKubernetes.On("ScaleDeployment", mock.Anything, serveName, int32(3)).Return(nil).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NoError(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("remove_all_replicas_deletes_deployment_service_config_secret", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              3,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(4)
		serveName := "test-service-serve-3"
		configName := serveName + "-config"
		secretName := serveName + "-secret"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(4, nil).Once()
		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteConfigMap", mock.Anything, configName).Return(nil).Once()
		mockKubernetes.On("DeleteSecret", mock.Anything, secretName).Return(nil).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NoError(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("remove_more_than_current_replicas_deletes_all_resources", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              2,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(10)
		serveName := "test-service-serve-2"
		configName := serveName + "-config"
		secretName := serveName + "-secret"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(3, nil).Once()
		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteConfigMap", mock.Anything, configName).Return(nil).Once()
		mockKubernetes.On("DeleteSecret", mock.Anything, secretName).Return(nil).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NoError(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("fail_GetDeploymentReplicas", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              2,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(5)
		serveName := "test-service-serve-2"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(0, errors.New("GetDeploymentReplicas failed")).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("fail_DeleteDeployment", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              2,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(5)
		serveName := "test-service-serve-2"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(5, nil).Once()
		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(errors.New("DeleteDeployment failed")).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("fail_ScaleDeployment", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              2,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(5)
		serveName := "test-service-serve-2"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(10, nil).Once()
		mockKubernetes.On("ScaleDeployment", mock.Anything, serveName, replicaToRemove).Return(errors.New("ScaleDeployment failed")).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("fail_DeleteService", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              2,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(5)
		serveName := "test-service-serve-2"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(5, nil).Once()
		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(errors.New("DeleteService failed")).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("fail_DeleteConfigMap", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              2,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(5)
		serveName := "test-service-serve-2"
		configName := serveName + "-config"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(5, nil).Once()
		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteConfigMap", mock.Anything, configName).Return(errors.New("DeleteConfigMap failed")).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})

	t.Run("fail_DeleteConfigMap", func(t *testing.T) {
		ctx := context.Background()
		mockKubernetes := new(kubermock.KuberneteseMock)
		service := NewProvisioningService(cfg, mockKubernetes)

		testScenario := &entity.TestScenario{
			ID:              2,
			MotherServiceID: 1,
			Name:            "testScenario1",
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                  2,
				TestScenarioID:      1,
				MaxRequests:         100,
				MaxDuration:         1000,
				BadValueRate:        10,
				NegativeValueRate:   20,
				RealValueRate:       20,
				ZeroValueRate:       20,
				StringValueRate:     20,
				LongStringValueRate: 10,
				NullValueRate:       10,
			},
		}

		replicaToRemove := int32(5)
		serveName := "test-service-serve-2"
		configName := serveName + "-config"
		secretName := serveName + "-secret"

		mockKubernetes.On("GetDeploymentReplicas", mock.Anything, serveName).Return(5, nil).Once()
		mockKubernetes.On("DeleteDeployment", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteService", mock.Anything, serveName).Return(nil).Once()
		mockKubernetes.On("DeleteConfigMap", mock.Anything, configName).Return(nil).Once()
		mockKubernetes.On("DeleteSecret", mock.Anything, secretName).Return(errors.New("DeleteSecret failed")).Once()

		err := service.DeprovisionTestService(ctx, testScenario, replicaToRemove)

		assert.NotNil(t, err)
		mockKubernetes.AssertExpectations(t)
	})
}
