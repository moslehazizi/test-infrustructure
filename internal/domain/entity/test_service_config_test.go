package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestServiceConfig_TableName(t *testing.T) {
	testSvcCfg := TestServiceConfig{}

	name := testSvcCfg.TableName()

	assert.Equal(t, name, "test_service_configs")
}

func TestTestServiceConfig_Validate(t *testing.T) {
	t.Run("success case", func(t *testing.T) {

	})
}
