package postgres

import (
	"control-panel-service/internal/domain/entity"
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

	repo := NewFactorialRepository(db)
	now := time.Now()

	factorial := &entity.Factorial{
		Model:  gorm.Model{CreatedAt: now, UpdatedAt: now},
		Input:  "5",
		Output: "120",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`INSERT INTO "%s" ("created_at","updated_at","deleted_at","input","output") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`, table))).
		WithArgs(now, now, nil, factorial.Input, factorial.Output).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(t.Context(), factorial)

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

	repo := NewFactorialRepository(db)
	now := time.Now()

	factorial := &entity.Factorial{
		Model:  gorm.Model{CreatedAt: now, UpdatedAt: now},
		Input:  "5",
		Output: "120",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`INSERT INTO "%s" ("created_at","updated_at","deleted_at","input","output") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`, table))).
		WithArgs(now, now, nil, factorial.Input, factorial.Output).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	err = repo.Create(t.Context(), factorial)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create factorial record")
	require.NoError(t, mock.ExpectationsWereMet())
}
