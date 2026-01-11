package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/postgres/mocks"
	"control-panel-service/pkg"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewTestCategoryRepository(t *testing.T) {
	conn := new(mocks.Connection)
	db, _, err := conn.OpenConnection()
	require.NoError(t, err)

	repo := NewTestCategoryRepository(db)
	assert.NotNil(t, repo)

	tr, ok := repo.(*testCategory)
	assert.True(t, ok)
	assert.NotNil(t, tr.db)
}

func TestGetAll(t *testing.T) {
	t.Run("success case with empty result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" ORDER BY id DESC`)).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_increase_rate",
			"created_at",
			"updated_at",
		}))

		repo := NewTestCategoryRepository(db)
		result, err := repo.GetAll(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
	t.Run("success case with some results", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		expected := []entity.TestCategory{
			{
				ID:                      1,
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
				Name:                    "load",
				Label:                   "Load Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
			},
			{
				ID:                      2,
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
				Name:                    "smoke",
				Label:                   "Smoke Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" ORDER BY id DESC`)).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_increase_rate",
			"created_at",
			"updated_at",
		}).AddRow(
			expected[0].ID,
			expected[0].Name,
			expected[0].Label,
			expected[0].HasMaxTestServiceCount,
			expected[0].HasExecutionDuration,
			expected[0].HasAutoStepIncreaseRate,
			expected[0].CreatedAt,
			expected[0].UpdatedAt,
		).AddRow(
			expected[1].ID,
			expected[1].Name,
			expected[1].Label,
			expected[1].HasMaxTestServiceCount,
			expected[1].HasExecutionDuration,
			expected[1].HasAutoStepIncreaseRate,
			expected[1].CreatedAt,
			expected[1].UpdatedAt,
		))

		repo := NewTestCategoryRepository(db)
		result, err := repo.GetAll(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		assert.True(t, reflect.DeepEqual(result, expected))
	})
	t.Run("error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" ORDER BY id DESC`)).WillReturnError(errors.New("database error"))

		repo := NewTestCategoryRepository(db)
		result, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTestCategory_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestCategoryRepository(db)

		now := time.Now()

		expectedTestCategory := &entity.TestCategory{
			ID:                      1,
			CreatedAt:               now,
			UpdatedAt:               now,
			Name:                    "load",
			Label:                   "my load",
			HasMaxTestServiceCount:  true,
			HasExecutionDuration:    true,
			HasAutoStepIncreaseRate: true,
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_categories" WHERE "test_categories"."id" = $1 ORDER BY "test_categories"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "name", "label",
				"has_max_test_service_count", "has_execution_duration",
				"has_auto_step_increase_rate",
			}).
				AddRow(
					expectedTestCategory.ID,
					expectedTestCategory.CreatedAt,
					expectedTestCategory.UpdatedAt,
					expectedTestCategory.Name,
					expectedTestCategory.Label,
					expectedTestCategory.HasAutoStepIncreaseRate,
					expectedTestCategory.HasExecutionDuration,
					expectedTestCategory.HasMaxTestServiceCount,
				))

		result, err := repo.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestCategory.ID, result.ID)
		assert.Equal(t, expectedTestCategory.Name, result.Name)
		assert.Equal(t, expectedTestCategory.Label, result.Label)
		assert.Equal(t, expectedTestCategory.HasAutoStepIncreaseRate, result.HasAutoStepIncreaseRate)
		assert.Equal(t, expectedTestCategory.HasExecutionDuration, result.HasExecutionDuration)
		assert.Equal(t, expectedTestCategory.HasMaxTestServiceCount, result.HasMaxTestServiceCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - not found", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestCategoryRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_categories" WHERE "test_categories"."id" = $1 ORDER BY "test_categories"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ErrTestCategoryNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - database error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestCategoryRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_categories" WHERE "test_categories"."id" = $1 ORDER BY "test_categories"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnError(errors.New("error happened"))

		result, err := repo.GetByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "error happened")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
