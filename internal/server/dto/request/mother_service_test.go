package request

import (
	"control-panel-service/internal/domain/entity"
	"reflect"
	"testing"
)

func TestToMotherServiceEntity(t *testing.T) {
	duration := 100
	minDelay := 10
	maxDelay := 500

	t.Run("Nil request", func(t *testing.T) {
		var req *MotherService
		got := req.ToMotherServiceEntity()
		if got != nil {
			t.Errorf("ToMotherServiceEntity() = %v, want nil", got)
		}
	})

	t.Run("All fields populated", func(t *testing.T) {
		req := &MotherService{
			Name:                   "Service A",
			ExceptionRate:          1,
			ResponseDelayRate:      1,
			ResponseDelayDuration:  &duration,
			RandomResponseDelayMin: &minDelay,
			RandomResponseDelayMax: &maxDelay,
			DatabaseName:           "db_main",
			DatabaseTableName:      "table_a",
		}

		expected := &entity.MotherService{
			Name:                   "Service A",
			ExceptionRate:          1,
			ResponseDelayRate:      1,
			ResponseDelayDuration:  &duration,
			RandomResponseDelayMin: &minDelay,
			RandomResponseDelayMax: &maxDelay,
			DatabaseName:           "db_main",
			DatabaseTableName:      "table_a",
		}

		got := req.ToMotherServiceEntity()
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("ToMotherServiceEntity() = %v, want %v", got, expected)
		}
	})

	t.Run("With nil pointers", func(t *testing.T) {
		req := &MotherService{
			Name:                   "Service B",
			ResponseDelayDuration:  nil,
			RandomResponseDelayMin: nil,
			RandomResponseDelayMax: nil,
		}

		expected := &entity.MotherService{
			Name:                   "Service B",
			ResponseDelayDuration:  nil,
			RandomResponseDelayMin: nil,
			RandomResponseDelayMax: nil,
		}

		got := req.ToMotherServiceEntity()
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("ToMotherServiceEntity() = %v, want %v", got, expected)
		}
	})
}
