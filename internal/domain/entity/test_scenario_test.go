package entity

import (
	"control-panel-service/pkg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestScenarioValidation(t *testing.T) {
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
			Status:              ScenarioStatusPending,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
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
			ID:              uint64(1),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			DeletedAt:       nil,
			Name:            "load1",
			TestCategoryID:  uint64(4),
			MotherServiceID: uint64(5),
			Status:          ScenarioStatusPending,
			NumSteps:        int64(1),
		}
		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasNumSteps: false,
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
			Status:              ScenarioStatusPending,
			MaxTestServiceCount: &sampleIntLessThanOne,
			NumSteps:            int64(2),
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
			ID:              uint64(1),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			DeletedAt:       nil,
			Name:            "load1",
			TestCategoryID:  uint64(4),
			MotherServiceID: uint64(5),
			Status:          ScenarioStatusPending,
			NumSteps:        int64(2),
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
			Status:              ScenarioStatusPending,
			MaxTestServiceCount: &sampleInt,
			NumSteps:            int64(2),
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
			NumSteps: 0,
		}
		// ...
		err := testSci.Validate(&TestCategory{HasNumSteps: true})
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNumStepsNotSet)
	})
	t.Run("failed_case_num_steps_set_to_negative_value", func(t *testing.T) {
		testSci := TestScenario{
			// ... other fields ...
			NumSteps: -2,
		}
		// ...
		err := testSci.Validate(&TestCategory{HasNumSteps: true})
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNumStepsNotSet)
	})
	t.Run("failed_case_num_steps_no_need_to_set", func(t *testing.T) {
		testSci := TestScenario{
			// ... other fields ...
			NumSteps: 2,
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
		Status:          ScenarioStatusPending,
		NumSteps:        int64(2),
	}

	name := testSci.TableName()

	assert.Equal(t, name, "test_scenarios")
}
