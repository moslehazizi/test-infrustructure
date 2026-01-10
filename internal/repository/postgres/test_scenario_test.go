package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/postgres/mocks"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestScenarioRepository_New(t *testing.T) {
	mockConn := new(mocks.Connection)
	db, _, err := mockConn.OpenConnection()
	require.NoError(t, err)

	repo := NewTestScenarioRepository(db)
	assert.NotNil(t, repo)

	ts, ok := repo.(*testScenario)
	assert.True(t, ok)
	assert.NotNil(t, ts.db)
}

func TestTestScenarioRepository_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)
		now := time.Now()

		testScenario := &entity.TestScenario{
			CreatedAt:            now,
			UpdatedAt:            now,
			DeletedAt:            nil,
			Name:                 "load1",
			TestCategoryID:       uint64(2),
			MotherServiceID:      uint64(1),
			Status:               entity.ScenarioStatus(entity.ScenarioStatusPending),
			MaxTestServiceCount:  nil,
			ExecutionDuration:    nil,
			AutoStepIncreaseRate: nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","max_test_service_count","execution_duration","auto_step_increase_rate") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
			WithArgs(
				testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				nil, nil, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err = repo.Create(context.Background(), testScenario)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)
		now := time.Now()

		testScenario := &entity.TestScenario{
			CreatedAt:            now,
			UpdatedAt:            now,
			DeletedAt:            nil,
			Name:                 "load1",
			TestCategoryID:       uint64(2),
			MotherServiceID:      uint64(1),
			Status:               entity.ScenarioStatus(entity.ScenarioStatusPending),
			MaxTestServiceCount:  nil,
			ExecutionDuration:    nil,
			AutoStepIncreaseRate: nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","max_test_service_count","execution_duration","auto_step_increase_rate") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
			WithArgs(testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				nil, nil, nil).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		err = repo.Create(context.Background(), testScenario)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create test scenario record")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
