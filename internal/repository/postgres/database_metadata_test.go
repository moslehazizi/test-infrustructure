package postgres

import (
	"context"
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

	repo := NewDatabaseMetadataRepository(db)
	assert.NotNil(t, repo)

	tr, ok := repo.(*databaseMetadata)
	assert.True(t, ok)
	assert.NotNil(t, tr.db)
}

func TestDatabaseMetadataRepository_GetAll(t *testing.T) {
	t.Run("success case - with some result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

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

		repo := NewDatabaseMetadataRepository(db)
		result, err := repo.GetAll(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
	})
	t.Run("success case - with empty result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT datname
		FROM pg_database
		WHERE datistemplate = false
		ORDER BY datname
	`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"datname"}),
			)

		repo := NewDatabaseMetadataRepository(db)
		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure case - database error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT datname
		FROM pg_database
		WHERE datistemplate = false
		ORDER BY datname
	`)).
			WillReturnError(errors.New("database is down"))

		repo := NewDatabaseMetadataRepository(db)
		result, err := repo.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabaseMetadataRepository_GetTablesByDBName(t *testing.T) {
	t.Run("success case - with result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"table_name"}).
					AddRow("events").
					AddRow("logs"),
			)

		repo := NewDatabaseMetadataRepository(db)
		result, err := repo.GetTablesByDBName(context.Background(), "load_test_db")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, []string{"events", "logs"}, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case - with empty result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"table_name"}),
			)

		repo := NewDatabaseMetadataRepository(db)
		result, err := repo.GetTablesByDBName(context.Background(), "load_test_db")

		assert.NoError(t, err)
		assert.Len(t, result, 0)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure case - database error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`)).
			WillReturnError(errors.New("query failed"))

		repo := NewDatabaseMetadataRepository(db)
		result, err := repo.GetTablesByDBName(context.Background(), "load_test_db")

		assert.Error(t, err)
		assert.Nil(t, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
