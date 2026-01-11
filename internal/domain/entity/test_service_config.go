package entity

import "time"

type TestServiceConfig struct {
	ID                    uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	TestScenarioID        uint64    `gorm:"column:test_scenario_id"`
	MaxRequests           int       `gorm:"column:max_requests"`
	MaxDuration           int       `gorm:"column:max_duration"`
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
}

func (TestServiceConfig) TableName() string {
	return "test_service_configs"
}
