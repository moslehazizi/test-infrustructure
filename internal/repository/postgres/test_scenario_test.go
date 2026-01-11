package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/postgres/mocks"
	"control-panel-service/pkg"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTestScenarioRepository_New(t *testing.T) {
	mockConn := new(mocks.Connection)
	db, _, err := mockConn.OpenConnection()
	require.NoError(t, err)

	repo := NewTestScenarioRepository(db)
	assert.NotNil(t, repo)

	ts, ok := repo.(*testScenario)
	assert.True(t, ok)
	assert.NotNil(t, ts.db)
}

func TestTestScenarioRepository_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)
		now := time.Now()

		testScenario := &entity.TestScenario{
			CreatedAt:           now,
			UpdatedAt:           now,
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			Status:              entity.ScenarioStatus(entity.ScenarioStatusPending),
			MaxTestServiceCount: nil,
			ExecutionDuration:   nil,
			AutoStepChangeRate:  nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","max_test_service_count","execution_duration","auto_step_change_rate") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
			WithArgs(
				testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				nil, nil, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		id, err := repo.Create(context.Background(), testScenario)

		assert.NoError(t, err)
		assert.Equal(t, id, uint64(1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)
		now := time.Now()

		testScenario := &entity.TestScenario{
			CreatedAt:           now,
			UpdatedAt:           now,
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			Status:              entity.ScenarioStatus(entity.ScenarioStatusPending),
			MaxTestServiceCount: nil,
			ExecutionDuration:   nil,
			AutoStepChangeRate:  nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","max_test_service_count","execution_duration","auto_step_change_rate") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
			WithArgs(testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				nil, nil, nil).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		id, err := repo.Create(context.Background(), testScenario)

		assert.Error(t, err)
		assert.Equal(t, id, uint64(0))
		assert.Contains(t, err.Error(), "failed to create test scenario record")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenario := &entity.TestScenario{
			ID:                  uint64(1),
			CreatedAt:           now,
			UpdatedAt:           now,
			Name:                "load1",
			TestCategoryID:      uint64(3),
			MotherServiceID:     uint64(2),
			Status:              entity.ScenarioStatusPending,
			MaxTestServiceCount: nil,
			ExecutionDuration:   nil,
			AutoStepChangeRate:  nil,
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."id" = $1 AND "test_scenarios"."deleted_at" IS NULL ORDER BY "test_scenarios"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}).
				AddRow(
					expectedTestScenario.ID,
					expectedTestScenario.CreatedAt,
					expectedTestScenario.UpdatedAt,
					nil,
					expectedTestScenario.Name,
					expectedTestScenario.TestCategoryID,
					expectedTestScenario.MotherServiceID,
					expectedTestScenario.Status,
					expectedTestScenario.MaxTestServiceCount,
					expectedTestScenario.ExecutionDuration,
					expectedTestScenario.AutoStepChangeRate,
				))

		result, err := repo.GetByID(context.Background(), uint64(1))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario.ID, result.ID)
		assert.Equal(t, expectedTestScenario.Name, result.Name)
		assert.Equal(t, expectedTestScenario.TestCategoryID, result.TestCategoryID)
		assert.Equal(t, expectedTestScenario.MotherServiceID, result.MotherServiceID)
		assert.Equal(t, expectedTestScenario.Status, result.Status)
		assert.Equal(t, expectedTestScenario.MaxTestServiceCount, result.MaxTestServiceCount)
		assert.Equal(t, expectedTestScenario.ExecutionDuration, result.ExecutionDuration)
		assert.Equal(t, expectedTestScenario.AutoStepChangeRate, result.AutoStepChangeRate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - record not found", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."id" = $1 AND "test_scenarios"."deleted_at" IS NULL ORDER BY "test_scenarios"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetByID(context.Background(), uint64(1))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ErrTestScenarioNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - database error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."id" = $1 AND "test_scenarios"."deleted_at" IS NULL ORDER BY "test_scenarios"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnError(errors.New("failed to get record"))

		result, err := repo.GetByID(context.Background(), uint64(1))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetPaginated(t *testing.T) {
	t.Run("success case - page 1 per page 2", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenarios := []*entity.TestScenario{
			{
				ID:                  uint64(5),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load1",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "smoke1",
				TestCategoryID:      uint64(5),
				MotherServiceID:     uint64(4),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
		}

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    1,
			PerPage: 2,
		} // LIMIT 2 (no OFFSET because offset=0)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}).
				AddRow(
					expectedTestScenarios[0].ID,
					expectedTestScenarios[0].CreatedAt,
					expectedTestScenarios[0].UpdatedAt,
					nil,
					expectedTestScenarios[0].Name,
					expectedTestScenarios[0].TestCategoryID,
					expectedTestScenarios[0].MotherServiceID,
					expectedTestScenarios[0].Status,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].ExecutionDuration,
					expectedTestScenarios[0].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[1].ID,
					expectedTestScenarios[1].CreatedAt,
					expectedTestScenarios[1].UpdatedAt,
					nil,
					expectedTestScenarios[1].Name,
					expectedTestScenarios[1].TestCategoryID,
					expectedTestScenarios[1].MotherServiceID,
					expectedTestScenarios[1].Status,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].ExecutionDuration,
					expectedTestScenarios[1].AutoStepChangeRate,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, uint64(5), result[0].ID)
		assert.Equal(t, expectedTestScenarios[0].Name, result[0].Name)
		assert.Equal(t, uint64(4), result[1].ID)
		assert.Equal(t, expectedTestScenarios[1].Name, result[1].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case - page 2 per page 2", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenarios := []*entity.TestScenario{
			{
				ID:                  uint64(5),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load1",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "smoke1",
				TestCategoryID:      uint64(5),
				MotherServiceID:     uint64(4),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
		}

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		} // LIMIT 2 OFFSET 2

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(2, 2).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}).
				AddRow(
					expectedTestScenarios[0].ID,
					expectedTestScenarios[0].CreatedAt,
					expectedTestScenarios[0].UpdatedAt,
					nil,
					expectedTestScenarios[0].Name,
					expectedTestScenarios[0].TestCategoryID,
					expectedTestScenarios[0].MotherServiceID,
					expectedTestScenarios[0].Status,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].ExecutionDuration,
					expectedTestScenarios[0].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[1].ID,
					expectedTestScenarios[1].CreatedAt,
					expectedTestScenarios[1].UpdatedAt,
					nil,
					expectedTestScenarios[1].Name,
					expectedTestScenarios[1].TestCategoryID,
					expectedTestScenarios[1].MotherServiceID,
					expectedTestScenarios[1].Status,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].ExecutionDuration,
					expectedTestScenarios[1].AutoStepChangeRate,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		assert.Equal(t, expectedTestScenarios[0].ID, result[0].ID)
		assert.Equal(t, expectedTestScenarios[0].Name, result[0].Name)
		assert.Equal(t, expectedTestScenarios[1].ID, result[1].ID)
		assert.Equal(t, expectedTestScenarios[1].Name, result[1].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case - page 2 per page 1", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenarios := []*entity.TestScenario{
			{
				ID:                  uint64(5),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load1",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
		}

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 1,
		} // LIMIT 1 OFFSET 1

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}).
				AddRow(
					expectedTestScenarios[0].ID,
					expectedTestScenarios[0].CreatedAt,
					expectedTestScenarios[0].UpdatedAt,
					nil,
					expectedTestScenarios[0].Name,
					expectedTestScenarios[0].TestCategoryID,
					expectedTestScenarios[0].MotherServiceID,
					expectedTestScenarios[0].Status,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].ExecutionDuration,
					expectedTestScenarios[0].AutoStepChangeRate,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)

		assert.Equal(t, expectedTestScenarios[0].ID, result[0].ID)
		assert.Equal(t, expectedTestScenarios[0].Name, result[0].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case - page 3 per page 5", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenarios := []*entity.TestScenario{
			{
				ID:                  uint64(6),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load1",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(5),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load2",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load3",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(3),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load4",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(2),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load5",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
		}

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    3,
			PerPage: 5,
		} // LIMIT 5 OFFSET 10

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(5, 10).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}).
				AddRow(
					expectedTestScenarios[0].ID,
					expectedTestScenarios[0].CreatedAt,
					expectedTestScenarios[0].UpdatedAt,
					nil,
					expectedTestScenarios[0].Name,
					expectedTestScenarios[0].TestCategoryID,
					expectedTestScenarios[0].MotherServiceID,
					expectedTestScenarios[0].Status,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].ExecutionDuration,
					expectedTestScenarios[0].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[1].ID,
					expectedTestScenarios[1].CreatedAt,
					expectedTestScenarios[1].UpdatedAt,
					nil,
					expectedTestScenarios[1].Name,
					expectedTestScenarios[1].TestCategoryID,
					expectedTestScenarios[1].MotherServiceID,
					expectedTestScenarios[1].Status,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].ExecutionDuration,
					expectedTestScenarios[1].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[2].ID,
					expectedTestScenarios[2].CreatedAt,
					expectedTestScenarios[2].UpdatedAt,
					nil,
					expectedTestScenarios[2].Name,
					expectedTestScenarios[2].TestCategoryID,
					expectedTestScenarios[2].MotherServiceID,
					expectedTestScenarios[2].Status,
					expectedTestScenarios[2].MaxTestServiceCount,
					expectedTestScenarios[2].ExecutionDuration,
					expectedTestScenarios[2].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[3].ID,
					expectedTestScenarios[3].CreatedAt,
					expectedTestScenarios[3].UpdatedAt,
					nil,
					expectedTestScenarios[3].Name,
					expectedTestScenarios[3].TestCategoryID,
					expectedTestScenarios[3].MotherServiceID,
					expectedTestScenarios[3].Status,
					expectedTestScenarios[3].MaxTestServiceCount,
					expectedTestScenarios[3].ExecutionDuration,
					expectedTestScenarios[3].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[4].ID,
					expectedTestScenarios[4].CreatedAt,
					expectedTestScenarios[4].UpdatedAt,
					nil,
					expectedTestScenarios[4].Name,
					expectedTestScenarios[4].TestCategoryID,
					expectedTestScenarios[4].MotherServiceID,
					expectedTestScenarios[4].Status,
					expectedTestScenarios[4].MaxTestServiceCount,
					expectedTestScenarios[4].ExecutionDuration,
					expectedTestScenarios[4].AutoStepChangeRate,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 5)

		assert.Equal(t, expectedTestScenarios[0].ID, result[0].ID)
		assert.Equal(t, expectedTestScenarios[0].Name, result[0].Name)

		assert.Equal(t, expectedTestScenarios[1].ID, result[1].ID)
		assert.Equal(t, expectedTestScenarios[1].Name, result[1].Name)

		assert.Equal(t, expectedTestScenarios[2].ID, result[2].ID)
		assert.Equal(t, expectedTestScenarios[2].Name, result[2].Name)

		assert.Equal(t, expectedTestScenarios[3].ID, result[3].ID)
		assert.Equal(t, expectedTestScenarios[3].Name, result[3].Name)

		assert.Equal(t, expectedTestScenarios[4].ID, result[4].ID)
		assert.Equal(t, expectedTestScenarios[4].Name, result[4].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case - page 0 per page 0 - return all records", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenarios := []*entity.TestScenario{
			{
				ID:                  uint64(6),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load1",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(5),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load2",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load3",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(3),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load4",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:                  uint64(2),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load5",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
		}

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    0,
			PerPage: 0,
		} // no limit no offset - return all records

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}).
				AddRow(
					expectedTestScenarios[0].ID,
					expectedTestScenarios[0].CreatedAt,
					expectedTestScenarios[0].UpdatedAt,
					nil,
					expectedTestScenarios[0].Name,
					expectedTestScenarios[0].TestCategoryID,
					expectedTestScenarios[0].MotherServiceID,
					expectedTestScenarios[0].Status,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].ExecutionDuration,
					expectedTestScenarios[0].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[1].ID,
					expectedTestScenarios[1].CreatedAt,
					expectedTestScenarios[1].UpdatedAt,
					nil,
					expectedTestScenarios[1].Name,
					expectedTestScenarios[1].TestCategoryID,
					expectedTestScenarios[1].MotherServiceID,
					expectedTestScenarios[1].Status,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].ExecutionDuration,
					expectedTestScenarios[1].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[2].ID,
					expectedTestScenarios[2].CreatedAt,
					expectedTestScenarios[2].UpdatedAt,
					nil,
					expectedTestScenarios[2].Name,
					expectedTestScenarios[2].TestCategoryID,
					expectedTestScenarios[2].MotherServiceID,
					expectedTestScenarios[2].Status,
					expectedTestScenarios[2].MaxTestServiceCount,
					expectedTestScenarios[2].ExecutionDuration,
					expectedTestScenarios[2].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[3].ID,
					expectedTestScenarios[3].CreatedAt,
					expectedTestScenarios[3].UpdatedAt,
					nil,
					expectedTestScenarios[3].Name,
					expectedTestScenarios[3].TestCategoryID,
					expectedTestScenarios[3].MotherServiceID,
					expectedTestScenarios[3].Status,
					expectedTestScenarios[3].MaxTestServiceCount,
					expectedTestScenarios[3].ExecutionDuration,
					expectedTestScenarios[3].AutoStepChangeRate,
				).
				AddRow(
					expectedTestScenarios[4].ID,
					expectedTestScenarios[4].CreatedAt,
					expectedTestScenarios[4].UpdatedAt,
					nil,
					expectedTestScenarios[4].Name,
					expectedTestScenarios[4].TestCategoryID,
					expectedTestScenarios[4].MotherServiceID,
					expectedTestScenarios[4].Status,
					expectedTestScenarios[4].MaxTestServiceCount,
					expectedTestScenarios[4].ExecutionDuration,
					expectedTestScenarios[4].AutoStepChangeRate,
				))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 5)

		assert.Equal(t, expectedTestScenarios[0].ID, result[0].ID)
		assert.Equal(t, expectedTestScenarios[0].Name, result[0].Name)

		assert.Equal(t, expectedTestScenarios[1].ID, result[1].ID)
		assert.Equal(t, expectedTestScenarios[1].Name, result[1].Name)

		assert.Equal(t, expectedTestScenarios[2].ID, result[2].ID)
		assert.Equal(t, expectedTestScenarios[2].Name, result[2].Name)

		assert.Equal(t, expectedTestScenarios[3].ID, result[3].ID)
		assert.Equal(t, expectedTestScenarios[3].Name, result[3].Name)

		assert.Equal(t, expectedTestScenarios[4].ID, result[4].ID)
		assert.Equal(t, expectedTestScenarios[4].Name, result[4].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case - empty result - page beyond available data", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    100,
			PerPage: 10,
		} // big number page - LIMIT 10 OFFSET 990

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(10, 990).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - database error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(2, 2).
			WillReturnError(errors.New("database connection failed"))

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenarios)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed case - negative page", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, _, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    -2,
			PerPage: 2,
		}

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrNegativePageOrPerPageNotAllowed)
	})

	t.Run("failed case - negative per page", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, _, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: -2,
		}

		result, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrNegativePageOrPerPageNotAllowed)
	})
}
