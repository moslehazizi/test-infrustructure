package response

import (
	"control-panel-service/internal/domain/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestScenario_FromTestScenarioEntity(t *testing.T) {
	now := time.Now()

	t.Run("success_update_all_fields_from_entity", func(t *testing.T) {
		// Arrange
		entityTS := &entity.TestScenario{
			ID:                  1,
			CreatedAt:           now,
			UpdatedAt:           now,
			StartedAt:           &now,
			Name:                "Full Scenario",
			Status:              "running",
			MaxTestServiceCount: new(int64(100)),
			IncreaseAgentNumber: 5,
			ExecNumMultiAgent:   3,
			Editable:            true,
			NumSteps:            10,
			TestCategory: &entity.TestCategory{
				ID:                     2,
				Name:                   "Category 1",
				Label:                  "cat-1",
				HasMaxTestServiceCount: true,
				HasNumSteps:            false,
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			MotherService: &entity.MotherService{
				ID:                       3,
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "Mother 1",
				ExceptionRate:            0,
				ResponseDelayRate:        0,
				ResponseDelayDuration:    new(100),
				RandomResponseDelayMin:   new(10),
				RandomResponseDelayMax:   new(50),
				Status:                   "active",
				ServiceDeploymentAddress: new("http://localhost:8080"),
				DatabaseName:             "mother_db",
				DatabaseTableName:        "mother_table",
			},
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     4,
				CreatedAt:              now,
				UpdatedAt:              now,
				MaxRequests:            1000,
				MaxDuration:            3600,
				RequestDelayDuration:   new(5),
				RandomRequestDelayMin:  new(1),
				RandomRequestDelayMax:  new(10),
				FixedTestNumber:        new(42),
				RandomTestNumberMin:    new(20),
				RandomTestNumberMax:    new(80),
				BadValueRate:           0,
				NegativeValueRate:      0,
				RealValueRate:          0,
				ZeroValueRate:          0,
				StringValueRate:        0,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "config_db",
				DatabaseTableName:      "config_table",
				IncreaseFixedInput:     2,
				ExecNumMultiFixedInput: 4,
			},
		}

		resp := &TestScenario{}

		// Act
		resp.FromTestScenarioEntity(entityTS)

		// Assert Scalars
		assert.Equal(t, resp.ID, entityTS.ID)
		assert.Equal(t, resp.CreatedAt, entityTS.CreatedAt)
		assert.Equal(t, resp.UpdatedAt, entityTS.UpdatedAt)
		assert.Equal(t, resp.StartedAt, entityTS.StartedAt)
		assert.Equal(t, resp.Name, entityTS.Name)
		assert.Equal(t, resp.Status, entityTS.Status)
		assert.Equal(t, resp.MaxTestServiceCount, entityTS.MaxTestServiceCount)
		assert.Equal(t, resp.IncreaseAgentNumber, entityTS.IncreaseAgentNumber)
		assert.Equal(t, resp.ExecNumMultiAgent, entityTS.ExecNumMultiAgent)
		assert.Equal(t, resp.Editable, entityTS.Editable)
		assert.Equal(t, resp.NumSteps, entityTS.NumSteps)

		// Assert TestCategory
		assert.NotNil(t, entityTS.TestCategory)
		assert.Equal(t, resp.TestCategory.ID, entityTS.TestCategory.ID)
		assert.Equal(t, resp.TestCategory.Name, entityTS.TestCategory.Name)
		assert.Equal(t, resp.TestCategory.Label, entityTS.TestCategory.Label)
		assert.Equal(t, resp.TestCategory.HasMaxTestServiceCount, entityTS.TestCategory.HasMaxTestServiceCount)
		assert.Equal(t, resp.TestCategory.HasNumSteps, entityTS.TestCategory.HasNumSteps)
		assert.Equal(t, resp.TestCategory.CreatedAt, entityTS.TestCategory.CreatedAt)
		assert.Equal(t, resp.TestCategory.UpdatedAt, entityTS.TestCategory.UpdatedAt)

		// Assert MotherService
		assert.NotNil(t, entityTS.MotherService)
		assert.Equal(t, resp.MotherService.ID, entityTS.MotherService.ID)
		assert.Equal(t, resp.MotherService.Name, entityTS.MotherService.Name)
		assert.Equal(t, resp.MotherService.ExceptionRate, entityTS.MotherService.ExceptionRate)
		assert.Equal(t, resp.MotherService.ResponseDelayRate, entityTS.MotherService.ResponseDelayRate)
		assert.Equal(t, resp.MotherService.ResponseDelayDuration, entityTS.MotherService.ResponseDelayDuration)
		assert.Equal(t, resp.MotherService.RandomResponseDelayMin, entityTS.MotherService.RandomResponseDelayMin)
		assert.Equal(t, resp.MotherService.RandomResponseDelayMax, entityTS.MotherService.RandomResponseDelayMax)
		assert.Equal(t, resp.MotherService.Status, entityTS.MotherService.Status)
		assert.Equal(t, resp.MotherService.ServiceDeploymentAddress, entityTS.MotherService.ServiceDeploymentAddress)
		assert.Equal(t, resp.MotherService.DatabaseName, entityTS.MotherService.DatabaseName)
		assert.Equal(t, resp.MotherService.DatabaseTableName, entityTS.MotherService.DatabaseTableName)

		// Assert TestServiceConfig
		assert.NotNil(t, entityTS.TestServiceConfig)
		assert.Equal(t, resp.TestServiceConfig.ID, entityTS.TestServiceConfig.ID)
		assert.Equal(t, resp.TestServiceConfig.MaxRequests, entityTS.TestServiceConfig.MaxRequests)
		assert.Equal(t, resp.TestServiceConfig.MaxDuration, entityTS.TestServiceConfig.MaxDuration)
		assert.Equal(t, resp.TestServiceConfig.BadValueRate, entityTS.TestServiceConfig.BadValueRate)
		assert.Equal(t, resp.TestServiceConfig.DatabaseName, entityTS.TestServiceConfig.DatabaseName)
		assert.Equal(t, resp.TestServiceConfig.ExecNumMultiFixedInput, entityTS.TestServiceConfig.ExecNumMultiFixedInput)
	})

	t.Run("success_nested_structs_are_nil", func(t *testing.T) {
		// Arrange
		entityTS := &entity.TestScenario{
			ID:                1,
			Name:              "Nil Configs Scenario",
			TestCategory:      nil,
			MotherService:     nil,
			TestServiceConfig: nil,
		}

		resp := &TestScenario{}

		// Act
		resp.FromTestScenarioEntity(entityTS)

		// Assert
		assert.Equal(t, resp.ID, entityTS.ID)
		assert.Equal(t, resp.Name, entityTS.Name)

		// Verify nil checks worked
		assert.Nil(t, entityTS.TestCategory)
		assert.Nil(t, entityTS.MotherService)
		assert.Nil(t, entityTS.TestServiceConfig)
	})
}
