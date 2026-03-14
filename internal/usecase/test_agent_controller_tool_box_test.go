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
		testSvcServe := "challenge-test-service"
		testSvcPort := 8080

		b := usecase.
			NewTestAgentControllerToolBox(new(mocks.MockProvisioningService), new(mocks.MockTestServiceSDK), testSvcServe, testSvcPort).
			Build(&entity.TestScenario{})

		assert.NotNil(t, b)
	})

}
