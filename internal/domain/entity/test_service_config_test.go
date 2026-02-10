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
	t.Run("failed case - invalid max request", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMaxRequest)
	})
	t.Run("failed case - invalid max duration", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests: 1,
			MaxDuration: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidMaxDuration)
	})
	t.Run("failed case - request delay duration ", func(t *testing.T) {
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
	t.Run("failed case - no request delay duration or fixed request delay duration", func(t *testing.T) {
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
	t.Run("failed case - min random delay request should be less than max random delay duration", func(t *testing.T) {
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

	t.Run("failed case - all delay request config couldn't be null at the same time", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests: 1,
			MaxDuration: 0,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidRequestDelayDurationConfig)
	})
	t.Run("success case: delay is 0 and random delays are null", func(t *testing.T) {
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

	t.Run("failed case - invalid fixed test number", func(t *testing.T) {
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
	t.Run("failed case - invalid fixed test number configuration", func(t *testing.T) {
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
	t.Run("failed case - random min test number couldn't be more than max", func(t *testing.T) {
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
	t.Run("failed case - random min or max test number less than zero", func(t *testing.T) {
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

	t.Run("failed case - all fixed and random value couldn't be null at the same time", func(t *testing.T) {
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
	t.Run("failed case - bad value rate less than 0 or more than 100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:  1,
			MaxDuration:  0,
			BadValueRate: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidBadValueRate)
	})
	t.Run("failed case - negative value rate less than 0 or more than 100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:       1,
			MaxDuration:       0,
			NegativeValueRate: 101,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidNegativeValueRate)
	})
	t.Run("failed case - real value rate less than 0 or more than 100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:   1,
			MaxDuration:   0,
			RealValueRate: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidRealValueRate)
	})
	t.Run("failed case - zero value rate less than 0 or more than 100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:   1,
			MaxDuration:   0,
			ZeroValueRate: 110,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidZeroValueRate)
	})
	t.Run("failed case - string value rate less than 0 or more than 100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:     1,
			MaxDuration:     0,
			StringValueRate: -1,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidStringValueRate)
	})
	t.Run("failed case - long string value rate less than 0 or more than 100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:         1,
			MaxDuration:         0,
			LongStringValueRate: 101,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidLongStringValueRate)
	})
	t.Run("failed case - null value rate less than 0 or more than 100", func(t *testing.T) {
		testSvcCfg := TestServiceConfig{
			MaxRequests:   1,
			MaxDuration:   0,
			NullValueRate: -2,
		}

		err := testSvcCfg.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrInvalidNullValueRate)
	})
	t.Run("success case", func(t *testing.T) {
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
	t.Run("failed case - if bad value rate is zero - sum of all bad value should be zero", func(t *testing.T) {
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

	t.Run("failed case - if bad value rate not zero - sum of all bad value should be 100", func(t *testing.T) {
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

	t.Run("failed case - database name is empty", func(t *testing.T) {
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

	t.Run("failed case - table name is empty", func(t *testing.T) {
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
