package postgres

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseMetadataRepositoryInitialization(t *testing.T) {
	conn := new(mocks.Connection)
	db, _, err := conn.OpenConnection()
	require.NoError(t, err)

	var cfg *config.Config

	repo := NewDatabaseMetadataRepository(db, cfg)
	assert.NotNil(t, repo)

	tr, ok := repo.(*databaseMetadata)
	assert.True(t, ok)
	assert.NotNil(t, tr.db)
}

func TestDatabaseMetadataRepository_GetAll(t *testing.T) {
t.Run("success_case_with_some_result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		var cfg *config.Config

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT datname
		FROM pg_database
		WHERE datistemplate = false
		ORDER BY datname
	`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"datname"}).
					AddRow("postgres").
					AddRow("control_panel").
					AddRow("load_test_db"),
			)

		repo := NewDatabaseMetadataRepository(db, cfg)
		result, err := repo.GetAll(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
	})
t.Run("success_case_with_empty_result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		var cfg *config.Config

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT datname
		FROM pg_database
		WHERE datistemplate = false
		ORDER BY datname
	`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"datname"}),
			)

		repo := NewDatabaseMetadataRepository(db, cfg)
		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)

		require.NoError(t, mock.ExpectationsWereMet())
	})

t.Run("failure_case_database_error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		var cfg *config.Config

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT datname
		FROM pg_database
		WHERE datistemplate = false
		ORDER BY datname
	`)).
			WillReturnError(errors.New("database is down"))

		repo := NewDatabaseMetadataRepository(db, cfg)
		result, err := repo.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

// I comment this because db connection not mockable for now and
// if I want to run tests every time, the method create new connection to postgres.

// func TestDatabaseMetadataRepository_GetTablesByDBName(t *testing.T) {
// 	t.Run("success case - with result", func(t *testing.T) {
// 		conn := new(mocks.Connection)
// 		db, mock, err := conn.OpenConnection()
// 		require.NoError(t, err)

// 		var cfg *config.Config

// 		mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT table_name
// 		FROM information_schema.tables
// 		WHERE table_catalog = 'load_test_db' AND table_schema = 'public'
// 		ORDER BY table_name
// 	`)).
// 			WillReturnRows(
// 				sqlmock.NewRows([]string{"table_name"}).
// 					AddRow("events").
// 					AddRow("logs"),
// 			)

// 		repo := NewDatabaseMetadataRepository(db, cfg)
// 		result, err := repo.GetTablesByDBName(context.Background(), "load_test_db")

// 		assert.NoError(t, err)
// 		assert.Len(t, result, 2)
// 		assert.Equal(t, []string{"events", "logs"}, result)

// 		require.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("success case - with empty result", func(t *testing.T) {
// 		conn := new(mocks.Connection)
// 		db, mock, err := conn.OpenConnection()
// 		require.NoError(t, err)
// 		var cfg *config.Config

// 		mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT table_name
// 		FROM information_schema.tables
// 		WHERE table_catalog = 'load_test_db' AND table_schema = 'public'
// 		ORDER BY table_name
// 	`)).
// 			WillReturnRows(
// 				sqlmock.NewRows([]string{"table_name"}),
// 			)

// 		repo := NewDatabaseMetadataRepository(db, cfg)
// 		result, err := repo.GetTablesByDBName(context.Background(), "load_test_db")

// 		assert.NoError(t, err)
// 		assert.Len(t, result, 0)

// 		require.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("failure case - database error", func(t *testing.T) {
// 		conn := new(mocks.Connection)
// 		db, mock, err := conn.OpenConnection()
// 		require.NoError(t, err)
// 		var cfg *config.Config

// 		mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT table_name
// 		FROM information_schema.tables
// 		WHERE table_catalog = 'load_test_db' AND table_schema = 'public'
// 		ORDER BY table_name
// 	`)).
// 			WillReturnError(errors.New("query failed"))

