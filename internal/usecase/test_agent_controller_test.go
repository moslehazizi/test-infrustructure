package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
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
	testSvcServe := "challenge-test-service-serve"

	ctrl := NewTestAgentController(new(mocks.MockProvisioningService), &entity.TestScenario{}, testSvcServe)

	assert.NotNil(t, ctrl)

	a, ok := ctrl.(*testAgentController)
	assert.True(t, ok)

	assert.NotNil(t, a.provisioningService)
}

func Test_testAgentController_provisionTestService(t *testing.T) {
	t.Run("fail case: provisioning service has error", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()
		testSvcServe := "challenge-test-service-serve"

		provSvc := new(mocks.MockProvisioningService)
		provSvc.On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, scenario, testSvcServe)

		c := ctrl.(*testAgentController)
		err := c.provisionTestService(context.Background(), scenario, id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToProvisionTestService)

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})

	t.Run("fail case: provisioning service has error => make sure retries work", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()
		testSvcServe := "challenge-test-service-serve"

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(3).
			Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, scenario, testSvcServe)

		c := ctrl.(*testAgentController)
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

	t.Run("success case: provisioning service has error on first try but works then", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()
		testSvcServe := "challenge-test-service-serve"

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(errors.New("something went wrong"))

		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(nil)

		ctrl := NewTestAgentController(provSvc, scenario, testSvcServe)

		c := ctrl.(*testAgentController)
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
	t.Run("fail case: unable to provision test service", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		testSvcServe := "challenge-test-service-serve"

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, scenario, testSvcServe)

		c := ctrl.(*testAgentController)
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50

		err := c.Run()
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToRunAgentControllerDueToProvisioningFailure)

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})
}

func Test_testAgentControler_StartTesting(t *testing.T) {
	t.Run("ssuccess case", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()
		testSvcServe := "challenge-test-service-serve"
		ctx := context.Background()

		provSvc := new(mocks.MockProvisioningService)
		ctrl := NewTestAgentController(provSvc, scenario, testSvcServe)

		c := ctrl.(*testAgentController)
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50

		go func() {
			provSvc.
				On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
				Times(1).
				Return(nil)

			err := c.Run()
			assert.NoError(t, err)
		}()

		err := c.StartTesting(ctx, request.RunRequest{
			StepNum:            1,
			ExecutionId:        id,
			ScenarioId:         2,
			StepIncrement:      10,
			MaxTxsCount:        1000,
			MaxTxsDuration:     10000,
			MaxDelayBetweenTxs: 0,
			MinDelayBetweenTxs: 0,
			MaxInputNum:        1,
			MinInputNum:        1,
			TotalErr:           0,
			RealNumErr:         0,
			NegativeNumErr:     0,
			ZeroNumErr:         0,
			ShortStrErr:        0,
			LongStrErr:         0,
			NilErr:             0,
		})

		assert.NoError(t, err)
	})
}

func Test_testAgentControler_startTesting(t *testing.T) {
	t.Run("success response", func(t *testing.T) {

	})
}
