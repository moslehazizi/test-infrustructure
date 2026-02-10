package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestServiceConfigRepo_Init(t *testing.T) {
	mockConn := new(mocks.Connection)
	db, _, err := mockConn.OpenConnection()
	require.NoError(t, err)

	repo := NewTestServiceConfigRepository(db)
	assert.NotNil(t, repo)

	ts, ok := repo.(*testServiceConfig)
	assert.True(t, ok)
	assert.NotNil(t, ts.db)
}

func TestTestServiceConfigRepo_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestServiceConfigRepository(db)
		now := time.Now()

		databaseName := "test_db"
		databaseTableName := "test_table"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:    uint64(1),
			CreatedAt:         now,
			UpdatedAt:         now,
			MaxRequests:       12,
			MaxDuration:       10,
			DatabaseName:      databaseName,
			DatabaseTableName: databaseTableName,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_service_configs" ("test_scenario_id","max_requests","max_duration","request_delay_duration","random_request_delay_min","random_request_delay_max","fixed_test_number","random_test_number_min","random_test_number_max","bad_value_rate","negative_value_rate","real_value_rate","zero_value_rate","string_value_rate","long_string_value_rate","null_value_rate","created_at","updated_at","database_name","database_table_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20) RETURNING "id"`)).
			WithArgs(
				testServiceConfig.TestScenarioID,
				testServiceConfig.MaxRequests,
				testServiceConfig.MaxDuration,
				nil, nil, nil, nil, nil, nil,
				0, 0, 0, 0, 0, 0, 0,
				testServiceConfig.CreatedAt,
				testServiceConfig.UpdatedAt,
				testServiceConfig.DatabaseName,
				testServiceConfig.DatabaseTableName,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err = repo.Create(context.Background(), testServiceConfig)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - db error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestServiceConfigRepository(db)
		now := time.Now()

		databaseName := "test_db"
		databaseTableName := "test_table"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:    uint64(1),
			CreatedAt:         now,
			UpdatedAt:         now,
			MaxRequests:       12,
			MaxDuration:       10,
			DatabaseName:      databaseName,
			DatabaseTableName: databaseTableName,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_service_configs" ("test_scenario_id","max_requests","max_duration","request_delay_duration","random_request_delay_min","random_request_delay_max","fixed_test_number","random_test_number_min","random_test_number_max","bad_value_rate","negative_value_rate","real_value_rate","zero_value_rate","string_value_rate","long_string_value_rate","null_value_rate","created_at","updated_at","database_name","database_table_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20) RETURNING "id"`)).
			WithArgs(
				testServiceConfig.TestScenarioID,
				testServiceConfig.MaxRequests,
				testServiceConfig.MaxDuration,
				nil, nil, nil, nil, nil, nil,
				0, 0, 0, 0, 0, 0, 0,
				testServiceConfig.CreatedAt,
				testServiceConfig.UpdatedAt,
				testServiceConfig.DatabaseName,
				testServiceConfig.DatabaseTableName,
			).
			WillReturnError(errors.New("db connection failed"))

		err = repo.Create(context.Background(), testServiceConfig)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db connection failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
