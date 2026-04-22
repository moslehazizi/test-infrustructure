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

func (resp *TestScenario) FromTestScenarioEntity(entityTS *entity.TestScenario) {
	resp.ID = entityTS.ID
	resp.Name = entityTS.Name
	resp.CreatedAt = entityTS.CreatedAt
	resp.UpdatedAt = entityTS.UpdatedAt
	resp.StartedAt = entityTS.StartedAt
	resp.Status = entityTS.Status
	resp.MaxTestServiceCount = entityTS.MaxTestServiceCount
	resp.NumSteps = entityTS.NumSteps
	resp.IncreaseAgentNumber = entityTS.IncreaseAgentNumber
	resp.ExecNumMultiAgent = entityTS.ExecNumMultiAgent
	resp.Editable = entityTS.Editable
	resp.TestCategory = func() *TestCategory {
		if entityTS.TestCategory == nil {
			return nil
		}

		return &TestCategory{
			ID:                     entityTS.TestCategory.ID,
			Name:                   entityTS.TestCategory.Name,
			Label:                  entityTS.TestCategory.Label,
			HasMaxTestServiceCount: entityTS.TestCategory.HasMaxTestServiceCount,
			HasNumSteps:            entityTS.TestCategory.HasNumSteps,
			CreatedAt:              entityTS.TestCategory.CreatedAt,
			UpdatedAt:              entityTS.TestCategory.UpdatedAt,
		}
	}()
	resp.MotherService = func() *MotherService {
		if entityTS.MotherService == nil {
			return nil
		}

		return &MotherService{
			ID:                       entityTS.MotherService.ID,
			CreatedAt:                entityTS.MotherService.CreatedAt,
			UpdatedAt:                entityTS.MotherService.UpdatedAt,
			Name:                     entityTS.MotherService.Name,
			ExceptionRate:            entityTS.MotherService.ExceptionRate,
			ResponseDelayRate:        entityTS.MotherService.ResponseDelayRate,
			ResponseDelayDuration:    entityTS.MotherService.ResponseDelayDuration,
			RandomResponseDelayMin:   entityTS.MotherService.RandomResponseDelayMin,
			RandomResponseDelayMax:   entityTS.MotherService.RandomResponseDelayMax,
			Status:                   entityTS.MotherService.Status,
			ServiceDeploymentAddress: entityTS.MotherService.ServiceDeploymentAddress,
			DatabaseName:             entityTS.MotherService.DatabaseName,
			DatabaseTableName:        entityTS.MotherService.DatabaseTableName,
		}
	}()
	resp.TestServiceConfig = func() *TestServiceConfig {
		if entityTS.TestServiceConfig == nil {
			return nil
		}

		return &TestServiceConfig{
			ID:                     entityTS.TestServiceConfig.ID,
			CreatedAt:              entityTS.TestServiceConfig.CreatedAt,
			UpdatedAt:              entityTS.TestServiceConfig.UpdatedAt,
			MaxRequests:            entityTS.TestServiceConfig.MaxRequests,
			MaxDuration:            entityTS.TestServiceConfig.MaxDuration,
			RequestDelayDuration:   entityTS.TestServiceConfig.RequestDelayDuration,
			RandomRequestDelayMin:  entityTS.TestServiceConfig.RandomRequestDelayMin,
			RandomRequestDelayMax:  entityTS.TestServiceConfig.RandomRequestDelayMax,
			FixedTestNumber:        entityTS.TestServiceConfig.FixedTestNumber,
			RandomTestNumberMin:    entityTS.TestServiceConfig.RandomTestNumberMin,
			RandomTestNumberMax:    entityTS.TestServiceConfig.RandomTestNumberMax,
			BadValueRate:           entityTS.TestServiceConfig.BadValueRate,
			NegativeValueRate:      entityTS.TestServiceConfig.NegativeValueRate,
			RealValueRate:          entityTS.TestServiceConfig.RealValueRate,
			ZeroValueRate:          entityTS.TestServiceConfig.ZeroValueRate,
			StringValueRate:        entityTS.TestServiceConfig.StringValueRate,
			LongStringValueRate:    entityTS.TestServiceConfig.LongStringValueRate,
			NullValueRate:          entityTS.TestServiceConfig.NullValueRate,
			DatabaseName:           entityTS.TestServiceConfig.DatabaseName,
			DatabaseTableName:      entityTS.TestServiceConfig.DatabaseTableName,
			IncreaseFixedInput:     entityTS.TestServiceConfig.IncreaseFixedInput,
			ExecNumMultiFixedInput: entityTS.TestServiceConfig.ExecNumMultiFixedInput,
		}
	}()
}
