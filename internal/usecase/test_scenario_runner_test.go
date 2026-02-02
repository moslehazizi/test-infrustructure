package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryScenarioExecutorEngine(t *testing.T) {
	eng := NewInMemoryScenarioExecutorEngine()
	assert.NotNil(t, eng)
}
