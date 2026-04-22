package request

import (
	"control-panel-service/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestScenario_ToTestScenarioEntity(t *testing.T) {
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
