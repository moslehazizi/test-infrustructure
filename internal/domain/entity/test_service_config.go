package entity

import (
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/pkg"
	"time"
)

type TestServiceConfig struct {
	ID                    uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	TestScenarioID        uint64    `gorm:"column:test_scenario_id"`
	MaxRequests           int       `gorm:"column:max_requests"`
	MaxDuration           int64     `gorm:"column:max_duration"`
	RequestDelayDuration  *int      `gorm:"column:request_delay_duration"`
	RandomRequestDelayMin *int      `gorm:"column:random_request_delay_min"`
	RandomRequestDelayMax *int      `gorm:"column:random_request_delay_max"`
	FixedTestNumber       *int      `gorm:"column:fixed_test_number"`
	RandomTestNumberMin   *int      `gorm:"column:random_test_number_min"`
	RandomTestNumberMax   *int      `gorm:"column:random_test_number_max"`
	BadValueRate          int       `gorm:"column:bad_value_rate"`
	NegativeValueRate     int       `gorm:"column:negative_value_rate"`
	RealValueRate         int       `gorm:"column:real_value_rate"`
	ZeroValueRate         int       `gorm:"column:zero_value_rate"`
	StringValueRate       int       `gorm:"column:string_value_rate"`
	LongStringValueRate   int       `gorm:"column:long_string_value_rate"`
	NullValueRate         int       `gorm:"column:null_value_rate"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
	DatabaseName          string    `gorm:"column:database_name"`
	DatabaseTableName     string    `gorm:"column:database_table_name"`
	// This field specifies the increment value for the fixed test number.
	// It is used when we want to re‑run the entire scenario,
	// and in each re-run we want the fixed test number to increase.
	// This number represents the increment amount.
	// For example, if the fixed test number is 1000 and this increment value is 2000,
	// then in the first run, the fixed number will be 1000, and in the next run of the scenario,
	// the number will become 3000.
	IncreaseFixedInput int `gorm:"column:increase_fixed_input"`
	// This field works together with the previous one and determines
	// how many times the fixed test number should be incremented.
	ExecNumMultiFixedInput int `gorm:"column:execution_number_multi_fixed_input"`
}

func (TestServiceConfig) TableName() string {
	return "test_service_configs"
}

// nolint
func (t *TestServiceConfig) Validate() error {
	// validate public rules
	if t.MaxRequests < 0 {
		return pkg.ErrInvalidMaxRequest
	}
	if t.MaxDuration < 0 {
		return pkg.ErrInvalidMaxDuration
	}
	if t.BadValueRate < 0 || t.BadValueRate > 100 {
		return pkg.ErrInvalidBadValueRate
	}
	if t.NegativeValueRate < 0 || t.NegativeValueRate > 100 {
		return pkg.ErrInvalidNegativeValueRate
	}
	if t.RealValueRate < 0 || t.RealValueRate > 100 {
		return pkg.ErrInvalidRealValueRate
	}
	if t.ZeroValueRate < 0 || t.ZeroValueRate > 100 {
		return pkg.ErrInvalidZeroValueRate
	}
	if t.StringValueRate < 0 || t.StringValueRate > 100 {
		return pkg.ErrInvalidStringValueRate
	}
	if t.LongStringValueRate < 0 || t.LongStringValueRate > 100 {
		return pkg.ErrInvalidLongStringValueRate
	}
	if t.NullValueRate < 0 || t.NullValueRate > 100 {
		return pkg.ErrInvalidNullValueRate
	}

	// validate delay duration between two transaction
	// if t.RequestDelayDuration == nil && t.RandomRequestDelayMax == nil && t.RandomRequestDelayMin == nil {
	// 	return pkg.ErrInvalidRequestDelayDurationConfig
	// }
	if t.RequestDelayDuration != nil && *t.RequestDelayDuration < 0 {
		return pkg.ErrInvalidRequestDelayDuration
	}

	if t.RequestDelayDuration != nil && *t.RequestDelayDuration > 0 {
		if t.RandomRequestDelayMin != nil || t.RandomRequestDelayMax != nil {
			return pkg.ErrInvalidRequestDelayDurationConfig
		}
	}

	if t.RandomRequestDelayMin != nil && t.RandomRequestDelayMax != nil {
		if *t.RandomRequestDelayMin >= *t.RandomRequestDelayMax {
			return pkg.ErrMinDelayDurationMoreThanMax
		}
	}

	// validate number that send via transaction - fix number or random number
	if t.FixedTestNumber == nil && t.RandomTestNumberMin == nil && t.RandomTestNumberMax == nil {
		return pkg.ErrInvalidTestNumberConfig
	}
	if t.FixedTestNumber != nil && *t.FixedTestNumber <= 0 {
		return pkg.ErrInvalidFixedTestNumber
	}

	if t.FixedTestNumber != nil {
		if t.RandomTestNumberMin != nil || t.RandomTestNumberMax != nil {
			return pkg.ErrInvalidFixedTestNumberConfig
		}
	}

	if t.RandomTestNumberMin != nil && t.RandomTestNumberMax != nil {
		if *t.RandomTestNumberMin < 0 || *t.RandomTestNumberMax < 0 {
			return pkg.ErrInvalidMinOrMaxRandomTestNumber
		}
		if *t.RandomTestNumberMin >= *t.RandomTestNumberMax {
			return pkg.ErrMinRandomTestNumberMoreThanMax
		}
	}

	if t.BadValueRate == 0 {
		if t.NegativeValueRate+t.RealValueRate+t.StringValueRate+t.NullValueRate+t.LongStringValueRate+t.ZeroValueRate != 0 {
			return pkg.ErrInvalidZeroSumOfBadValues
		}
	}
	if t.BadValueRate > 0 {
		if t.NegativeValueRate+t.RealValueRate+t.StringValueRate+t.NullValueRate+t.LongStringValueRate+t.ZeroValueRate != 100 {
			return pkg.ErrInvalid100SumOfBadValues
		}
	}

	if t.DatabaseName == "" {
		return pkg.ErrInvalidDatabaseName
	}

	if t.DatabaseTableName == "" {
		return pkg.ErrInvalidDatabaseTableName
	}
	if t.IncreaseFixedInput < 0 {
		return pkg.ErrIncreaseFixedInputNotBeNegative
	}
	if t.ExecNumMultiFixedInput < 0 {
		return pkg.ErrExecNumMultiFixedInputShouldBePositive
	}
	if (t.IncreaseFixedInput == 0 && t.ExecNumMultiFixedInput != 0) || (t.ExecNumMultiFixedInput == 0 && t.IncreaseFixedInput != 0) {
		return pkg.ErrMultiFixedInputConfigNotTrue
	}

	return nil
}

// TODO: convert request to entity instead.
func (tsc *TestServiceConfig) ApplyUpdateFromRequest(cfg *request.TestServiceConfigRequest) {
	tsc.MaxRequests = cfg.MaxRequests
	tsc.MaxDuration = int64(cfg.MaxDuration)
	tsc.RequestDelayDuration = cfg.RequestDelayDuration
	tsc.RandomRequestDelayMin = cfg.RandomRequestDelayMin
	tsc.RandomRequestDelayMax = cfg.RandomRequestDelayMax
	tsc.FixedTestNumber = cfg.FixedTestNumber
	tsc.RandomTestNumberMin = cfg.RandomTestNumberMin
	tsc.RandomTestNumberMax = cfg.RandomTestNumberMax
	tsc.BadValueRate = cfg.BadValueRate
	tsc.NegativeValueRate = cfg.NegativeValueRate
	tsc.ZeroValueRate = cfg.ZeroValueRate
	tsc.StringValueRate = cfg.StringValueRate
	tsc.RealValueRate = cfg.RealValueRate
	tsc.LongStringValueRate = cfg.LongStringValueRate
	tsc.NullValueRate = cfg.NullValueRate
	tsc.DatabaseName = cfg.DatabaseName
	tsc.DatabaseTableName = cfg.DatabaseTableName
	tsc.IncreaseFixedInput = cfg.IncreaseFixedInput
	tsc.ExecNumMultiFixedInput = cfg.ExecNumMultiFixedInput
}

// in: internal/entity/test_service_config.go
