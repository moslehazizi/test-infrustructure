package request

import (
	"control-panel-service/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestScenario_TestScenarioUpdateRequestToTestScenarioEntity(t *testing.T) {
	t.Run("success_update_all_fields", func(t *testing.T) {
		maxSvcCount := int64(5)
		motherSvc := &entity.MotherService{ID: 10}

		ts := &entity.TestScenario{
			NumSteps: 0,
		}

		req := &TestScenarioUpdateRequest{
			Name:                "Updated Scenario",
			MotherServiceID:     10,
			IncreaseAgentNumber: 2,
			ExecNumMultiAgent:   3,
			MaxTestServiceCount: &maxSvcCount,
			NumSteps:            10,
		}

		req.ToTestScenarioEntity(ts, motherSvc)

		assert.Equal(t, req.Name, ts.Name)
		assert.Equal(t, req.MotherServiceID, ts.MotherServiceID)
		assert.Equal(t, motherSvc, ts.MotherService)
		assert.Equal(t, req.IncreaseAgentNumber, ts.IncreaseAgentNumber)
		assert.Equal(t, req.ExecNumMultiAgent, ts.ExecNumMultiAgent)
		assert.Equal(t, entity.ScenarioStatusReady, ts.Status)
		assert.Equal(t, req.MaxTestServiceCount, ts.MaxTestServiceCount)
		assert.Equal(t, req.NumSteps, ts.NumSteps)
	})

	t.Run("success_update_skip_optional_fields_when_nil_or_less_than_one", func(t *testing.T) {
		initialMaxCount := int64(2)
		ts := &entity.TestScenario{
			MaxTestServiceCount: &initialMaxCount,
			NumSteps:            5,
		}

		req := &TestScenarioUpdateRequest{
			Name:                "Updated Scenario",
			MaxTestServiceCount: nil,
			NumSteps:            0,
		}

		motherSvc := &entity.MotherService{}
		req.ToTestScenarioEntity(ts, motherSvc)

		assert.Equal(t, "Updated Scenario", ts.Name)
		assert.Equal(t, &initialMaxCount, ts.MaxTestServiceCount)
		assert.Equal(t, int64(5), ts.NumSteps)
		assert.Equal(t, entity.ScenarioStatusReady, ts.Status)
	})
}

func TestTestServiceConfig_ToTestServiceConfigEntity(t *testing.T) {
	t.Run("success_update_all_fields_from_request", func(t *testing.T) {
		sampleInt1 := 10
		sampleInt2 := 20
		sampleInt3 := 30

		tsc := &entity.TestServiceConfig{}

		req := &TestServiceConfigRequest{
			MaxRequests:            100,
			MaxDuration:            5000,
			RequestDelayDuration:   &sampleInt1,
			RandomRequestDelayMin:  &sampleInt1,
			RandomRequestDelayMax:  &sampleInt2,
			FixedTestNumber:        &sampleInt3,
			RandomTestNumberMin:    &sampleInt1,
			RandomTestNumberMax:    &sampleInt2,
			BadValueRate:           1,
			NegativeValueRate:      2,
			ZeroValueRate:          3,
			StringValueRate:        4,
			RealValueRate:          5,
			LongStringValueRate:    6,
			NullValueRate:          7,
			DatabaseName:           "test_db",
			DatabaseTableName:      "test_table",
			IncreaseFixedInput:     2,
			ExecNumMultiFixedInput: 3,
		}

		req.ToTestServiceConfigEntity(tsc)

		assert.Equal(t, req.MaxRequests, tsc.MaxRequests)
		assert.Equal(t, int64(req.MaxDuration), tsc.MaxDuration)
		assert.Equal(t, req.RequestDelayDuration, tsc.RequestDelayDuration)
		assert.Equal(t, req.RandomRequestDelayMin, tsc.RandomRequestDelayMin)
		assert.Equal(t, req.RandomRequestDelayMax, tsc.RandomRequestDelayMax)
		assert.Equal(t, req.FixedTestNumber, tsc.FixedTestNumber)
		assert.Equal(t, req.RandomTestNumberMin, tsc.RandomTestNumberMin)
		assert.Equal(t, req.RandomTestNumberMax, tsc.RandomTestNumberMax)
		assert.Equal(t, req.BadValueRate, tsc.BadValueRate)
		assert.Equal(t, req.NegativeValueRate, tsc.NegativeValueRate)
		assert.Equal(t, req.ZeroValueRate, tsc.ZeroValueRate)
		assert.Equal(t, req.StringValueRate, tsc.StringValueRate)
		assert.Equal(t, req.RealValueRate, tsc.RealValueRate)
		assert.Equal(t, req.LongStringValueRate, tsc.LongStringValueRate)
		assert.Equal(t, req.NullValueRate, tsc.NullValueRate)
		assert.Equal(t, req.DatabaseName, tsc.DatabaseName)
		assert.Equal(t, req.DatabaseTableName, tsc.DatabaseTableName)
		assert.Equal(t, req.IncreaseFixedInput, tsc.IncreaseFixedInput)
		assert.Equal(t, req.ExecNumMultiFixedInput, tsc.ExecNumMultiFixedInput)
	})
}

