package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/provider/dto/response"
	"control-panel-service/internal/provider/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTestAgentController(t *testing.T) {
	serviceHost := "challenge-test-service-serve"
	ingressHost := "127.0.0.1"
	ingressPort := 8081

	ctrl := NewTestAgentController(new(mocks.MockProvisioningService), new(mocks.MockTestServiceSDK), &entity.TestScenario{}, serviceHost, ingressHost, ingressPort)

	assert.NotNil(t, ctrl)

	a := ctrl


	assert.NotNil(t, a.provisioningService)
}

func Test_testAgentController_provisionTestService(t *testing.T) {
	t.Run("fail_case_provisioning_service_has_error", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provSvc.On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, new(mocks.MockTestServiceSDK), scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		err := c.provisionTestService(context.Background(), scenario, id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToProvisionTestService)

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})

	t.Run("fail_case_provisioning_service_has_error_make_sure_retries_work", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(3).
			Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, new(mocks.MockTestServiceSDK), scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 3
		c.provisioningRetriesSleep = time.Millisecond * 50

		start := time.Now()

		err := c.provisionTestService(context.Background(), scenario, id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToProvisionTestService)

		dur := time.Since(start)

		assert.GreaterOrEqual(t, dur, time.Millisecond*100) // 3 retries, sleep on first nad second one each one 50 ms

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})

	t.Run("success_case_provisioning_service_has_error_on_first_try_but_works_then", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(errors.New("something went wrong"))

		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(nil)

		ctrl := NewTestAgentController(provSvc, new(mocks.MockTestServiceSDK), scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 3
		c.provisioningRetriesSleep = time.Millisecond * 50

		start := time.Now()

		err := c.provisionTestService(context.Background(), scenario, id)
		assert.NoError(t, err)

		dur := time.Since(start)

		assert.GreaterOrEqual(t, dur, time.Millisecond*50) // 1 retries, each one 50 ms

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})
}

func Test_testAgentController_Run(t *testing.T) {
	t.Run("fail_case:_unable_to_provision_test_service", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, new(mocks.MockTestServiceSDK), scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50

		err := c.Run()
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToRunAgentControllerDueToProvisioningFailure)

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})
}

func Test_testAgentControler_StartTesting(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()
		zero := 0
		one := 1

		req := request.RunRequest{
			StepNum:               1,
			ExecutionId:           id,
			TestScenarioID:        5,
			MaxRequests:           1000,
			MaxDuration:           10000,
			RequestDelayDuration:  &zero,
			RandomRequestDelayMin: nil,
			RandomRequestDelayMax: nil,
			FixedTestNumber:       &one,
			RandomTestNumberMin:   nil,
			RandomTestNumberMax:   nil,
			BadValueRate:          0,
			RealValueRate:         0,
			NegativeValueRate:     0,
			ZeroValueRate:         0,
			StringValueRate:       0,
			LongStringValueRate:   0,
			NullValueRate:         0,
			DatabaseName:          "test-service",
			DatabaseTableName:     "executor",
		}

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("RunExecute", mock.Anything, baseUrl, req).
			Return(response.FactorialExecutionResult{}, nil)

		err := c.StartTesting(ctx, req)
		assert.NoError(t, err)
		provTest.AssertExpectations(t)
	})

	t.Run("fail_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()
		zero := 0
		one := 1

		req := request.RunRequest{
			StepNum:               1,
			ExecutionId:           id,
			TestScenarioID:        5,
			MaxRequests:           1000,
			MaxDuration:           10000,
			RequestDelayDuration:  &zero,
			RandomRequestDelayMin: nil,
			RandomRequestDelayMax: nil,
			FixedTestNumber:       &one,
			RandomTestNumberMin:   nil,
			RandomTestNumberMax:   nil,
			BadValueRate:          0,
			RealValueRate:         0,
			NegativeValueRate:     0,
			ZeroValueRate:         0,
			StringValueRate:       0,
			LongStringValueRate:   0,
			NullValueRate:         0,
			DatabaseName:          "test-service",
			DatabaseTableName:     "executor",
		}

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("RunExecute", mock.Anything, baseUrl, req).
			Return(response.FactorialExecutionResult{}, errors.New("run execute error"))

		err := c.StartTesting(ctx, req)
		assert.Error(t, err)
		provTest.AssertExpectations(t)
	})
}

