package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

func TestTestServiceConfigRepo_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestServiceConfigRepository(db)

		now := time.Now()
		sampleNumber := 1

		expectedTestServiceConfig := &entity.TestServiceConfig{
			ID:                    uint64(1),
			CreatedAt:             now,
			UpdatedAt:             now,
			TestScenarioID:        2,
			MaxRequests:           100,
			MaxDuration:           1200,
			RequestDelayDuration:  nil,
			RandomRequestDelayMin: nil,
			RandomRequestDelayMax: nil,
			FixedTestNumber:       &sampleNumber,
			RandomTestNumberMin:   nil,
			RandomTestNumberMax:   nil,
			BadValueRate:          0,
			NegativeValueRate:     0,
			RealValueRate:         0,
			ZeroValueRate:         0,
			StringValueRate:       0,
			LongStringValueRate:   0,
			NullValueRate:         0,
			DatabaseName:          "db-name",
			DatabaseTableName:     "db-tb-name",
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_service_configs" WHERE "test_service_configs"."id" = $1 ORDER BY "test_service_configs"."id" LIMIT $2`)).
			WithArgs(2, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "test_scenario_id", "max_requests",
				"max_duration", "request_delay_duration",
				"random_request_delay_min", "random_request_delay_max", "fixed_test_number", "random_test_number_min", "random_test_number_max", "bad_value_rate", "negative_value_rate", "real_value_rate", "zero_value_rate", "string_value_rate", "long_string_value_rate", "null_value_rate", "database_name", "database_table_name",
			}).
				AddRow(
					expectedTestServiceConfig.ID,
					expectedTestServiceConfig.CreatedAt,
					expectedTestServiceConfig.UpdatedAt,
					expectedTestServiceConfig.TestScenarioID,
					expectedTestServiceConfig.MaxRequests,
					expectedTestServiceConfig.MaxDuration,
					expectedTestServiceConfig.RequestDelayDuration,
					expectedTestServiceConfig.RandomRequestDelayMin,
					expectedTestServiceConfig.RandomRequestDelayMax,
					expectedTestServiceConfig.FixedTestNumber,
					expectedTestServiceConfig.RandomTestNumberMin,
					expectedTestServiceConfig.RandomTestNumberMax,
					expectedTestServiceConfig.BadValueRate,
					expectedTestServiceConfig.NegativeValueRate,
					expectedTestServiceConfig.RealValueRate,
					expectedTestServiceConfig.ZeroValueRate,
					expectedTestServiceConfig.StringValueRate,
					expectedTestServiceConfig.LongStringValueRate,
					expectedTestServiceConfig.NullValueRate,
					expectedTestServiceConfig.DatabaseName,
					expectedTestServiceConfig.DatabaseTableName,
				))

		result, err := repo.GetByID(context.Background(), uint64(2))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestServiceConfig.ID, result.ID)
		assert.Equal(t, expectedTestServiceConfig.TestScenarioID, result.TestScenarioID)
		assert.Equal(t, expectedTestServiceConfig.MaxRequests, result.MaxRequests)
		assert.Equal(t, expectedTestServiceConfig.MaxDuration, result.MaxDuration)
		assert.Equal(t, expectedTestServiceConfig.DatabaseTableName, result.DatabaseTableName)
		assert.Equal(t, expectedTestServiceConfig.BadValueRate, result.BadValueRate)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - not found", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestServiceConfigRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_service_configs" WHERE "test_service_configs"."id" = $1 ORDER BY "test_service_configs"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ErrTestServiceConfigNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - database error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestServiceConfigRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_service_configs" WHERE "test_service_configs"."id" = $1 ORDER BY "test_service_configs"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnError(errors.New("error happened"))

		result, err := repo.GetByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "error happened")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTestServiceConfigRepository_UpdateByScenarioID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestServiceConfigRepository(db)

		scenarioID := uint64(1)
		sampleInt := 5

		cfg := &entity.TestServiceConfig{
			MaxRequests:           100,
			MaxDuration:           60,
			RequestDelayDuration:  &sampleInt,
			RandomRequestDelayMin: nil,
			RandomRequestDelayMax: nil,
			FixedTestNumber:       &sampleInt,
			RandomTestNumberMin:   nil,
			RandomTestNumberMax:   nil,
			BadValueRate:          1,
			NegativeValueRate:     2,
			RealValueRate:         3,
			ZeroValueRate:         4,
			StringValueRate:       5,
			LongStringValueRate:   6,
			NullValueRate:         7,
			DatabaseName:          "db1",
			DatabaseTableName:     "table1",
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_service_configs" SET`,
		)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err = repo.UpdateByScenarioID(context.Background(), scenarioID, cfg)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update fails", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestServiceConfigRepository(db)

		scenarioID := uint64(1)

		cfg := &entity.TestServiceConfig{
			MaxRequests: 100,
			MaxDuration: 60,
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_service_configs" SET`,
		)).
			WillReturnError(errors.New("update failed"))

		mock.ExpectRollback()

		err = repo.UpdateByScenarioID(context.Background(), scenarioID, cfg)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
