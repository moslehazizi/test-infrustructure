package response

import (
	"control-panel-service/internal/domain/entity"
	"reflect"
	"testing"
	"time"
)

func TestFromMotherServiceEntity(t *testing.T) {
	now := time.Now()
	duration := 100
	minDelay := 10
	maxDelay := 500

	t.Run("All_fields_populated", func(t *testing.T) {
		entityMother := &entity.MotherService{
			ID:                       1,
			CreatedAt:                now,
			UpdatedAt:                now,
			Name:                     "Service A",
			ExceptionRate:            1,
			ResponseDelayRate:        1,
			ResponseDelayDuration:    &duration,
			RandomResponseDelayMin:   &minDelay,
			RandomResponseDelayMax:   &maxDelay,
			Status:                   "ACTIVE",
			ServiceDeploymentAddress: new("127.0.0.1:8080"),
			DatabaseName:             "db_main",
			DatabaseTableName:        "table_a",
		}

		result := &MotherService{}
		result.FromMotherServiceEntity(entityMother)

		expected := &MotherService{
			ID:                       1,
			CreatedAt:                now,
			UpdatedAt:                now,
			Name:                     "Service A",
			ExceptionRate:            1,
			ResponseDelayRate:        1,
			ResponseDelayDuration:    &duration,
			RandomResponseDelayMin:   &minDelay,
			RandomResponseDelayMax:   &maxDelay,
			Status:                   "ACTIVE",
			ServiceDeploymentAddress: new("127.0.0.1:8080"),
			DatabaseName:             "db_main",
			DatabaseTableName:        "table_a",
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("FromMotherServiceEntity() result = %v, want %v", result, expected)
		}
	})

	t.Run("With_nil_pointers", func(t *testing.T) {
		entityMother := &entity.MotherService{
			ID:                     2,
			Name:                   "Service B",
			ResponseDelayDuration:  nil,
			RandomResponseDelayMin: nil,
			RandomResponseDelayMax: nil,
		}

		result := &MotherService{}
		result.FromMotherServiceEntity(entityMother)

		expected := &MotherService{
			ID:                     2,
			Name:                   "Service B",
			ResponseDelayDuration:  nil,
			RandomResponseDelayMin: nil,
			RandomResponseDelayMax: nil,
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("FromMotherServiceEntity() result = %v, want %v", result, expected)
		}
	})
}
