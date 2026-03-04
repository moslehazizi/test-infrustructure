package request

import (
	"control-panel-service/internal/domain/entity"

	"github.com/google/uuid"
)

type RunRequest struct {
	StepNum               int       `json:"step_num"`
	ExecutionId           uuid.UUID `json:"execution_id"`
	TestScenarioID        uint64    `json:"test_scenario_id"`
	MaxRequests           int       `json:"max_requests"`
	MaxDuration           int64     `json:"max_duration"`
	RequestDelayDuration  *int      `json:"request_delay_duration"`
	RandomRequestDelayMin *int      `json:"random_request_delay_min"`
	RandomRequestDelayMax *int      `json:"random_request_delay_max"`
	FixedTestNumber       *int      `json:"fixed_test_number"`
	RandomTestNumberMin   *int      `json:"random_test_number_min"`
	RandomTestNumberMax   *int      `json:"random_test_number_max"`
	BadValueRate          int       `json:"bad_value_rate"`
	NegativeValueRate     int       `json:"negative_value_rate"`
	RealValueRate         int       `json:"real_value_rate"`
	ZeroValueRate         int       `json:"zero_value_rate"`
	StringValueRate       int       `json:"string_value_rate"`
	LongStringValueRate   int       `json:"long_string_value_rate"`
	NullValueRate         int       `json:"null_value_rate"`
	DatabaseName          string    `json:"database_name"`
	DatabaseTableName     string    `json:"database_table_name"`
}

func NewRunRequestFromTestServiceConfig(stepNum int, executionID uuid.UUID, cfg *entity.TestServiceConfig) *RunRequest {
	return &RunRequest{
		StepNum:               stepNum,
		ExecutionId:           executionID,
		TestScenarioID:        cfg.TestScenarioID,
		MaxRequests:           stepNum * cfg.MaxRequests,
		MaxDuration:           cfg.MaxDuration,
		RequestDelayDuration:  cfg.RequestDelayDuration,
		RandomRequestDelayMin: cfg.RandomRequestDelayMin,
		RandomRequestDelayMax: cfg.RandomRequestDelayMax,
		FixedTestNumber:       cfg.FixedTestNumber,
		RandomTestNumberMin:   cfg.RandomTestNumberMin,
		RandomTestNumberMax:   cfg.RandomTestNumberMax,
		BadValueRate:          cfg.BadValueRate,
		NegativeValueRate:     cfg.NegativeValueRate,
		RealValueRate:         cfg.RealValueRate,
		ZeroValueRate:         cfg.ZeroValueRate,
		StringValueRate:       cfg.StringValueRate,
		LongStringValueRate:   cfg.LongStringValueRate,
		NullValueRate:         cfg.NullValueRate,
		DatabaseName:          cfg.DatabaseName,
		DatabaseTableName:     cfg.DatabaseTableName,
	}
}
