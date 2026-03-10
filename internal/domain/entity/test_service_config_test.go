package entity

import (
	"control-panel-service/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestServiceConfig_TableName(t *testing.T) {
	testSvcCfg := TestServiceConfig{}

	name := testSvcCfg.TableName()

	assert.Equal(t, name, "test_service_configs")
}

func TestTestServiceConfig_Validate(t *testing.T) {
	t.Run("failed_case_invalid_max_request", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMaxRequest)
	})
	t.Run("failed_case_invalid_max_duration", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests: 1,
			MaxDuration: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMaxDuration)
	})
	t.Run("failed_case_request_delay_duration", func(t *testing.T) {
		sampleInt := -1
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidRequestDelayDuration)
	})
	t.Run("failed_case_no_request_delay_duration_or_fixed_request_delay_duration", func(t *testing.T) {
		sampleInt := 0
		testSvcCfg := TestServiceConfig{
			MaxRequests:           1,
			MaxDuration:           0,
			RequestDelayDuration:  &sampleInt,
			RandomRequestDelayMin: &sampleInt,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidRequestDelayDurationConfig)
	})
	t.Run("failed_case_min_random_delay_request_should_be_less_than_max_random_delay_duration", func(t *testing.T) {
		min := 20
		max := 10
		testSvcCfg := TestServiceConfig{
			MaxRequests:           1,
			MaxDuration:           0,
			RandomRequestDelayMin: &min,
			RandomRequestDelayMax: &max,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMinDelayDurationMoreThanMax)
	})

	t.Run("failed_case_all_delay_request_config_couldn't_be_null_at_the_same_time", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests: 1,
			MaxDuration: 0,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidRequestDelayDurationConfig)
	})
	t.Run("success_case_delay_is_0_and_random_delays_are_null", func(t *testing.T) {
		d := 0
		n := 10
		databaseName := "db1"
		databaseTableName := "table1"
		testSvcCfg := TestServiceConfig{
			MaxRequests:           1,
			MaxDuration:           0,
			RequestDelayDuration:  &d,
			RandomRequestDelayMax: nil,
			RandomRequestDelayMin: nil,
			FixedTestNumber:       &n,
			DatabaseName:          databaseName,
			DatabaseTableName:     databaseTableName,
		}

		err := testSvcCfg.Validate()

		assert.NoError(t, err)
	})

	t.Run("failed_case_invalid_fixed_test_number", func(t *testing.T) {
		sampleInt := -1
		sampleUInt := 1
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleUInt,
			FixedTestNumber:      &sampleInt,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidFixedTestNumber)
	})
	t.Run("failed_case_invalid_fixed_test_number_configuration", func(t *testing.T) {
		sampleInt := 2
		min := 10
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			RandomTestNumberMin:  &min,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidFixedTestNumberConfig)
	})
	t.Run("failed_case_random_min_test_number_couldn't_be_more_than_max", func(t *testing.T) {
		min := 20
		max := 10
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &min,
			RandomTestNumberMin:  &min,
			RandomTestNumberMax:  &max,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMinRandomTestNumberMoreThanMax)
	})
	t.Run("failed_case_random_min_or_max_test_number_less_than_zero", func(t *testing.T) {
		min := -10
		max := 10
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &max,
			RandomTestNumberMin:  &min,
			RandomTestNumberMax:  &max,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMinOrMaxRandomTestNumber)
	})

	t.Run("failed_case_all_fixed_and_random_value_couldn't_be_null_at_the_same_time", func(t *testing.T) {
		sampleInt := 10
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidTestNumberConfig)
	})
	t.Run("failed_case_bad_value_rate_less_than_0_or_more_than_100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:  1,
			MaxDuration:  0,
			BadValueRate: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidBadValueRate)
	})
	t.Run("failed_case_negative_value_rate_less_than_0_or_more_than_100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:       1,
			MaxDuration:       0,
			NegativeValueRate: 101,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidNegativeValueRate)
	})
	t.Run("failed_case_real_value_rate_less_than_0_or_more_than_100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:   1,
			MaxDuration:   0,
			RealValueRate: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidRealValueRate)
	})
	t.Run("failed_case_zero_value_rate_less_than_0_or_more_than_100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:   1,
			MaxDuration:   0,
			ZeroValueRate: 110,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidZeroValueRate)
	})
	t.Run("failed_case_string_value_rate_less_than_0_or_more_than_100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:     1,
			MaxDuration:     0,
			StringValueRate: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidStringValueRate)
	})
	t.Run("failed_case_long_string_value_rate_less_than_0_or_more_than_100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:         1,
			MaxDuration:         0,
			LongStringValueRate: 101,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidLongStringValueRate)
	})
	t.Run("failed_case_null_value_rate_less_than_0_or_more_than_100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:   1,
			MaxDuration:   0,
			NullValueRate: -2,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidNullValueRate)
	})
	t.Run("success_case", func(t *testing.T) {
		sampleInt := 2
		databaseName := "db1"
		databaseTableName := "table1"
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			BadValueRate:         0,
			NegativeValueRate:    0,
			RealValueRate:        0,
			ZeroValueRate:        0,
			StringValueRate:      0,
			LongStringValueRate:  0,
			NullValueRate:        0,
			DatabaseName:         databaseName,
			DatabaseTableName:    databaseTableName,
		}

		err := testSvcCfg.Validate()

		assert.Nil(t, err)
	})
	t.Run("failed_case_if_bad_value_rate_is_zero_sum_of_all_bad_value_should_be_zero", func(t *testing.T) {
		sampleInt := 2
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			BadValueRate:         0,
			NegativeValueRate:    0,
			RealValueRate:        0,
			ZeroValueRate:        0,
			StringValueRate:      1,
			LongStringValueRate:  0,
			NullValueRate:        0,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidZeroSumOfBadValues)
	})

	t.Run("failed_case_if_bad_value_rate_not_zero_sum_of_all_bad_value_should_be_100", func(t *testing.T) {
		sampleInt := 2
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			BadValueRate:         10,
			NegativeValueRate:    0,
			RealValueRate:        0,
			ZeroValueRate:        0,
			StringValueRate:      5,
			LongStringValueRate:  0,
			NullValueRate:        0,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalid100SumOfBadValues)
	})

	t.Run("failed_case_database_name_is_empty", func(t *testing.T) {
		sampleInt := 2
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			BadValueRate:         0,
			NegativeValueRate:    0,
			RealValueRate:        0,
			ZeroValueRate:        0,
			StringValueRate:      0,
			LongStringValueRate:  0,
			NullValueRate:        0,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDatabaseName)
	})

	t.Run("failed_case_table_name_is_empty", func(t *testing.T) {
		sampleInt := 2
		databaseName := "db1"
		testSvcCfg := TestServiceConfig{
			MaxRequests:          1,
			MaxDuration:          0,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
			BadValueRate:         0,
			NegativeValueRate:    0,
			RealValueRate:        0,
			ZeroValueRate:        0,
			StringValueRate:      0,
			LongStringValueRate:  0,
			NullValueRate:        0,
			DatabaseName:         databaseName,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidDatabaseTableName)
	})
}