func Test_testAgentControler_Healthy(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Health", mock.Anything, baseUrl).
			Return(response.HealthResponse{OK: true}, nil)

		healthy := c.Healthy()
		assert.True(t, healthy)
		provTest.AssertExpectations(t)
	})

	t.Run("fail_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Health", mock.Anything, baseUrl).
			Return(response.HealthResponse{}, errors.New("Health error"))

		healthy := c.Healthy()
		assert.False(t, healthy)
		provTest.AssertExpectations(t)
	})

	t.Run("success_case_is_not_healthy", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Health", mock.Anything, baseUrl).
			Return(response.HealthResponse{OK: false}, nil)

		healthy := c.Healthy()
		assert.False(t, healthy)
		provTest.AssertExpectations(t)
	})
}

func Test_testAgentControler_DeleteTesting(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		ctx := context.Background()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		provSvc.
			On("DeprovisionTestServiceByName", mock.Anything, c.scenario, c.uniqueID).
			Return(nil)

		err := c.DeleteTesting(ctx)
		assert.NoError(t, err)
		provSvc.AssertExpectations(t)
	})

	t.Run("fail_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		ctx := context.Background()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		provSvc.
			On("DeprovisionTestServiceByName", mock.Anything, c.scenario, c.uniqueID).
			Return(errors.New("faield to deprovision"))

		err := c.DeleteTesting(ctx)
		assert.Error(t, err)
		provSvc.AssertExpectations(t)
	})
}

func Test_testAgentControler_ReadyForTesting(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{
			ID:                5,
			TestServiceConfig: &entity.TestServiceConfig{MaxRequests: 1000},
		}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("ReadyForTest", mock.Anything, baseUrl).
			Return(response.HealthResponse{OK: true}, nil)

		ready := c.ReadyForTesting()
		assert.True(t, ready)
		provTest.AssertExpectations(t)
	})

	t.Run("fail_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("ReadyForTest", mock.Anything, baseUrl).
			Return(response.HealthResponse{}, errors.New("Health error"))

		ready := c.ReadyForTesting()
		assert.False(t, ready)
		provTest.AssertExpectations(t)
	})

	t.Run("success_case_is_not_healthy", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("ReadyForTest", mock.Anything, baseUrl).
			Return(response.HealthResponse{OK: false}, nil)

		ready := c.ReadyForTesting()
		assert.False(t, ready)
		provTest.AssertExpectations(t)
	})
}

func Test_testAgentControler_Pause(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Pause", mock.Anything, baseUrl).
			Return(&response.PauseResponse{
				Message: "done",
			}, nil)

		err := c.PauseTesting(ctx)
		assert.NoError(t, err)
		provTest.AssertExpectations(t)
	})

	t.Run("fail_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Pause", mock.Anything, baseUrl).
			Return(nil, errors.New("run execute error"))

		err := c.PauseTesting(ctx)
		assert.Error(t, err)
		provTest.AssertExpectations(t)
	})
}

func Test_testAgentControler_Resume(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Resume", mock.Anything, baseUrl).
			Return(&response.ResumeResponse{
				Message: "done",
			}, nil)

		err := c.ResumeTesting(ctx)
		assert.NoError(t, err)
		provTest.AssertExpectations(t)
	})

	t.Run("fail_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Resume", mock.Anything, baseUrl).
			Return(nil, errors.New("run execute error"))

		err := c.ResumeTesting(ctx)
		assert.Error(t, err)
		provTest.AssertExpectations(t)
	})
}

func Test_testAgentControler_Stop(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Stop", mock.Anything, baseUrl).
			Return(&response.StopResponse{
				Message: "done",
			}, nil)

		err := c.StopTesting(ctx)
		assert.NoError(t, err)
		provTest.AssertExpectations(t)
	})

	t.Run("failed_case", func(t *testing.T) {
		scenario := &entity.TestScenario{ID: 5}
		id := uuid.New()
		serviceHost := "challenge-test-service-serve"
		ingressHost := "127.0.0.1"
		ingressPort := 8081
		ctx := context.Background()

		provSvc := new(mocks.MockProvisioningService)
		provTest := new(mocks.MockTestServiceSDK)

		ctrl := NewTestAgentController(provSvc, provTest, scenario, serviceHost, ingressHost, ingressPort)

		c := ctrl
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50
		c.uniqueID = id

		baseUrl := c.urlGenerator()
		provTest.
			On("Stop", mock.Anything, baseUrl).
			Return(nil, errors.New("run execute error"))

		err := c.StopTesting(ctx)
		assert.Error(t, err)
		provTest.AssertExpectations(t)
	})
}
