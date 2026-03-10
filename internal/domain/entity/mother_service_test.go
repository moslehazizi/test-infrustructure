package entity

import (
	"control-panel-service/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMotherService_Validate(t *testing.T) {
	t.Run("success_case_no_delay", func(t *testing.T) {
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

	t.Run("success_case_fixed_delay", func(t *testing.T) {
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

	t.Run("success_case_random_delay", func(t *testing.T) {
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

	t.Run("failed_case_name_is_missing", func(t *testing.T) {
		service := MotherService{
			DatabaseName: "db",
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMotherServiceName)
	})

	t.Run("failed_case_delay_rate_is_negative", func(t *testing.T) {
		service := MotherService{
			Name:              "test",
			ResponseDelayRate: -1,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidResponseDelayRate)
	})

	t.Run("failed_case_delay_rate_is_bigger_than_100", func(t *testing.T) {
		service := MotherService{
			Name:              "test",
			ResponseDelayRate: 101,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidResponseDelayRate)
	})

	t.Run("failed_case_exception_rate_is_negative", func(t *testing.T) {
		service := MotherService{
			Name:          "test",
			ExceptionRate: -1,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidExceptionRate)
	})

	t.Run("failed_case_exception_rate_is_bigger_than_100", func(t *testing.T) {
		service := MotherService{
			Name:          "test",
			ExceptionRate: 101,
		}

		err := service.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidExceptionRate)
	})

	t.Run("failed_case_fixed_delay_is_set_but_rate_is_0", func(t *testing.T) {
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

	t.Run("failed_case_random_delay_is_set_but_rate_is_0", func(t *testing.T) {
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

	t.Run("failed_case_both_fixed_and_random_delay_fields_are_set", func(t *testing.T) {
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

	t.Run("failed_case_min_is_greater_than_max", func(t *testing.T) {
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

	t.Run("failed_case_only_min_random_delay_is_provided", func(t *testing.T) {
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

	t.Run("failed_case_only_max_random_delay_is_provided", func(t *testing.T) {
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
