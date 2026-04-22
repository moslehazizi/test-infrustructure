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

func (ts *TestScenario) ToTestScenarioEntity(entityTS *entity.TestScenario) {
	entityTS.Name = ts.Name
	entityTS.TestCategoryID = ts.TestCategoryID
	entityTS.MotherServiceID = ts.MotherServiceID
	entityTS.MaxTestServiceCount = ts.MaxTestServiceCount
	entityTS.NumSteps = ts.NumSteps
	entityTS.IncreaseAgentNumber = ts.IncreaseAgentNumber
	entityTS.ExecNumMultiAgent = ts.ExecNumMultiAgent
	entityTS.TestServiceConfig = func() *entity.TestServiceConfig {
		if ts.Config == nil {
			return nil
		}

		return &entity.TestServiceConfig{
			MaxRequests:            ts.Config.MaxRequests,
			MaxDuration:            int64(ts.Config.MaxDuration),
			RequestDelayDuration:   ts.Config.RequestDelayDuration,
			RandomRequestDelayMin:  ts.Config.RandomRequestDelayMin,
			RandomRequestDelayMax:  ts.Config.RandomRequestDelayMax,
			FixedTestNumber:        ts.Config.FixedTestNumber,
			RandomTestNumberMin:    ts.Config.RandomTestNumberMin,
			RandomTestNumberMax:    ts.Config.RandomTestNumberMax,
			BadValueRate:           ts.Config.BadValueRate,
			NegativeValueRate:      ts.Config.NegativeValueRate,
			ZeroValueRate:          ts.Config.ZeroValueRate,
			StringValueRate:        ts.Config.StringValueRate,
			RealValueRate:          ts.Config.RealValueRate,
			LongStringValueRate:    ts.Config.LongStringValueRate,
			NullValueRate:          ts.Config.NullValueRate,
			DatabaseName:           ts.Config.DatabaseName,
			DatabaseTableName:      ts.Config.DatabaseTableName,
			IncreaseFixedInput:     ts.Config.IncreaseFixedInput,
			ExecNumMultiFixedInput: ts.Config.ExecNumMultiFixedInput,
		}
	}()
}
