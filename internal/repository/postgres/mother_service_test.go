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

func TestMotherServiceRepository_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		expectedMotherService := &entity.MotherService{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name: "mother1",
			// TODO: add other fields based on migration file
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL ORDER BY "mother_services"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
				AddRow(expectedMotherService.ID, expectedMotherService.CreatedAt, expectedMotherService.UpdatedAt, nil, expectedMotherService.Name))

		result, err := repo.GetByID(context.Background(), uint64(1))
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMotherService.ID, result.ID)
		assert.Equal(t, expectedMotherService.Name, result.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("not found case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL ORDER BY "mother_services"."id" LIMIT $2`)).
			WithArgs(uint64(999), 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetByID(context.Background(), uint64(999))
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL ORDER BY "mother_services"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnError(errors.New("database connection failed"))

		result, err := repo.GetByID(context.Background(), uint64(1))
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMotherServiceRepository_GetAll(t *testing.T) {
	t.Run("success case with multiple records", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		expectedMotherServices := []*entity.MotherService{
			{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name: "mother1",
			},
			{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name: "mother2",
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
				AddRow(expectedMotherServices[0].ID, expectedMotherServices[0].CreatedAt, expectedMotherServices[0].UpdatedAt, nil, expectedMotherServices[0].Name).
				AddRow(expectedMotherServices[1].ID, expectedMotherServices[1].CreatedAt, expectedMotherServices[1].UpdatedAt, nil, expectedMotherServices[1].Name))

		result, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedMotherServices[0].ID, result[0].ID)
		assert.Equal(t, expectedMotherServices[0].Name, result[0].Name)
		assert.Equal(t, expectedMotherServices[1].ID, result[1].ID)
		assert.Equal(t, expectedMotherServices[1].Name, result[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case with empty result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}))

		result, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL`)).
			WillReturnError(errors.New("database connection failed"))

		result, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get mother service records")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
