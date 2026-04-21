package entity

import (
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/pkg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestScenarioValidation(t *testing.T) {
	t.Run("success_case_execution_number_multi_agent_and_increase_number_both_can_be_zero", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   0,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.NoError(t, err)
	})
	t.Run("failed_case_execution_number_multi_agent_and_increase_number_both_or_neither_sould_be_zero", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   0,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMultiAgentConfigNotTrue)
	})
	t.Run("failed_case_execution_number_multi_agent_and_increase_number_both_or_neither_sould_be_zero", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMultiAgentConfigNotTrue)
	})
	t.Run("success_case_execution_number_multi_agent_can_be_than_zero", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   0,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.NoError(t, err)
	})
	t.Run("failed_case_execution_number_multi_agent_should_be_more_than_zero", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   -1,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrExecNumMultiAgentShouldBePositive)
	})
	t.Run("failed_case_increase_agent_number_not_be_negative", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: -2,
			ExecNumMultiAgent:   3,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrIncreaseAgentNumNotBeNegative)
	})
	t.Run("success_case_test_category_config_matches_inputs", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Nil(t, err)
	})
	t.Run("success_case_nullable_values_not_set", func(t *testing.T) {
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			NumSteps:            int64(1),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasNumSteps:            false,
		}

		err := testSci.Validate(TestCategory)

		assert.Nil(t, err)
	})
	t.Run("failed_case_max_test_service_count_less_than_one", func(t *testing.T) {
		sampleIntLessThanOne := int64(-1)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleIntLessThanOne,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
	})
	t.Run("failed_case_max_test_service_count_not_set", func(t *testing.T) {
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})
	t.Run("failed_case_no_need_to_max_test_service_count", func(t *testing.T) {
		sampleInt := int64(2)
		testSci := TestScenario{
			ID:                  uint64(1),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(4),
			MotherServiceID:     uint64(5),
			Status:              ScenarioStatusReady,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasNumSteps:            true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNoNeedMaxTestServiceCount)
	})
	t.Run("failed_case_num_steps_not_set_(zero)", func(t *testing.T) {
		testSci := TestScenario{
			// ... other fields ...
			NumSteps:            0,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		// ...
		err := testSci.Validate(&TestCategory{HasNumSteps: true})
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNumStepsNotSet)
	})
	t.Run("failed_case_num_steps_set_to_negative_value", func(t *testing.T) {
		testSci := TestScenario{
			// ... other fields ...
			NumSteps:            -2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		// ...
		err := testSci.Validate(&TestCategory{HasNumSteps: true})
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNumStepsNotSet)
	})
	t.Run("failed_case_num_steps_no_need_to_set", func(t *testing.T) {
		testSci := TestScenario{
			// ... other fields ...
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		// ...
		err := testSci.Validate(&TestCategory{HasNumSteps: false})
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNumStepsShouldBeOne)
	})
}

func TestTableName(t *testing.T) {
	testSci := TestScenario{
		ID:              uint64(1),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
		Name:            "load1",
		TestCategoryID:  uint64(4),
		MotherServiceID: uint64(5),
		Status:          ScenarioStatusReady,
		NumSteps:        int64(2),
	}

	name := testSci.TableName()

	assert.Equal(t, name, "test_scenarios")
}

func TestTestScenario_ApplyUpdateRequest(t *testing.T) {
	t.Run("success_update_all_fields", func(t *testing.T) {
		maxSvcCount := int64(5)
		motherSvc := &MotherService{ID: 10}

		ts := &TestScenario{
			NumSteps: 0,
		}

		req := &request.TestScenarioUpdateRequest{
			Name:                "Updated Scenario",
			MotherServiceID:     10,
			IncreaseAgentNumber: 2,
			ExecNumMultiAgent:   3,
			MaxTestServiceCount: &maxSvcCount,
			NumSteps:            10,
		}

		ts.ApplyUpdateRequest(req, motherSvc)

		assert.Equal(t, req.Name, ts.Name)
		assert.Equal(t, req.MotherServiceID, ts.MotherServiceID)
		assert.Equal(t, motherSvc, ts.MotherService)
		assert.Equal(t, req.IncreaseAgentNumber, ts.IncreaseAgentNumber)
		assert.Equal(t, req.ExecNumMultiAgent, ts.ExecNumMultiAgent)
		assert.Equal(t, ScenarioStatusReady, ts.Status)
		assert.Equal(t, req.MaxTestServiceCount, ts.MaxTestServiceCount)
		assert.Equal(t, req.NumSteps, ts.NumSteps)
	})

	t.Run("success_update_skip_optional_fields_when_nil_or_less_than_one", func(t *testing.T) {
		initialMaxCount := int64(2)
		ts := &TestScenario{
			MaxTestServiceCount: &initialMaxCount,
			NumSteps:            5,
		}

		req := &request.TestScenarioUpdateRequest{
			Name:                "Updated Scenario",
			MaxTestServiceCount: nil,
			NumSteps:            0,
		}

		motherSvc := &MotherService{}
		ts.ApplyUpdateRequest(req, motherSvc)

		assert.Equal(t, "Updated Scenario", ts.Name)
		assert.Equal(t, &initialMaxCount, ts.MaxTestServiceCount)
		assert.Equal(t, int64(5), ts.NumSteps)
		assert.Equal(t, ScenarioStatusReady, ts.Status)
	})
}
