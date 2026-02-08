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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMotherServiceRepository_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewMotherServiceRepository(db)
		now := time.Now()

		responseDelayDuration := 100
		randomDelayMin := 50
		randomDelayMax := 150
		serviceAddress := "http://service.example.com"

		motherService := &entity.MotherService{
			CreatedAt:                now,
			UpdatedAt:                now,
			Name:                     "mother1",
			ExceptionRate:            10,
			ResponseDelayRate:        20,
			ResponseDelayDuration:    &responseDelayDuration,
			RandomResponseDelayMin:   &randomDelayMin,
			RandomResponseDelayMax:   &randomDelayMax,
			Status:                   entity.MotherServiceStatusPending,
			ServiceDeploymentAddress: &serviceAddress,
			DatabaseName:             "test_db",
			DatabaseTableName:        "test_table",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "mother_services" ("created_at","updated_at","deleted_at","name","exception_rate","response_delay_rate","response_delay_duration","random_response_delay_min","random_response_delay_max","status","service_deployment_address","database_name","database_table_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`)).
			WithArgs(
				now, now, nil,
				motherService.Name,
				motherService.ExceptionRate,
				motherService.ResponseDelayRate,
				motherService.ResponseDelayDuration,
				motherService.RandomResponseDelayMin,
				motherService.RandomResponseDelayMax,
				motherService.Status,
				motherService.ServiceDeploymentAddress,
				motherService.DatabaseName,
				motherService.DatabaseTableName,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		id, err := repo.Create(context.Background(), motherService)
		assert.NoError(t, err)
		assert.Equal(t, id, uint64(1))
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
			CreatedAt:         now,
			UpdatedAt:         now,
			Name:              "mother1",
			ExceptionRate:     0.0,
			ResponseDelayRate: 0.0,
			Status:            entity.MotherServiceStatusPending,
			DatabaseName:      "test_db",
			DatabaseTableName: "test_table",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "mother_services" ("created_at","updated_at","deleted_at","name","exception_rate","response_delay_rate","response_delay_duration","random_response_delay_min","random_response_delay_max","status","service_deployment_address","database_name","database_table_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`)).
			WithArgs(
				now, now, nil,
				motherService.Name,
				motherService.ExceptionRate,
				motherService.ResponseDelayRate,
				motherService.ResponseDelayDuration,
				motherService.RandomResponseDelayMin,
				motherService.RandomResponseDelayMax,
				motherService.Status,
				motherService.ServiceDeploymentAddress,
				motherService.DatabaseName,
				motherService.DatabaseTableName,
			).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		id, err := repo.Create(context.Background(), motherService)
		assert.Equal(t, id, uint64(0))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create mother service record")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate name error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewMotherServiceRepository(db)
		now := time.Now()

		motherService := &entity.MotherService{
			CreatedAt:         now,
			UpdatedAt:         now,
			Name:              "mother1",
			ExceptionRate:     0.0,
			ResponseDelayRate: 0.0,
			Status:            entity.MotherServiceStatusPending,
			DatabaseName:      "test_db",
			DatabaseTableName: "test_table",
		}

		duplicateError := &pgconn.PgError{
			Code:           "23505",
			Message:        "duplicate key value violates unique constraint",
			ConstraintName: "mother_services_name_key",
			Detail:         `Key (name)=(mother1) already exists.`,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "mother_services" ("created_at","updated_at","deleted_at","name","exception_rate","response_delay_rate","response_delay_duration","random_response_delay_min","random_response_delay_max","status","service_deployment_address","database_name","database_table_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`)).
			WithArgs(
				now, now, nil,
				motherService.Name,
				motherService.ExceptionRate,
				motherService.ResponseDelayRate,
				motherService.ResponseDelayDuration,
				motherService.RandomResponseDelayMin,
				motherService.RandomResponseDelayMax,
				motherService.Status,
				motherService.ServiceDeploymentAddress,
				motherService.DatabaseName,
				motherService.DatabaseTableName,
			).
			WillReturnError(duplicateError)
		mock.ExpectRollback()

		id, err := repo.Create(context.Background(), motherService)
		assert.Equal(t, id, uint64(0))
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceAlreadyExist)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			assert.Equal(t, "23505", pgErr.Code)
			assert.Contains(t, pgErr.Detail, "mother1")
		}
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
		responseDelayDuration := 100
		randomDelayMin := 50
		randomDelayMax := 150
		serviceAddress := "http://service.example.com"

		expectedMotherService := &entity.MotherService{
			ID:                       uint64(1),
			CreatedAt:                now,
			UpdatedAt:                now,
			Name:                     "mother1",
			ExceptionRate:            10,
			ResponseDelayRate:        20,
			ResponseDelayDuration:    &responseDelayDuration,
			RandomResponseDelayMin:   &randomDelayMin,
			RandomResponseDelayMax:   &randomDelayMax,
			Status:                   entity.MotherServiceStatusRunning,
			ServiceDeploymentAddress: &serviceAddress,
			DatabaseName:             "test_db",
			DatabaseTableName:        "test_table",
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL ORDER BY "mother_services"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"status", "service_deployment_address",
				"database_name", "database_table_name",
			}).
				AddRow(
					expectedMotherService.ID,
					expectedMotherService.CreatedAt,
					expectedMotherService.UpdatedAt,
					nil,
					expectedMotherService.Name,
					expectedMotherService.ExceptionRate,
					expectedMotherService.ResponseDelayRate,
					expectedMotherService.ResponseDelayDuration,
					expectedMotherService.RandomResponseDelayMin,
					expectedMotherService.RandomResponseDelayMax,
					expectedMotherService.Status,
					expectedMotherService.ServiceDeploymentAddress,
					expectedMotherService.DatabaseName,
					expectedMotherService.DatabaseTableName,
				))

		result, err := repo.GetByID(context.Background(), uint64(1))
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMotherService.ID, result.ID)
		assert.Equal(t, expectedMotherService.Name, result.Name)
		assert.Equal(t, expectedMotherService.ExceptionRate, result.ExceptionRate)
		assert.Equal(t, expectedMotherService.ResponseDelayRate, result.ResponseDelayRate)
		assert.Equal(t, expectedMotherService.Status, result.Status)
		assert.Equal(t, expectedMotherService.DatabaseName, result.DatabaseName)
		assert.Equal(t, expectedMotherService.DatabaseTableName, result.DatabaseTableName)
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
		assert.Equal(t, pkg.ErrMotherServiceNotFound, err)
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

func TestMotherServiceRepository_GetPaginated(t *testing.T) {
	t.Run("success case with pagination - page 1", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		serviceAddress1 := "http://service1.example.com"
		serviceAddress2 := "http://service2.example.com"

		expectedMotherServices := []*entity.MotherService{
			{
				ID:                       uint64(5),
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "mother5",
				ExceptionRate:            0.0,
				ResponseDelayRate:        0.0,
				Status:                   entity.MotherServiceStatusPending,
				ServiceDeploymentAddress: &serviceAddress2,
				DatabaseName:             "test_db5",
				DatabaseTableName:        "test_table5",
			},
			{
				ID:                       uint64(4),
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "mother4",
				ExceptionRate:            10,
				ResponseDelayRate:        20,
				Status:                   entity.MotherServiceStatusPending,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db4",
				DatabaseTableName:        "test_table4",
			},
		}

		paginationRequest := entity.PaginationRequest{
			Page:    1,
			PerPage: 2,
		} // LIMIT 2 (no OFFSET because offset=0)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"status", "service_deployment_address",
				"database_name", "database_table_name",
			}).
				AddRow(
					expectedMotherServices[0].ID,
					expectedMotherServices[0].CreatedAt,
					expectedMotherServices[0].UpdatedAt,
					nil,
					expectedMotherServices[0].Name,
					expectedMotherServices[0].ExceptionRate,
					expectedMotherServices[0].ResponseDelayRate,
					expectedMotherServices[0].ResponseDelayDuration,
					expectedMotherServices[0].RandomResponseDelayMin,
					expectedMotherServices[0].RandomResponseDelayMax,
					expectedMotherServices[0].Status,
					expectedMotherServices[0].ServiceDeploymentAddress,
					expectedMotherServices[0].DatabaseName,
					expectedMotherServices[0].DatabaseTableName,
				).
				AddRow(
					expectedMotherServices[1].ID,
					expectedMotherServices[1].CreatedAt,
					expectedMotherServices[1].UpdatedAt,
					nil,
					expectedMotherServices[1].Name,
					expectedMotherServices[1].ExceptionRate,
					expectedMotherServices[1].ResponseDelayRate,
					expectedMotherServices[1].ResponseDelayDuration,
					expectedMotherServices[1].RandomResponseDelayMin,
					expectedMotherServices[1].RandomResponseDelayMax,
					expectedMotherServices[1].Status,
					expectedMotherServices[1].ServiceDeploymentAddress,
					expectedMotherServices[1].DatabaseName,
					expectedMotherServices[1].DatabaseTableName,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		assert.Equal(t, uint64(5), result[0].ID)
		assert.Equal(t, "mother5", result[0].Name)
		assert.Equal(t, uint64(4), result[1].ID)
		assert.Equal(t, "mother4", result[1].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case with pagination - page 2", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		serviceAddress1 := "http://service1.example.com"
		serviceAddress2 := "http://service2.example.com"

		expectedMotherServices := []*entity.MotherService{
			{
				ID:                       uint64(3),
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "mother3",
				ExceptionRate:            0.0,
				ResponseDelayRate:        0.0,
				Status:                   entity.MotherServiceStatusPending,
				ServiceDeploymentAddress: &serviceAddress2,
				DatabaseName:             "test_db3",
				DatabaseTableName:        "test_table3",
			},
			{
				ID:                       uint64(2),
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "mother2",
				ExceptionRate:            10,
				ResponseDelayRate:        20,
				Status:                   entity.MotherServiceStatusPending,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db2",
				DatabaseTableName:        "test_table2",
			},
		}

		paginationRequest := entity.PaginationRequest{
			Page:    2,
			PerPage: 2,
		} // LIMIT 2 OFFSET 2

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(2, 2).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"status", "service_deployment_address",
				"database_name", "database_table_name",
			}).
				AddRow(
					expectedMotherServices[0].ID,
					expectedMotherServices[0].CreatedAt,
					expectedMotherServices[0].UpdatedAt,
					nil,
					expectedMotherServices[0].Name,
					expectedMotherServices[0].ExceptionRate,
					expectedMotherServices[0].ResponseDelayRate,
					expectedMotherServices[0].ResponseDelayDuration,
					expectedMotherServices[0].RandomResponseDelayMin,
					expectedMotherServices[0].RandomResponseDelayMax,
					expectedMotherServices[0].Status,
					expectedMotherServices[0].ServiceDeploymentAddress,
					expectedMotherServices[0].DatabaseName,
					expectedMotherServices[0].DatabaseTableName,
				).
				AddRow(
					expectedMotherServices[1].ID,
					expectedMotherServices[1].CreatedAt,
					expectedMotherServices[1].UpdatedAt,
					nil,
					expectedMotherServices[1].Name,
					expectedMotherServices[1].ExceptionRate,
					expectedMotherServices[1].ResponseDelayRate,
					expectedMotherServices[1].ResponseDelayDuration,
					expectedMotherServices[1].RandomResponseDelayMin,
					expectedMotherServices[1].RandomResponseDelayMax,
					expectedMotherServices[1].Status,
					expectedMotherServices[1].ServiceDeploymentAddress,
					expectedMotherServices[1].DatabaseName,
					expectedMotherServices[1].DatabaseTableName,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		assert.Equal(t, uint64(3), result[0].ID)
		assert.Equal(t, "mother3", result[0].Name)
		assert.Equal(t, uint64(2), result[1].ID)
		assert.Equal(t, "mother2", result[1].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case with pagination - page 3 with 10 per page", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		serviceAddress := "http://service.example.com"

		expectedMotherService := &entity.MotherService{
			ID:                       uint64(25),
			CreatedAt:                now,
			UpdatedAt:                now,
			Name:                     "mother25",
			ExceptionRate:            0.0,
			ResponseDelayRate:        0.0,
			Status:                   entity.MotherServiceStatusPending,
			ServiceDeploymentAddress: &serviceAddress,
			DatabaseName:             "test_db25",
			DatabaseTableName:        "test_table25",
		}

		paginationRequest := entity.PaginationRequest{
			Page:    3,
			PerPage: 10,
		} // Page 3, perPage 10: LIMIT 10 OFFSET 20

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(10, 20).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"status", "service_deployment_address",
				"database_name", "database_table_name",
			}).
				AddRow(
					expectedMotherService.ID,
					expectedMotherService.CreatedAt,
					expectedMotherService.UpdatedAt,
					nil,
					expectedMotherService.Name,
					expectedMotherService.ExceptionRate,
					expectedMotherService.ResponseDelayRate,
					expectedMotherService.ResponseDelayDuration,
					expectedMotherService.RandomResponseDelayMin,
					expectedMotherService.RandomResponseDelayMax,
					expectedMotherService.Status,
					expectedMotherService.ServiceDeploymentAddress,
					expectedMotherService.DatabaseName,
					expectedMotherService.DatabaseTableName,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, uint64(25), result[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case without pagination - return all items (page=0, perPage=0)", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		serviceAddress1 := "http://service1.example.com"
		serviceAddress2 := "http://service2.example.com"
		serviceAddress3 := "http://service3.example.com"

		expectedMotherServices := []*entity.MotherService{
			{
				ID:                       uint64(3),
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "mother3",
				ExceptionRate:            0.0,
				ResponseDelayRate:        0.0,
				Status:                   entity.MotherServiceStatusPending,
				ServiceDeploymentAddress: &serviceAddress3,
				DatabaseName:             "test_db3",
				DatabaseTableName:        "test_table3",
			},
			{
				ID:                       uint64(2),
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "mother2",
				ExceptionRate:            0.0,
				ResponseDelayRate:        0.0,
				Status:                   entity.MotherServiceStatusPending,
				ServiceDeploymentAddress: &serviceAddress2,
				DatabaseName:             "test_db2",
				DatabaseTableName:        "test_table2",
			},
			{
				ID:                       uint64(1),
				CreatedAt:                now,
				UpdatedAt:                now,
				Name:                     "mother1",
				ExceptionRate:            10,
				ResponseDelayRate:        20,
				Status:                   entity.MotherServiceStatusPending,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db1",
				DatabaseTableName:        "test_table1",
			},
		}

		paginationRequest := entity.PaginationRequest{
			Page:    0,
			PerPage: 0,
		} // No LIMIT or OFFSET

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"status", "service_deployment_address",
				"database_name", "database_table_name",
			}).
				AddRow(
					expectedMotherServices[0].ID,
					expectedMotherServices[0].CreatedAt,
					expectedMotherServices[0].UpdatedAt,
					nil,
					expectedMotherServices[0].Name,
					expectedMotherServices[0].ExceptionRate,
					expectedMotherServices[0].ResponseDelayRate,
					expectedMotherServices[0].ResponseDelayDuration,
					expectedMotherServices[0].RandomResponseDelayMin,
					expectedMotherServices[0].RandomResponseDelayMax,
					expectedMotherServices[0].Status,
					expectedMotherServices[0].ServiceDeploymentAddress,
					expectedMotherServices[0].DatabaseName,
					expectedMotherServices[0].DatabaseTableName,
				).
				AddRow(
					expectedMotherServices[1].ID,
					expectedMotherServices[1].CreatedAt,
					expectedMotherServices[1].UpdatedAt,
					nil,
					expectedMotherServices[1].Name,
					expectedMotherServices[1].ExceptionRate,
					expectedMotherServices[1].ResponseDelayRate,
					expectedMotherServices[1].ResponseDelayDuration,
					expectedMotherServices[1].RandomResponseDelayMin,
					expectedMotherServices[1].RandomResponseDelayMax,
					expectedMotherServices[1].Status,
					expectedMotherServices[1].ServiceDeploymentAddress,
					expectedMotherServices[1].DatabaseName,
					expectedMotherServices[1].DatabaseTableName,
				).
				AddRow(
					expectedMotherServices[2].ID,
					expectedMotherServices[2].CreatedAt,
					expectedMotherServices[2].UpdatedAt,
					nil,
					expectedMotherServices[2].Name,
					expectedMotherServices[2].ExceptionRate,
					expectedMotherServices[2].ResponseDelayRate,
					expectedMotherServices[2].ResponseDelayDuration,
					expectedMotherServices[2].RandomResponseDelayMin,
					expectedMotherServices[2].RandomResponseDelayMax,
					expectedMotherServices[2].Status,
					expectedMotherServices[2].ServiceDeploymentAddress,
					expectedMotherServices[2].DatabaseName,
					expectedMotherServices[2].DatabaseTableName,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, uint64(3), result[0].ID)
		assert.Equal(t, uint64(2), result[1].ID)
		assert.Equal(t, uint64(1), result[2].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case with empty result - page beyond available data", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		paginationRequest := entity.PaginationRequest{
			Page:    100,
			PerPage: 10,
		} // LIMIT 10 OFFSET 990

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(10, 990).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"status", "service_deployment_address",
				"database_name", "database_table_name",
			}))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case with pagination", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		paginationRequest := entity.PaginationRequest{
			Page:    2,
			PerPage: 10,
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(10, 10).
			WillReturnError(errors.New("database connection failed"))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get mother service records")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case if page is negative", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		paginationRequest := entity.PaginationRequest{
			Page:    -2,
			PerPage: 10,
		}

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrNegativePageOrPerPageNotAllowed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case if perPage is negative", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		paginationRequest := entity.PaginationRequest{
			Page:    2,
			PerPage: -10,
		}

		result, err := repo.GetPaginated(context.Background(), paginationRequest)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrNegativePageOrPerPageNotAllowed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
