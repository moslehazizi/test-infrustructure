package request

import "control-panel-service/internal/domain/entity"

type TestScenario struct {
	Name                string                    `json:"name"`
	TestCategoryID      uint64                    `json:"test_category_id"`
	MotherServiceID     uint64                    `json:"mother_service_id"`
	MaxTestServiceCount *int64                    `json:"max_test_service_count"`
	Config              *TestServiceConfigRequest `json:"test_service_config"`
	NumSteps            int64                     `json:"num_steps"`
	IncreaseAgentNumber int64                     `json:"increase_agent_number"`
	ExecNumMultiAgent   int64                     `json:"execution_number_multi_agent"`
}

type TestServiceConfigRequest struct {
	MaxRequests            int    `json:"max_requests"`
	MaxDuration            int    `json:"max_duration"`
	RequestDelayDuration   *int   `json:"request_delay_duration"`
	RandomRequestDelayMin  *int   `json:"random_request_delay_min"`
	RandomRequestDelayMax  *int   `json:"random_request_delay_max"`
	FixedTestNumber        *int   `json:"fixed_test_number"`
	RandomTestNumberMin    *int   `json:"random_test_number_min"`
	RandomTestNumberMax    *int   `json:"random_test_number_max"`
	BadValueRate           int    `json:"bad_value_rate"`
	NegativeValueRate      int    `json:"negative_value_rate"`
	RealValueRate          int    `json:"real_value_rate"`
	ZeroValueRate          int    `json:"zero_value_rate"`
	StringValueRate        int    `json:"string_value_rate"`
	LongStringValueRate    int    `json:"long_string_value_rate"`
	NullValueRate          int    `json:"null_value_rate"`
	DatabaseName           string `json:"database_name"`
	DatabaseTableName      string `json:"database_table_name"`
	IncreaseFixedInput     int    `json:"increase_fixed_input"`
	ExecNumMultiFixedInput int    `json:"execution_number_multi_fixed_input"`
}

type TestScenarioUpdateRequest struct {
	ID                  uint64                    `json:"id"`
	Name                string                    `json:"name"`
	MotherServiceID     uint64                    `json:"mother_service_id"`
	MaxTestServiceCount *int64                    `json:"max_test_service_count"`
	Config              *TestServiceConfigRequest `json:"test_service_config"`
	NumSteps            int64                     `json:"num_steps"`
	IncreaseAgentNumber int64                     `json:"increase_agent_number"`
	ExecNumMultiAgent   int64                     `json:"execution_number_multi_agent"`
}

type TestScenarioPaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (tsr *TestScenarioUpdateRequest) ToTestScenarioEntity(entitySCI *entity.TestScenario, entityMother *entity.MotherService) {
	entitySCI.Name = tsr.Name
	entitySCI.MotherServiceID = tsr.MotherServiceID
	entitySCI.MotherService = entityMother
	entitySCI.IncreaseAgentNumber = tsr.IncreaseAgentNumber
	entitySCI.ExecNumMultiAgent = tsr.ExecNumMultiAgent
	entitySCI.Status = entity.ScenarioStatusReady

	if tsr.MaxTestServiceCount != nil {
		entitySCI.MaxTestServiceCount = tsr.MaxTestServiceCount
	}

	if tsr.NumSteps >= 1 {
		entitySCI.NumSteps = tsr.NumSteps
	}
}

func (tscr *TestServiceConfigRequest) ToTestServiceConfigEntity(entityCFG *entity.TestServiceConfig) {
	entityCFG.MaxRequests = tscr.MaxRequests
	entityCFG.MaxDuration = int64(tscr.MaxDuration)
	entityCFG.RequestDelayDuration = tscr.RequestDelayDuration
	entityCFG.RandomRequestDelayMin = tscr.RandomRequestDelayMin
	entityCFG.RandomRequestDelayMax = tscr.RandomRequestDelayMax
	entityCFG.FixedTestNumber = tscr.FixedTestNumber
	entityCFG.RandomTestNumberMin = tscr.RandomTestNumberMin
	entityCFG.RandomTestNumberMax = tscr.RandomTestNumberMax
	entityCFG.BadValueRate = tscr.BadValueRate
	entityCFG.NegativeValueRate = tscr.NegativeValueRate
	entityCFG.ZeroValueRate = tscr.ZeroValueRate
	entityCFG.StringValueRate = tscr.StringValueRate
	entityCFG.RealValueRate = tscr.RealValueRate
	entityCFG.LongStringValueRate = tscr.LongStringValueRate
	entityCFG.NullValueRate = tscr.NullValueRate
	entityCFG.DatabaseName = tscr.DatabaseName
	entityCFG.DatabaseTableName = tscr.DatabaseTableName
	entityCFG.IncreaseFixedInput = tscr.IncreaseFixedInput
	entityCFG.ExecNumMultiFixedInput = tscr.ExecNumMultiFixedInput
}
