package response

import (
	"control-panel-service/internal/domain/entity"
	"time"
)

type TestScenario struct {
	ID                  uint64                `json:"id"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
	StartedAt           *time.Time            `json:"started_at"`
	Name                string                `json:"name"`
	TestCategory        *TestCategory         `json:"test_category"`
	MotherService       *MotherService        `json:"mother_service"`
	Status              entity.ScenarioStatus `json:"status"`
	MaxTestServiceCount *int64                `json:"max_test_service_count"`
	TestServiceConfig   *TestServiceConfig    `json:"test_service_config"`
	Editable            bool                  `json:"editable"`
	NumSteps            int64                 `json:"num_steps"`
	IncreaseAgentNumber int64                 `json:"increase_agent_number"`
	ExecNumMultiAgent   int64                 `json:"execution_number_multi_agent"`
}

type TestServiceConfig struct {
	ID                     uint64    `json:"id"`
	MaxRequests            int       `json:"max_requests"`
	MaxDuration            int64     `json:"max_duration"`
	RequestDelayDuration   *int      `json:"request_delay_duration"`
	RandomRequestDelayMin  *int      `json:"random_request_delay_min"`
	RandomRequestDelayMax  *int      `json:"random_request_delay_max"`
	FixedTestNumber        *int      `json:"fixed_test_number"`
	RandomTestNumberMin    *int      `json:"random_test_number_min"`
	RandomTestNumberMax    *int      `json:"random_test_number_max"`
	BadValueRate           int       `json:"bad_value_rate"`
	NegativeValueRate      int       `json:"negative_value_rate"`
	RealValueRate          int       `json:"real_value_rate"`
	ZeroValueRate          int       `json:"zero_value_rate"`
	StringValueRate        int       `json:"string_value_rate"`
	LongStringValueRate    int       `json:"long_string_value_rate"`
	NullValueRate          int       `json:"null_value_rate"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	DatabaseName           string    `json:"database_name"`
	DatabaseTableName      string    `json:"database_table_name"`
	IncreaseFixedInput     int       `json:"increase_fixed_input"`
	ExecNumMultiFixedInput int       `json:"execution_number_multi_fixed_input"`
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
