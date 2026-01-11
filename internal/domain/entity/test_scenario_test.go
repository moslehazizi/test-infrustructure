package entity

import (
	"control-panel-service/pkg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestScenarioValidation(t *testing.T) {
	t.Run("success case - all values set", func(t *testing.T) {
		sampleInt := 2
		testSci := TestScenario{
			ID:                   uint64(1),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			DeletedAt:            nil,
			Name:                 "load1",
			TestCategoryID:       uint64(4),
			MotherServiceID:      uint64(5),
			Status:               ScenarioStatusPending,
			MaxTestServiceCount:  &sampleInt,
			ExecutionDuration:    &sampleInt,
			AutoStepIncreaseRate: &sampleInt,
		}

		err := testSci.Validate()

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

		err := testSci.Validate()

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

		err := testSci.Validate()

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

		err := testSci.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrExecutionDurationLessThanOne)
	})

	t.Run("failed case - auto step increase rate should be more than 1", func(t *testing.T) {
		sampleIntLessThanOne := -1
		testSci := TestScenario{
			ID:                   uint64(1),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			DeletedAt:            nil,
			Name:                 "load1",
			TestCategoryID:       uint64(4),
			MotherServiceID:      uint64(5),
			Status:               ScenarioStatusPending,
			AutoStepIncreaseRate: &sampleIntLessThanOne,
		}

		err := testSci.Validate()

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrAutoStepIncreaseRateLessThanOne)
	})
}