// 		repo := NewDatabaseMetadataRepository(db, cfg)
// 		result, err := repo.GetTablesByDBName(context.Background(), "load_test_db")

// 		assert.Error(t, err)
// 		assert.Nil(t, result)

// 		require.NoError(t, mock.ExpectationsWereMet())
// 	})
// }

func TestIsTestScenarioTable(t *testing.T) {
	t.Run("success_case_exact_required_columns", func(t *testing.T) {
		columns := []string{
			"mother_service_id",
			"test_service_id",
			"start_tx_time",
			"step_num",
			"execution_id",
			"scenario_id",
			"duration_tx",
			"delay_before_tx",
			"http_status_code",
		}
		result := isTestScenarioTable(columns)
		if !result {
			t.Errorf("isTestScenarioTable() = false, expected true")
		}
	})

	t.Run("success_case_with_extra_columns", func(t *testing.T) {
		columns := []string{
			"mother_service_id",
			"test_service_id",
			"start_tx_time",
			"step_num",
			"execution_id",
			"scenario_id",
			"duration_tx",
			"delay_before_tx",
			"http_status_code",
			"extra_field_1",
			"extra_field_2",
		}
		result := isTestScenarioTable(columns)
		if !result {
			t.Errorf("isTestScenarioTable() = false, expected true")
		}
	})

	t.Run("failed_case_missing_one_required_column", func(t *testing.T) {
		columns := []string{
			"mother_service_id",
			"test_service_id",
			"start_tx_time",
			"step_num",
			"execution_id",
			"scenario_id",
			"duration_tx",
			"delay_before_tx",
		}
		result := isTestScenarioTable(columns)
		if result {
			t.Errorf("isTestScenarioTable() = true, expected false")
		}
	})

	t.Run("failed_case_empty_columns", func(t *testing.T) {
		columns := []string{}
		result := isTestScenarioTable(columns)
		if result {
			t.Errorf("isTestScenarioTable() = true, expected false")
		}
	})

	t.Run("failed_case_only_some_required_columns", func(t *testing.T) {
		columns := []string{
			"mother_service_id",
			"test_service_id",
		}
		result := isTestScenarioTable(columns)
		if result {
			t.Errorf("isTestScenarioTable() = true, expected false")
		}
	})
}

func TestIsMotherTable(t *testing.T) {
	t.Run("success_case_exact_required_columns", func(t *testing.T) {
		columns := []string{
			"event_id",
			"input",
			"output",
		}
		result := isMotherTable(columns)
		if !result {
			t.Errorf("isMotherTable() = false, expected true")
		}
	})

	t.Run("success_case_with_extra_columns", func(t *testing.T) {
		columns := []string{
			"event_id",
			"input",
			"output",
			"timestamp",
			"user_id",
		}
		result := isMotherTable(columns)
		if !result {
			t.Errorf("isMotherTable() = false, expected true")
		}
	})

	t.Run("failed_case_missing_one_required_column", func(t *testing.T) {
		columns := []string{
			"event_id",
			"input",
		}
		result := isMotherTable(columns)
		if result {
			t.Errorf("isMotherTable() = true, expected false")
		}
	})

	t.Run("failed_case_missing_two_required_columns", func(t *testing.T) {
		columns := []string{
			"event_id",
		}
		result := isMotherTable(columns)
		if result {
			t.Errorf("isMotherTable() = true, expected false")
		}
	})

	t.Run("failed_case_empty_columns", func(t *testing.T) {
		columns := []string{}
		result := isMotherTable(columns)
		if result {
			t.Errorf("isMotherTable() = true, expected false")
		}
	})

	t.Run("failed_case_completely_different_columns", func(t *testing.T) {
		columns := []string{
			"random_field_1",
			"random_field_2",
		}
		result := isMotherTable(columns)
		if result {
			t.Errorf("isMotherTable() = true, expected false")
		}
	})
}