package entity

import (
	"control-panel-service/pkg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestScenarioValidation(t *testing.T) {
	t.Run("success case - test category config matches inputs", func(t *testing.T) {
		sampleInt := 2
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
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Nil(t, err)
	})

	t.Run("success case - nullable values not set", func(t *testing.T) {
		testSci := TestScenario{
			ID:              uint64(1),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			DeletedAt:       nil,
			Name:            "load1",
			TestCategoryID:  uint64(4),
			MotherServiceID: uint64(5),
			Status:          ScenarioStatusPending,
		}
		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  false,
		}

		err := testSci.Validate(TestCategory)

		assert.Nil(t, err)
	})

	t.Run("failed case - max test service count less than one", func(t *testing.T) {
		sampleIntLessThanOne := -1
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
		}
		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
	})

	t.Run("failed case - execution duration should be more than 1", func(t *testing.T) {
		sampleIntLessThanOne := -1
		testSci := TestScenario{
			ID:                uint64(1),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
			Name:              "load1",
			TestCategoryID:    uint64(4),
			MotherServiceID:   uint64(5),
			Status:            ScenarioStatusPending,
			ExecutionDuration: &sampleIntLessThanOne,
		}
		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrExecutionDurationLessThanOne)
	})

	t.Run("failed case - auto step change rate should be more than 1", func(t *testing.T) {
		sampleIntLessThanOne := -1
		testSci := TestScenario{
			ID:                 uint64(1),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			DeletedAt:          nil,
			Name:               "load1",
			TestCategoryID:     uint64(4),
			MotherServiceID:    uint64(5),
			Status:             ScenarioStatusPending,
			AutoStepChangeRate: &sampleIntLessThanOne,
		}
		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrAutoStepChangeRateLessThanOne)
	})

	t.Run("failed case - max test service count not set", func(t *testing.T) {
		sampleInt := 2
		testSci := TestScenario{
			ID:                 uint64(1),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			DeletedAt:          nil,
			Name:               "load1",
			TestCategoryID:     uint64(4),
			MotherServiceID:    uint64(5),
			Status:             ScenarioStatusPending,
			ExecutionDuration:  &sampleInt,
			AutoStepChangeRate: &sampleInt,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("failed case - no need to max test service count", func(t *testing.T) {
		sampleInt := 2
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
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNoNeedMaxTestServiceCount)
	})

	t.Run("failed case - execution duration not set", func(t *testing.T) {
		sampleInt := 2
		testSci := TestScenario{
			ID:                 uint64(1),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			DeletedAt:          nil,
			Name:               "load1",
			TestCategoryID:     uint64(4),
			MotherServiceID:    uint64(5),
			Status:             ScenarioStatusPending,
			AutoStepChangeRate: &sampleInt,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   true,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrExecutionDurationNotSet)
	})

	t.Run("failed case - no need execution duration", func(t *testing.T) {
		sampleInt := 2
		testSci := TestScenario{
			ID:                 uint64(1),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			DeletedAt:          nil,
			Name:               "load1",
			TestCategoryID:     uint64(4),
			MotherServiceID:    uint64(5),
			Status:             ScenarioStatusPending,
			ExecutionDuration:  &sampleInt,
			AutoStepChangeRate: &sampleInt,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNoNeedExecutionDuration)
	})

	t.Run("failed case - auto step change not set", func(t *testing.T) {
		testSci := TestScenario{
			ID:              uint64(1),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			DeletedAt:       nil,
			Name:            "load1",
			TestCategoryID:  uint64(4),
			MotherServiceID: uint64(5),
			Status:          ScenarioStatusPending,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  true,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrAutoStepChangeNotSet)
	})
	t.Run("failed case - no need to auto step change", func(t *testing.T) {
		sampleInt := 1
		testSci := TestScenario{
			ID:                 uint64(1),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			DeletedAt:          nil,
			Name:               "load1",
			TestCategoryID:     uint64(4),
			MotherServiceID:    uint64(5),
			Status:             ScenarioStatusPending,
			AutoStepChangeRate: &sampleInt,
		}

		TestCategory := &TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasExecutionDuration:   false,
			HasAutoStepChangeRate:  false,
		}

		err := testSci.Validate(TestCategory)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNoNeedAutoStepChange)
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
	}

	name := testSci.TableName()

	assert.Equal(t, name, "test_scenarios")
}
