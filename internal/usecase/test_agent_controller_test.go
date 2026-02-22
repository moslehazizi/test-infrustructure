package usecase_test

import (
	"control-panel-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestAgentController(t *testing.T) {
	ctrl := usecase.NewTestAgentController()

	assert.NotNil(t, ctrl)
}
