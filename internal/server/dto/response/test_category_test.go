package response

import (
	"control-panel-service/internal/domain/entity"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestCategory_FromTestCategoryEntity(t *testing.T) {
	now := time.Now()

	t.Run("nil_input", func(t *testing.T) {
		var entityCat *entity.TestCategory

		var result *TestCategory
		result.FromTestCategoryEntity(entityCat)

		assert.Nil(t, result)
	})

	t.Run("All_fields_populated", func(t *testing.T) {
		entityCategory := &entity.TestCategory{
			ID:                     uint64(1),
			CreatedAt:              now,
			UpdatedAt:              now,
			Name:                   "Category Name",
			Label:                  "Category Label",
			HasMaxTestServiceCount: true,
			HasNumSteps:            false,
			Active:                 true,
		}

		result := &TestCategory{}
		result.FromTestCategoryEntity(entityCategory)

		expected := &TestCategory{
			ID:                     uint64(1),
			CreatedAt:              now,
			UpdatedAt:              now,
			Name:                   "Category Name",
			Label:                  "Category Label",
			HasMaxTestServiceCount: true,
			HasNumSteps:            false,
			Active:                 true,
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}
