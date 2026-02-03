package response

import (
	"control-panel-service/internal/domain/entity"
	"time"
)

type TestScenario struct {
	ID                  uint64                `json:"id"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
	Name                string                `json:"name"`
	TestCategory        *TestCategory         `json:"test_category"`
	MotherService       *MotherService        `json:"mother_service"`
	Status              entity.ScenarioStatus `json:"status"`
	MaxTestServiceCount *int                  `json:"max_test_service_count"`
	ExecutionDuration   *int                  `json:"execution_duration"`
	AutoStepChangeRate  *int                  `json:"auto_step_change_rate"`
	TestServiceConfig   *TestServiceConfig    `json:"test_service_config"`
}

type TestServiceConfig struct {
	ID                    uint64    `json:"id"`
	MaxRequests           int       `json:"max_requests"`
	MaxDuration           int       `json:"max_duration"`
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
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type PaginatedTestScenario struct {
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
	Data    []TestScenario `json:"data"`
	Total   int64          `json:"total"`
}

type TestScenarioResponseByID struct {
	Data TestScenario `json:"data"`
}
