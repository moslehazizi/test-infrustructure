package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
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
	ctrl := NewTestAgentController(new(mocks.MockProvisioningService), &entity.TestScenario{})

	assert.NotNil(t, ctrl)

	a, ok := ctrl.(*testAgentController)
	assert.True(t, ok)

	assert.NotNil(t, a.provisioningService)
}

func Test_testAgentController_provisionTestService(t *testing.T) {
	t.Run("fail case: provisioning service has error", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()

		provSvc := new(mocks.MockProvisioningService)
		provSvc.On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, scenario)

		c := ctrl.(*testAgentController)
		err := c.provisionTestService(context.Background(), scenario, id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToProvisionTestService)

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})

	t.Run("fail case: provisioning service has error => make sure retries work", func(t *testing.T) {
		scenario := &entity.TestScenario{}
		id := uuid.New()

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(3).
			Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, scenario)

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

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(errors.New("something went wrong"))

		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(nil)

		ctrl := NewTestAgentController(provSvc, scenario)

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

		provSvc := new(mocks.MockProvisioningService)
		provSvc.
			On("ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything).
			Times(1).
			Return(errors.New("something went wrong"))

		ctrl := NewTestAgentController(provSvc, scenario)

		c := ctrl.(*testAgentController)
		c.provisioningRetries = 1
		c.provisioningRetriesSleep = time.Millisecond * 50

		err := c.Run()
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToRunAgentControllerDueToProvisioningFailure)

		provSvc.AssertCalled(t, "ProvisionTestServiceByName", mock.Anything, scenario, mock.Anything)
	})
}
