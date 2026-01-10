package entity

import (
	"control-panel-service/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMotherService_Validate(t *testing.T) {
	t.Run("success case - no delay", func(t *testing.T) {
		service := MotherService{
			Name:              "mother",
			DatabaseName:      "mother1",
			DatabaseTableName: "table",
			ExceptionRate:     10,
			ResponseDelayRate: 0,
		}

		err := service.Validate()

		assert.Nil(t, err)
	})

	t.Run("success case - fixed delay", func(t *testing.T) {
		duration := 100
		service := MotherService{
			Name:                  "mother",
			DatabaseName:          "mother1",
			DatabaseTableName:     "table",
			ResponseDelayRate:     10,
			ResponseDelayDuration: &duration,
		}

		err := service.Validate()

		assert.Nil(t, err)
	})

	t.Run("success case - random delay", func(t *testing.T) {
		minDelay := 10
		maxDelay := 50
		service := MotherService{
			Name:                   "mother",
			DatabaseName:           "mother1",
			DatabaseTableName:      "table",
			ResponseDelayRate:      10,
			RandomResponseDelayMin: &minDelay,
			RandomResponseDelayMax: &maxDelay,
		}

		err := service.Validate()

		assert.Nil(t, err)
	})

	t.Run("failed case - name is missing", func(t *testing.T) {
		service := MotherService{
			DatabaseName: "db",
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidName)
	})

	t.Run("failed case - delay rate is negative", func(t *testing.T) {
		service := MotherService{
			Name:              "test",
			ResponseDelayRate: -1,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidResponseDelayRate)
	})

	t.Run("failed case - delay rate is bigger than 100", func(t *testing.T) {
		service := MotherService{
			Name:              "test",
			ResponseDelayRate: 101,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidResponseDelayRate)
	})

	t.Run("failed case - exception rate is negative", func(t *testing.T) {
		service := MotherService{
			Name:          "test",
			ExceptionRate: -1,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidExceptionRate)
	})

	t.Run("failed case - exception rate is bigger than 100", func(t *testing.T) {
		service := MotherService{
			Name:          "test",
			ExceptionRate: 101,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidExceptionRate)
	})

	t.Run("failed case - fixed delay is set but rate is 0", func(t *testing.T) {
		duration := 100
		service := MotherService{
			Name:                  "mother",
			DatabaseName:          "mother1",
			DatabaseTableName:     "table",
			ResponseDelayRate:     0,
			ResponseDelayDuration: &duration,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDelayConfiguration)
	})

	t.Run("failed case - random delay is set but rate is 0", func(t *testing.T) {
		minDelay := 10
		maxDelay := 50
		service := MotherService{
			Name:                   "mother",
			DatabaseName:           "mother1",
			DatabaseTableName:      "table",
			ResponseDelayRate:      0,
			RandomResponseDelayMin: &minDelay,
			RandomResponseDelayMax: &maxDelay,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDelayConfiguration)
	})

	t.Run("failed case - both fixed and random delay fields are set", func(t *testing.T) {
		minDelay := 10
		duration := 100
		service := MotherService{
			Name:                   "mother",
			DatabaseName:           "mother1",
			DatabaseTableName:      "table",
			ResponseDelayRate:      10,
			ResponseDelayDuration:  &duration,
			RandomResponseDelayMin: &minDelay,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDelayConfiguration)
	})

	t.Run("failed case - min is greater than max", func(t *testing.T) {
		max := 10
		min := 100
		service := MotherService{
			Name:                   "mother",
			DatabaseName:           "mother1",
			DatabaseTableName:      "table",
			ResponseDelayRate:      10,
			RandomResponseDelayMin: &min,
			RandomResponseDelayMax: &max,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidRandomDelayRange)
	})

	t.Run("failed case - only min random delay is provided", func(t *testing.T) {
		min := 10
		service := MotherService{
			Name:                   "mother",
			DatabaseName:           "mother1",
			DatabaseTableName:      "table",
			ResponseDelayRate:      10,
			RandomResponseDelayMin: &min,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDelayConfiguration)
	})

	t.Run("failed case - only max random delay is provided", func(t *testing.T) {
		max := 10
		service := MotherService{
			Name:                   "mother",
			DatabaseName:           "mother1",
			DatabaseTableName:      "table",
			ResponseDelayRate:      10,
			RandomResponseDelayMax: &max,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDelayConfiguration)
	})
}
