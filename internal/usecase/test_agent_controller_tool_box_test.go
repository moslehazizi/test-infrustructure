package usecase_test

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/mocks"
	"control-panel-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestAgentControllerToolBox_Build(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		serviceHost := "challenge-test-service"
		ingressHost := "127.0.0.1"
		ingressPort := 8081

		b := usecase.
			NewTestAgentControllerToolBox(new(mocks.MockProvisioningService), new(mocks.MockTestServiceSDK), serviceHost, ingressHost, ingressPort).
			Build(&entity.TestScenario{})

		assert.NotNil(t, b)
	})

}
