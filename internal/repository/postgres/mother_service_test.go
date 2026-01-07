package postgres

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/postgres/mocks"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMotherServiceRepository_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()

		repo := NewMotherServiceRepository(db)
		now := time.Now()

		motherService := &entity.MotherService{
			Model: gorm.Model{CreatedAt: now, UpdatedAt: now},
			Name:  "mother1",
			// TODO: add other fields based on migration file
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "mother_services" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
			WithArgs(now, now, nil, motherService.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err = repo.Create(t.Context(), motherService)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), motherService.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewMotherServiceRepository(db)
		now := time.Now()

		motherService := &entity.MotherService{
			Model: gorm.Model{CreatedAt: now, UpdatedAt: now},
			Name:  "mother1",
			// TODO: add other fields based on migration file
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "mother_services" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
			WithArgs(now, now, nil, motherService.Name).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		err = repo.Create(t.Context(), motherService)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create mother service record")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
