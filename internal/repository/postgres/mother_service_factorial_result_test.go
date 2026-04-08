package postgres

import (
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"fmt"

	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFactorialRepository_Create_Success(t *testing.T) {
	table := "factorials"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	assert.Nil(t, err)

	repo := NewMotherServiceFactorialResultRepository(&config.Config{
		Postgres: config.Postgres{},
	})
	now := time.Now()

	factorial := &entity.Factorial{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		Input:           "5",
		Output:          "120",
		MotherServiceId: 3,
		MotherService: &entity.MotherService{
			DatabaseName:      "mother",
			DatabaseTableName: table,
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`INSERT INTO "%s" ("created_at","updated_at","deleted_at","input","output","mother_service_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`, table))).
		WithArgs(now, now, nil, factorial.Input, factorial.Output, factorial.MotherServiceId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(t.Context(), factorial, func(cfg any) (database.Database, error) {
		return db, nil
	})

	require.NoError(t, err)
	require.Equal(t, uint(1), factorial.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFactorialRepository_Create_DBError(t *testing.T) {
	table := "factorials"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	require.NoError(t, err)

	repo := NewMotherServiceFactorialResultRepository(&config.Config{
		Postgres: config.Postgres{},
	})
	now := time.Now()

	factorial := &entity.Factorial{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		Input:           "5",
		Output:          "120",
		MotherServiceId: 3,
		MotherService: &entity.MotherService{
			DatabaseName:      "mother",
			DatabaseTableName: table,
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`INSERT INTO "%s" ("created_at","updated_at","deleted_at","input","output","mother_service_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`, table))).
		WithArgs(now, now, nil, factorial.Input, factorial.Output, factorial.MotherServiceId).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	err = repo.Create(t.Context(), factorial, func(cfg any) (database.Database, error) {
		return db, nil
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create factorial record")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFactorialRepository_Create_DBInitializerError(t *testing.T) {
	table := "factorials"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	require.NoError(t, err)

	repo := NewMotherServiceFactorialResultRepository(&config.Config{
		Postgres: config.Postgres{},
	})
	now := time.Now()

	factorial := &entity.Factorial{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		Input:           "5",
		Output:          "120",
		MotherServiceId: 3,
		MotherService: &entity.MotherService{
			DatabaseName:      "mother",
			DatabaseTableName: table,
		},
	}

	err = repo.Create(t.Context(), factorial, func(cfg any) (database.Database, error) {
		return db, errors.New("something went wrong")
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "something went wrong")
	require.NoError(t, mock.ExpectationsWereMet())
}
