package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_testService_GetRunningByScenario(t *testing.T) {
t.Run("success_case_with_empty_result_and_no_limit", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_services" WHERE "status" = $1 AND "test_services"."deleted_at" IS NULL`)).
			WithArgs(entity.TestServiceStatusRunning).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
			}))

		repo := NewTestServiceRepository(db)
		result, err := repo.GetRunningByScenario(context.Background(), uint64(1), int64(-1))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

t.Run("success_case_with_some_results_and_no_limit", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_services" WHERE "status" = $1 AND "test_services"."deleted_at" IS NULL`)).
			WithArgs(entity.TestServiceStatusRunning).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
			}).AddRow(1).AddRow(2).AddRow(3))

		repo := NewTestServiceRepository(db)
		result, err := repo.GetRunningByScenario(context.Background(), uint64(1), int64(-1))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
t.Run("success_case_with_some_results_and_applied_limit", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_services" WHERE "status" = $1 AND "test_services"."deleted_at" IS NULL LIMIT $2`)).
			WithArgs(entity.TestServiceStatusRunning, 3).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
			}).AddRow(1).AddRow(2).AddRow(3))

		repo := NewTestServiceRepository(db)
		result, err := repo.GetRunningByScenario(context.Background(), uint64(1), int64(3))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
t.Run("error_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_services" WHERE "status" = $1 AND "test_services"."deleted_at" IS NULL LIMIT $2`)).
			WithArgs(entity.TestServiceStatusRunning, 3).
			WillReturnError(errors.New("something went wrong"))

		repo := NewTestServiceRepository(db)
		result, err := repo.GetRunningByScenario(context.Background(), uint64(1), int64(3))
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_testService_GetCountAllRunningByScenario(t *testing.T) {
t.Run("error_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, _, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestServiceRepository(db)
		result, err := repo.GetCountAllRunningByScenario(context.Background(), uint64(1))
		assert.Error(t, err)
		assert.Equal(t, result, int64(0))
	})

t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "test_services" WHERE "status" = $1 AND "test_services"."deleted_at" IS NULL`)).
			WithArgs(entity.TestServiceStatusRunning).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		repo := NewTestServiceRepository(db)
		result, err := repo.GetCountAllRunningByScenario(context.Background(), uint64(1))
		assert.NoError(t, err)
		assert.Equal(t, result, int64(3))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
