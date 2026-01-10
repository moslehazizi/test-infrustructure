package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/postgres/mocks"
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
			MaxTestServiceCount:  nil,
			ExecutionDuration:    nil,
			AutoStepIncreaseRate: nil,
			StoppedAt:            nil,
			StartedAt:            nil,
			RestartedAt:          nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","max_test_service_count","execution_duration","auto_step_increase_rate","stopped_at","restarted_at","started_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err = repo.Create(context.Background(), testScenario)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
