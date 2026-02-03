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

		someTime := time.Date(2026, 01, 13, 14, 10, 0, 0, time.Now().Location())
		num := 10
		expectedTestScenario := &entity.TestScenario{
			ID:                  1,
			CreatedAt:           someTime,
			UpdatedAt:           someTime,
			Name:                "some test",
			Status:              entity.ScenarioStatusPending,
			MaxTestServiceCount: nil,
			ExecutionDuration:   nil,
			AutoStepChangeRate:  nil,
			TestCategoryID:      3,
			MotherServiceID:     2,
			TestCategory: &entity.TestCategory{
				ID:                     3,
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				Name:                   "peak",
				Label:                  "peak test",
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
			MotherService: &entity.MotherService{
				ID:   2,
				Name: "m2",
			},
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                    100,
				TestScenarioID:        1,
				CreatedAt:             someTime,
				UpdatedAt:             someTime,
				MaxRequests:           100,
				MaxDuration:           0,
				RequestDelayDuration:  nil,
				RandomRequestDelayMin: nil,
				RandomRequestDelayMax: nil,
				FixedTestNumber:       &num,
				RandomTestNumberMin:   nil,
				RandomTestNumberMax:   nil,
				BadValueRate:          0,
				NegativeValueRate:     0,
				RealValueRate:         0,
				ZeroValueRate:         0,
				StringValueRate:       0,
				LongStringValueRate:   0,
				NullValueRate:         0,
			},
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL`)).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
		}).AddRow(
			2,
			"m2",
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" = $1`)).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_change_rate",
		}).AddRow(
			3,
			someTime,
			someTime,
			"peak",
			"peak test",
			true,
			true,
			false,
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_service_configs" WHERE "test_service_configs"."test_scenario_id" = $1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"test_scenario_id",
			"max_requests",
			"max_duration",
			"request_delay_duration",
			"random_request_delay_min",
			"random_request_delay_max",
			"fixed_test_number",
			"random_test_number_min",
			"random_test_number_max",
			"bad_value_rate",
			"negative_value_rate",
			"real_value_rate",
			"zero_value_rate",
			"string_value_rate",
			"long_string_value_rate",
			"null_value_rate",
		}).AddRow(
			100,
			someTime,
			someTime,
			1,
			100,
			0,
			nil,
			nil,
			nil,
			10,
			nil,
			nil,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		))

		result, err := repo.GetByID(context.Background(), uint64(1))

		require.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario, result)
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
	t.Run("failed case - failed to get database records count", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		paginationRequest := entity.TestScenarioPaginationRequest{
			Page:    1,
			PerPage: 2,
		} // LIMIT 2 (no OFFSET because offset=0)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnError(errors.New("error happened"))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get test scenario records total count")
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("success case - page 1 per page 2", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenarios := []*entity.TestScenario{
			{
				ID:              5,
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "the test",
				TestCategoryID:  4,
				MotherServiceID: 3,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(4),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "new test",
				TestCategoryID:  uint64(5),
				MotherServiceID: uint64(4),
				TestCategory: &entity.TestCategory{
					ID:                     5,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "spike",
					Label:                  "spike test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   4,
					Name: "m4",
				},
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
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).AddRow(2),
			)

		mock.MatchExpectationsInOrder(false)
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" IN ($1,$2) AND "mother_services"."deleted_at" IS NULL`)).WithArgs(3, 4).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
		}).AddRow(
			3,
			"m3",
		).AddRow(
			4,
			"m4",
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" IN ($1,$2)`)).WithArgs(4, 5).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_change_rate",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
			false,
		).AddRow(
			5,
			now,
			now,
			"spike",
			"spike test",
			true,
			true,
			false,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, count, int64(2))
		assert.Len(t, result, 2)
		assert.Equal(t, expectedTestScenarios, result)
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
				ID:              uint64(5),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load1",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(4),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "smoke1",
				TestCategoryID:  uint64(5),
				MotherServiceID: uint64(4),
				TestCategory: &entity.TestCategory{
					ID:                     5,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "spike",
					Label:                  "spike test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   4,
					Name: "m4",
				},
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
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).AddRow(2),
			)

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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" IN ($1,$2) AND "mother_services"."deleted_at" IS NULL`)).WithArgs(3, 4).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
		}).AddRow(
			3,
			"m3",
		).AddRow(
			4,
			"m4",
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" IN ($1,$2)`)).WithArgs(4, 5).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_change_rate",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
			false,
		).AddRow(
			5,
			now,
			now,
			"spike",
			"spike test",
			true,
			true,
			false,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, count, int64(2))
		assert.Len(t, result, 2)

		assert.Equal(t, expectedTestScenarios, result)

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
				ID:              uint64(5),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load1",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
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
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).AddRow(1),
			)

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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL`)).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
		}).AddRow(
			3,
			"m3",
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" = $1`)).WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_change_rate",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
			false,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, count, int64(1))

		assert.Equal(t, expectedTestScenarios, result)

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
				ID:              uint64(6),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load1",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(5),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load2",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(4),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load3",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(3),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load4",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(2),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load5",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
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
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).AddRow(5),
			)

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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL`)).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
		}).AddRow(
			3,
			"m3",
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" = $1`)).WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_change_rate",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
			false,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 5)
		assert.Equal(t, count, int64(5))

		assert.Equal(t, expectedTestScenarios, result)

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
				ID:              uint64(6),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load1",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(5),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load2",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(4),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load3",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(3),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load4",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: nil,
				ExecutionDuration:   nil,
				AutoStepChangeRate:  nil,
			},
			{
				ID:              uint64(2),
				CreatedAt:       now,
				UpdatedAt:       now,
				Name:            "load5",
				TestCategoryID:  uint64(4),
				MotherServiceID: uint64(3),
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
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
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).AddRow(10),
			)

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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL`)).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
		}).AddRow(
			3,
			"m3",
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" = $1`)).WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_execution_duration",
			"has_auto_step_change_rate",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
			false,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 5)
		assert.Equal(t, count, int64(10))

		assert.Equal(t, expectedTestScenarios, result)

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
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).AddRow(0),
			)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(10, 990).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"max_test_service_count", "execution_duration",
				"auto_step_change_rate",
			}))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Equal(t, count, int64(0))
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
			`SELECT count(*) FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL`,
		)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).AddRow(0),
			)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`)).
			WithArgs(2, 2).
			WillReturnError(errors.New("database connection failed"))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))
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

		result, _, err := repo.GetPaginated(context.Background(), paginationRequest)

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

		result, _, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrNegativePageOrPerPageNotAllowed)
	})
}

func TestTestScenarioRepository_SetStatus(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				entity.ScenarioStatusRunning,
				sqlmock.AnyArg(),
				uint64(1),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err = repo.SetStatus(context.Background(), uint64(1), entity.ScenarioStatusRunning)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("failed case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				entity.ScenarioStatusRunning,
				sqlmock.AnyArg(),
				uint64(1),
			).
			WillReturnError(errors.New("something went wrong"))

		mock.ExpectRollback()

		err = repo.SetStatus(context.Background(), uint64(1), entity.ScenarioStatusRunning)

		assert.Error(t, err)
	})

}