func TestTestScenario_TestScenarioToTestScenarioEntity(t *testing.T) {
	t.Run("success_update_all_fields_from_request", func(t *testing.T) {
		// Arrange
		ts := &TestScenario{
			Name:                "Load Test Scenario",
			TestCategoryID:      10,
			MotherServiceID:     20,
			MaxTestServiceCount: new(int64(100)),
			NumSteps:            5,
			IncreaseAgentNumber: 2,
			ExecNumMultiAgent:   3,
			Config: &TestServiceConfigRequest{
				MaxRequests:            500,
				MaxDuration:            3600,
				RequestDelayDuration:   new(10),
				RandomRequestDelayMin:  new(1),
				RandomRequestDelayMax:  new(5),
				FixedTestNumber:        new(42),
				RandomTestNumberMin:    new(10),
				RandomTestNumberMax:    new(20),
				BadValueRate:           0,
				NegativeValueRate:      0,
				ZeroValueRate:          0,
				StringValueRate:        0,
				RealValueRate:          0,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "test_db",
				DatabaseTableName:      "test_table",
				IncreaseFixedInput:     2,
				ExecNumMultiFixedInput: 3,
			},
		}

		entityTS := &entity.TestScenario{}

		// Act
		ts.ToTestScenarioEntity(entityTS)

		// Assert
		assert.Equal(t, ts.Name, entityTS.Name)
		assert.Equal(t, ts.TestCategoryID, entityTS.TestCategoryID)
		assert.Equal(t, ts.MotherServiceID, entityTS.MotherServiceID)
		assert.Equal(t, ts.MaxTestServiceCount, entityTS.MaxTestServiceCount)
		assert.Equal(t, ts.NumSteps, entityTS.NumSteps)
		assert.Equal(t, ts.IncreaseAgentNumber, entityTS.IncreaseAgentNumber)
		assert.Equal(t, ts.ExecNumMultiAgent, entityTS.ExecNumMultiAgent)

		// Assert Config mapping
		assert.NotNil(t, entityTS.TestServiceConfig)
		assert.Equal(t, ts.Config.MaxRequests, entityTS.TestServiceConfig.MaxRequests)
		assert.Equal(t, int64(ts.Config.MaxDuration), entityTS.TestServiceConfig.MaxDuration)
		assert.Equal(t, ts.Config.RequestDelayDuration, entityTS.TestServiceConfig.RequestDelayDuration)
		assert.Equal(t, ts.Config.RandomRequestDelayMin, entityTS.TestServiceConfig.RandomRequestDelayMin)
		assert.Equal(t, ts.Config.RandomRequestDelayMax, entityTS.TestServiceConfig.RandomRequestDelayMax)
		assert.Equal(t, ts.Config.FixedTestNumber, entityTS.TestServiceConfig.FixedTestNumber)
		assert.Equal(t, ts.Config.RandomTestNumberMin, entityTS.TestServiceConfig.RandomTestNumberMin)
		assert.Equal(t, ts.Config.RandomTestNumberMax, entityTS.TestServiceConfig.RandomTestNumberMax)
		assert.Equal(t, ts.Config.BadValueRate, entityTS.TestServiceConfig.BadValueRate)
		assert.Equal(t, ts.Config.NegativeValueRate, entityTS.TestServiceConfig.NegativeValueRate)
		assert.Equal(t, ts.Config.ZeroValueRate, entityTS.TestServiceConfig.ZeroValueRate)
		assert.Equal(t, ts.Config.StringValueRate, entityTS.TestServiceConfig.StringValueRate)
		assert.Equal(t, ts.Config.RealValueRate, entityTS.TestServiceConfig.RealValueRate)
		assert.Equal(t, ts.Config.LongStringValueRate, entityTS.TestServiceConfig.LongStringValueRate)
		assert.Equal(t, ts.Config.NullValueRate, entityTS.TestServiceConfig.NullValueRate)
		assert.Equal(t, ts.Config.DatabaseName, entityTS.TestServiceConfig.DatabaseName)
		assert.Equal(t, ts.Config.DatabaseTableName, entityTS.TestServiceConfig.DatabaseTableName)
		assert.Equal(t, ts.Config.IncreaseFixedInput, entityTS.TestServiceConfig.IncreaseFixedInput)
		assert.Equal(t, ts.Config.ExecNumMultiFixedInput, entityTS.TestServiceConfig.ExecNumMultiFixedInput)
	})

	t.Run("success_config_is_nil", func(t *testing.T) {
		// Arrange
		ts := &TestScenario{
			Name:           "Scenario without config",
			TestCategoryID: 5,
			Config:         nil, // Explicitly nil
		}

		entityTS := &entity.TestScenario{}

		// Act
		ts.ToTestScenarioEntity(entityTS)

		// Assert
		assert.Equal(t, ts.Name, entityTS.Name)
		assert.Equal(t, ts.TestCategoryID, entityTS.TestCategoryID)

		// This is the main assertion for this test block
		assert.Nil(t, entityTS.TestServiceConfig, "TestServiceConfig should be nil when ts.Config is nil")
	})
}
