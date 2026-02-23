package usecase_test

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/mocks"
	"control-panel-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestAgentControllerBuilder_Build(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		testSvcServe := "challenge-test-service"
		b := usecase.
			NewTestAgentControllerBuilder(new(mocks.MockProvisioningService), testSvcServe).
			Build(&entity.TestScenario{})

		assert.NotNil(t, b)
	})

}
