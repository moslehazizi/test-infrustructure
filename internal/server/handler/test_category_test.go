package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestCategoryHandler(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("ok", func(t *testing.T) {
		h := NewTestCategoryHandler(&cfg, &mocks.MockTestCategoryService{})
		require.NotNil(t, h.testCategoryService)
		require.NotNil(t, h.cfg)
	})
}
