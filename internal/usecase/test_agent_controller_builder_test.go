package usecase_test

import (
	"control-panel-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestAgentControllerBuilder_Build(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		b := usecase.NewTestAgentControllerBuilder().Build()

		assert.NotNil(t, b)
	})

}
