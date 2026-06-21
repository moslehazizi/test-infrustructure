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

	assert.NotNil(t, repo.db)
}

func TestTestScenarioRepository_Create(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
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
			Status:              entity.ScenarioStatus(entity.ScenarioStatusReady),
			MaxTestServiceCount: nil,
			StartedAt:           nil,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			// `INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","max_test_service_count","deployment_number","started_at","editable") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`)).
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","num_steps","max_test_service_count","deployment_number","started_at","editable","increase_agent_number","execution_number_multi_agent") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`)).
			WithArgs(
				testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				testScenario.NumSteps,
				nil,
				int32(0),
				testScenario.StartedAt,
				true,
				testScenario.IncreaseAgentNumber,
				testScenario.ExecNumMultiAgent).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		id, err := repo.Create(context.Background(), testScenario)

		assert.NoError(t, err)
		assert.Equal(t, id, uint64(1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case", func(t *testing.T) {
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
			Status:              entity.ScenarioStatus(entity.ScenarioStatusReady),
			MaxTestServiceCount: nil,
			StartedAt:           nil,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			// `INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","max_test_service_count","deployment_number","started_at","editable") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`)).
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","num_steps","max_test_service_count","deployment_number","started_at","editable","increase_agent_number","execution_number_multi_agent") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`)).
			WithArgs(testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				testScenario.NumSteps,
				nil, int32(0), testScenario.StartedAt,
				true,
				testScenario.IncreaseAgentNumber,
				testScenario.ExecNumMultiAgent).
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
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		sampleName := "test-name"

		someTime := time.Date(2026, 01, 13, 14, 10, 0, 0, time.Now().Location())
		num := 10
		expectedTestScenario := &entity.TestScenario{
			ID:                  1,
			CreatedAt:           someTime,
			UpdatedAt:           someTime,
			Name:                "some test",
			Status:              entity.ScenarioStatusReady,
			MaxTestServiceCount: nil,
			TestCategoryID:      3,
			MotherServiceID:     2,
			StartedAt:           nil,
			Editable:            true,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestCategory: &entity.TestCategory{
				ID:                     3,
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				Name:                   "peak",
				Label:                  "peak test",
				HasMaxTestServiceCount: true,
				HasNumSteps:            true,
			},
			MotherService: &entity.MotherService{
				ID:   2,
				Name: "m2",
			},
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     100,
				TestScenarioID:         1,
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				MaxRequests:            100,
				MaxDuration:            0,
				RequestDelayDuration:   nil,
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        &num,
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           0,
				NegativeValueRate:      0,
				RealValueRate:          0,
				ZeroValueRate:          0,
				StringValueRate:        0,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           sampleName,
				DatabaseTableName:      sampleName,
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "test_scenarios" WHERE "test_scenarios"."id" = $1 AND "test_scenarios"."deleted_at" IS NULL ORDER BY "test_scenarios"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
				"created_at",
				"updated_at",
				"deleted_at",
				"name",
				"test_category_id",
				"mother_service_id",
				"status",
				"num_steps",
				"max_test_service_count",
				"started_at",
				"editable",
				"increase_agent_number",
				"execution_number_multi_agent",
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
					expectedTestScenario.NumSteps,
					expectedTestScenario.MaxTestServiceCount,
					expectedTestScenario.StartedAt,
					expectedTestScenario.Editable,
					expectedTestScenario.IncreaseAgentNumber,
					expectedTestScenario.ExecNumMultiAgent,
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
			"has_num_steps",
		}).AddRow(
			3,
			someTime,
			someTime,
			"peak",
			"peak test",
			true,
			true,
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
			"database_name",
			"database_table_name",
			"increase_fixed_input",
			"execution_number_multi_fixed_input",
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
			sampleName,
			sampleName,
			expectedTestScenario.TestServiceConfig.IncreaseFixedInput,
			expectedTestScenario.TestServiceConfig.ExecNumMultiFixedInput,
		))

		result, err := repo.GetByID(context.Background(), uint64(1))

		require.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario, result)
	})

	t.Run("failed_case_record_not_found", func(t *testing.T) {
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

	t.Run("failed_case_database_error", func(t *testing.T) {
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
	t.Run("failed_case_failed_to_get_database_records_count", func(t *testing.T) {
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

	t.Run("success_case_page_1_per_page_2", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		now := time.Now()

		expectedTestScenarios := []*entity.TestScenario{
			{
				ID:                  5,
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "the test",
				TestCategoryID:      4,
				MotherServiceID:     3,
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "new test",
				TestCategoryID:      uint64(5),
				MotherServiceID:     uint64(4),
				StartedAt:           nil,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     5,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "spike",
					Label:                  "spike test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   4,
					Name: "m4",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
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
				"id",
				"created_at",
				"updated_at",
				"deleted_at",
				"name",
				"test_category_id",
				"mother_service_id",
				"status",
				"num_steps",
				"max_test_service_count",
				"started_at",
				"editable",
				"increase_agent_number",
				"execution_number_multi_agent",
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
					expectedTestScenarios[0].NumSteps,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].StartedAt,
					expectedTestScenarios[0].Editable,
					expectedTestScenarios[0].IncreaseAgentNumber,
					expectedTestScenarios[0].ExecNumMultiAgent,
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
					expectedTestScenarios[1].NumSteps,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].StartedAt,
					expectedTestScenarios[1].Editable,
					expectedTestScenarios[1].IncreaseAgentNumber,
					expectedTestScenarios[1].ExecNumMultiAgent,
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
			"has_num_steps",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
		).AddRow(
			5,
			now,
			now,
			"spike",
			"spike test",
			true,
			true,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, count, int64(2))
		assert.Len(t, result, 2)
		assert.Equal(t, expectedTestScenarios, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_case_page_2_per_page_2", func(t *testing.T) {
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
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "smoke1",
				TestCategoryID:      uint64(5),
				MotherServiceID:     uint64(4),
				StartedAt:           nil,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     5,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "spike",
					Label:                  "spike test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   4,
					Name: "m4",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
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
				"num_steps",
				"max_test_service_count", "started_at", "editable",
				"increase_agent_number", "execution_number_multi_agent",
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
					expectedTestScenarios[0].NumSteps,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].StartedAt,
					expectedTestScenarios[0].Editable,
					expectedTestScenarios[0].IncreaseAgentNumber,
					expectedTestScenarios[0].ExecNumMultiAgent,
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
					expectedTestScenarios[1].NumSteps,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].StartedAt,
					expectedTestScenarios[1].Editable,
					expectedTestScenarios[1].IncreaseAgentNumber,
					expectedTestScenarios[1].ExecNumMultiAgent,
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
			"has_num_steps",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
		).AddRow(
			5,
			now,
			now,
			"spike",
			"spike test",
			true,
			true,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, count, int64(2))
		assert.Len(t, result, 2)

		assert.Equal(t, expectedTestScenarios, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_case_page_2_per_page_1", func(t *testing.T) {
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
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
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
				"num_steps",
				"max_test_service_count", "started_at", "editable", "increase_agent_number",
				"execution_number_multi_agent",
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
					expectedTestScenarios[0].NumSteps,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].StartedAt,
					expectedTestScenarios[0].Editable,
					expectedTestScenarios[0].IncreaseAgentNumber,
					expectedTestScenarios[0].ExecNumMultiAgent,
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
			"has_num_steps",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, count, int64(1))

		assert.Equal(t, expectedTestScenarios, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_case_page_3_per_page_5", func(t *testing.T) {
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
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(5),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load2",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load3",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(3),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load4",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(2),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load5",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				StartedAt:           nil,
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
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
				"num_steps",
				"max_test_service_count", "started_at", "editable", "increase_agent_number",
				"execution_number_multi_agent",
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
					expectedTestScenarios[1].NumSteps,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].StartedAt,
					expectedTestScenarios[0].Editable,
					expectedTestScenarios[0].IncreaseAgentNumber,
					expectedTestScenarios[0].ExecNumMultiAgent,
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
					expectedTestScenarios[1].NumSteps,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].StartedAt,
					expectedTestScenarios[1].Editable,
					expectedTestScenarios[1].IncreaseAgentNumber,
					expectedTestScenarios[1].ExecNumMultiAgent,
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
					expectedTestScenarios[2].NumSteps,
					expectedTestScenarios[2].MaxTestServiceCount,
					expectedTestScenarios[2].StartedAt,
					expectedTestScenarios[2].Editable,
					expectedTestScenarios[2].IncreaseAgentNumber,
					expectedTestScenarios[2].ExecNumMultiAgent,
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
					expectedTestScenarios[3].NumSteps,
					expectedTestScenarios[3].MaxTestServiceCount,
					expectedTestScenarios[3].StartedAt,
					expectedTestScenarios[3].Editable,
					expectedTestScenarios[3].IncreaseAgentNumber,
					expectedTestScenarios[3].ExecNumMultiAgent,
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
					expectedTestScenarios[4].NumSteps,
					expectedTestScenarios[4].MaxTestServiceCount,
					expectedTestScenarios[4].StartedAt,
					expectedTestScenarios[4].Editable,
					expectedTestScenarios[4].IncreaseAgentNumber,
					expectedTestScenarios[4].ExecNumMultiAgent,
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
			"has_num_steps",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 5)
		assert.Equal(t, count, int64(5))

		assert.Equal(t, expectedTestScenarios, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_case_page_0_per_page_0_return_all_records", func(t *testing.T) {
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
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(5),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load2",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(4),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load3",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(3),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load4",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
			},
			{
				ID:                  uint64(2),
				CreatedAt:           now,
				UpdatedAt:           now,
				Name:                "load5",
				TestCategoryID:      uint64(4),
				MotherServiceID:     uint64(3),
				NumSteps:            2,
				IncreaseAgentNumber: 0,
				ExecNumMultiAgent:   1,
				TestCategory: &entity.TestCategory{
					ID:                     4,
					CreatedAt:              now,
					UpdatedAt:              now,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   3,
					Name: "m3",
				},
				Status:              entity.ScenarioStatusReady,
				MaxTestServiceCount: nil,
				Editable:            true,
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
				"num_steps",
				"max_test_service_count", "started_at", "editable", "increase_agent_number", "execution_number_multi_agent",
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
					expectedTestScenarios[4].NumSteps,
					expectedTestScenarios[0].MaxTestServiceCount,
					expectedTestScenarios[0].StartedAt,
					expectedTestScenarios[0].Editable,
					expectedTestScenarios[0].IncreaseAgentNumber,
					expectedTestScenarios[0].ExecNumMultiAgent,
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
					expectedTestScenarios[4].NumSteps,
					expectedTestScenarios[1].MaxTestServiceCount,
					expectedTestScenarios[1].StartedAt,
					expectedTestScenarios[1].Editable,
					expectedTestScenarios[1].IncreaseAgentNumber,
					expectedTestScenarios[1].ExecNumMultiAgent,
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
					expectedTestScenarios[4].NumSteps,
					expectedTestScenarios[2].MaxTestServiceCount,
					expectedTestScenarios[2].StartedAt,
					expectedTestScenarios[2].Editable,
					expectedTestScenarios[2].IncreaseAgentNumber,
					expectedTestScenarios[2].ExecNumMultiAgent,
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
					expectedTestScenarios[4].NumSteps,
					expectedTestScenarios[3].MaxTestServiceCount,
					expectedTestScenarios[3].StartedAt,
					expectedTestScenarios[3].Editable,
					expectedTestScenarios[3].IncreaseAgentNumber,
					expectedTestScenarios[3].ExecNumMultiAgent,
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
					expectedTestScenarios[4].NumSteps,
					expectedTestScenarios[4].MaxTestServiceCount,
					expectedTestScenarios[4].StartedAt,
					expectedTestScenarios[4].Editable,
					expectedTestScenarios[4].IncreaseAgentNumber,
					expectedTestScenarios[4].ExecNumMultiAgent,
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
			"has_num_steps",
		}).AddRow(
			4,
			now,
			now,
			"peak",
			"peak test",
			true,
			true,
		))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 5)
		assert.Equal(t, count, int64(10))

		assert.Equal(t, expectedTestScenarios, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_case_empty_result_page_beyond_available_data", func(t *testing.T) {
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
				"max_test_service_count", "started_at", "editable",
			}))

		result, count, err := repo.GetPaginated(context.Background(), paginationRequest)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Equal(t, count, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case_database_error", func(t *testing.T) {
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

	t.Run("failed_case_negative_page", func(t *testing.T) {
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

	t.Run("failed_case_negative_per_page", func(t *testing.T) {
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
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "editable"=$1,"status"=$2,"updated_at"=$3 WHERE "id" = $4 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				false,
				entity.ScenarioStatusRunning,
				sqlmock.AnyArg(),
				uint64(1),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err = repo.SetStatus(context.Background(), uint64(1), entity.ScenarioStatusRunning, false)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("failed_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "editable"=$2,"status"=$1,"updated_at"=$3 WHERE "id" = $4 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				false,
				entity.ScenarioStatusRunning,
				sqlmock.AnyArg(),
				uint64(1),
			).
			WillReturnError(errors.New("something went wrong"))

		mock.ExpectRollback()

		err = repo.SetStatus(context.Background(), uint64(1), entity.ScenarioStatusRunning, false)

		assert.Error(t, err)
	})

}

func TestTestScenarioRepository_GetByStatus(t *testing.T) {
	t.Run("success_case_get_running_status", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		someTime := time.Date(2026, 01, 13, 14, 10, 0, 0, time.Now().Location())
		num := 10

		expectedTestScenario := []*entity.TestScenario{
			{
				ID:                  1,
				CreatedAt:           someTime,
				UpdatedAt:           someTime,
				Name:                "some test",
				Status:              entity.ScenarioStatusRunning,
				MaxTestServiceCount: nil,
				TestCategoryID:      3,
				MotherServiceID:     2,
				NumSteps:            2,
				TestCategory: &entity.TestCategory{
					ID:                     3,
					CreatedAt:              someTime,
					UpdatedAt:              someTime,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
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
				Editable: false,
			},
			{
				ID:                  4,
				CreatedAt:           someTime,
				UpdatedAt:           someTime,
				Name:                "some test",
				Status:              entity.ScenarioStatusRunning,
				MaxTestServiceCount: nil,
				TestCategoryID:      6,
				MotherServiceID:     5,
				NumSteps:            2,
				TestCategory: &entity.TestCategory{
					ID:                     6,
					CreatedAt:              someTime,
					UpdatedAt:              someTime,
					Name:                   "spike",
					Label:                  "spike test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   5,
					Name: "m5",
				},
				TestServiceConfig: &entity.TestServiceConfig{
					ID:                    54,
					TestScenarioID:        4,
					CreatedAt:             someTime,
					UpdatedAt:             someTime,
					MaxRequests:           54,
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
				Editable: false},
		}

		mock.ExpectQuery(`SELECT \* FROM "test_scenarios" WHERE status = .+ AND "test_scenarios"\."deleted_at" IS NULL ORDER BY id DESC`).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"num_steps",
				"max_test_service_count", "editable",
			}).
				AddRow(
					expectedTestScenario[0].ID,
					expectedTestScenario[0].CreatedAt,
					expectedTestScenario[0].UpdatedAt,
					nil,
					expectedTestScenario[0].Name,
					expectedTestScenario[0].TestCategoryID,
					expectedTestScenario[0].MotherServiceID,
					expectedTestScenario[0].Status,
					expectedTestScenario[0].NumSteps,
					expectedTestScenario[0].MaxTestServiceCount,
					expectedTestScenario[0].Editable,
				).
				AddRow(
					expectedTestScenario[1].ID,
					expectedTestScenario[1].CreatedAt,
					expectedTestScenario[1].UpdatedAt,
					nil,
					expectedTestScenario[1].Name,
					expectedTestScenario[1].TestCategoryID,
					expectedTestScenario[1].MotherServiceID,
					expectedTestScenario[1].Status,
					expectedTestScenario[1].NumSteps,
					expectedTestScenario[1].MaxTestServiceCount,
					expectedTestScenario[1].Editable,
				))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" IN ($1,$2) AND "mother_services"."deleted_at" IS NULL`)).WithArgs(2, 5).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
		}).AddRow(
			2,
			"m2",
		).AddRow(
			5,
			"m5",
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" IN ($1,$2)`)).WithArgs(3, 6).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_num_steps",
		}).AddRow(
			3,
			someTime,
			someTime,
			"peak",
			"peak test",
			true,
			true,
		).AddRow(
			6,
			someTime,
			someTime,
			"spike",
			"spike test",
			true,
			true,
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_service_configs" WHERE "test_service_configs"."test_scenario_id" IN ($1,$2)`)).WithArgs(1, 4).WillReturnRows(sqlmock.NewRows([]string{
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
		).AddRow(
			54,
			someTime,
			someTime,
			4,
			54,
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

		result, err := repo.GetByStatus(context.Background(), entity.ScenarioStatusRunning)

		require.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario, result)
	})

	t.Run("failed_case_database_connection_error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)

		mock.ExpectQuery(`SELECT \* FROM "test_scenarios" WHERE status = .+ AND "test_scenarios"\."deleted_at" IS NULL ORDER BY id DESC`).
			WillReturnError(errors.New("error happened"))

		result, err := repo.GetByStatus(context.Background(), entity.ScenarioStatusRunning)

		require.NoError(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenariosByStatus)
	})
}

func TestGetDeploymentNumberByScenarioID(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT "deployment_number" FROM "test_scenarios" WHERE id = $1 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(uint64(1)).
			WillReturnRows(
				sqlmock.NewRows([]string{"deployment_number"}).
					AddRow(5),
			)

		result, err := repo.GetDeploymentNumberByScenarioID(context.Background(), 1)

		require.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, int32(5), result)
	})

	t.Run("failed_case_record_not_found", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT "deployment_number" FROM "test_scenarios" WHERE id = $1 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(uint64(999)).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetDeploymentNumberByScenarioID(context.Background(), 999)

		require.NoError(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		assert.Equal(t, int32(0), result)
	})

	t.Run("failed_case_database_error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT "deployment_number" FROM "test_scenarios" WHERE id = $1 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(uint64(1)).
			WillReturnError(errors.New("db is down"))

		result, err := repo.GetDeploymentNumberByScenarioID(context.Background(), 1)

		require.NoError(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		assert.Equal(t, int32(0), result)
	})
}

func TestUpdateDeploymentNumber(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "deployment_number"=$1,"updated_at"=$2 WHERE id = $3 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(int32(42), sqlmock.AnyArg(), uint64(1)). // ✅ AnyArg for timestamp
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err = repo.UpdateDeploymentNumber(context.Background(), 1, 42)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("record_not_found", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "deployment_number"=$1,"updated_at"=$2 WHERE id = $3 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(int32(42), sqlmock.AnyArg(), uint64(999)).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err = repo.UpdateDeploymentNumber(context.Background(), 999, 42)
		require.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "deployment_number"=$1,"updated_at"=$2 WHERE id = $3 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(int32(42), sqlmock.AnyArg(), uint64(1)).
			WillReturnError(errors.New("db is down"))
		mock.ExpectRollback()

		err = repo.UpdateDeploymentNumber(context.Background(), 1, 42)
		require.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestTestScenarioRepository_Update(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		maxCount := int64(10)

		scenario := &entity.TestScenario{
			ID:                  1,
			Name:                "updated-name",
			MotherServiceID:     2,
			Status:              entity.ScenarioStatusRunning,
			MaxTestServiceCount: &maxCount,
			DeploymentNumber:    3,
			Editable:            false,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "execution_number_multi_agent"=$1,"increase_agent_number"=$2,"max_test_service_count"=$3,"mother_service_id"=$4,"name"=$5,"num_steps"=$6,"status"=$7,"updated_at"=$8 WHERE id = $9 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				scenario.ExecNumMultiAgent,
				scenario.IncreaseAgentNumber,
				scenario.MaxTestServiceCount,
				scenario.MotherServiceID,
				scenario.Name,
				scenario.NumSteps,
				scenario.Status,
				sqlmock.AnyArg(),
				scenario.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err = repo.Update(context.Background(), scenario)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update_fails", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewTestScenarioRepository(db)

		scenario := &entity.TestScenario{
			ID:                  1,
			Name:                "updated-name",
			MotherServiceID:     2,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET`,
		)).
			WillReturnError(errors.New("update failed"))

		mock.ExpectRollback()

		err = repo.Update(context.Background(), scenario)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTestScenarioRepository_GetByMotherServiceId(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)
		motherId := uint64(1)

		someTime := time.Date(2026, 01, 13, 14, 10, 0, 0, time.Now().Location())
		num := 10

		expectedTestScenario := []*entity.TestScenario{
			{
				ID:                  1,
				CreatedAt:           someTime,
				UpdatedAt:           someTime,
				Name:                "some test",
				Status:              entity.ScenarioStatusRunning,
				MaxTestServiceCount: nil,
				TestCategoryID:      3,
				MotherServiceID:     motherId,
				NumSteps:            2,
				TestCategory: &entity.TestCategory{
					ID:                     3,
					CreatedAt:              someTime,
					UpdatedAt:              someTime,
					Name:                   "peak",
					Label:                  "peak test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   motherId,
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
				Editable: false,
			},
			{
				ID:                  4,
				CreatedAt:           someTime,
				UpdatedAt:           someTime,
				Name:                "some test",
				Status:              entity.ScenarioStatusRunning,
				MaxTestServiceCount: nil,
				TestCategoryID:      6,
				MotherServiceID:     motherId,
				NumSteps:            2,
				TestCategory: &entity.TestCategory{
					ID:                     6,
					CreatedAt:              someTime,
					UpdatedAt:              someTime,
					Name:                   "spike",
					Label:                  "spike test",
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   motherId,
					Name: "m2",
				},
				TestServiceConfig: &entity.TestServiceConfig{
					ID:                    54,
					TestScenarioID:        4,
					CreatedAt:             someTime,
					UpdatedAt:             someTime,
					MaxRequests:           54,
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
				Editable: false},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_scenarios" WHERE mother_service_id = $1 AND "test_scenarios"."deleted_at" IS NULL ORDER BY id DESC`)).
			WithArgs(motherId).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"test_category_id", "mother_service_id", "status",
				"num_steps",
				"max_test_service_count", "editable",
			}).
				AddRow(
					expectedTestScenario[0].ID,
					expectedTestScenario[0].CreatedAt,
					expectedTestScenario[0].UpdatedAt,
					nil,
					expectedTestScenario[0].Name,
					expectedTestScenario[0].TestCategoryID,
					expectedTestScenario[0].MotherServiceID,
					expectedTestScenario[0].Status,
					expectedTestScenario[0].NumSteps,
					expectedTestScenario[0].MaxTestServiceCount,
					expectedTestScenario[0].Editable,
				).
				AddRow(
					expectedTestScenario[1].ID,
					expectedTestScenario[1].CreatedAt,
					expectedTestScenario[1].UpdatedAt,
					nil,
					expectedTestScenario[1].Name,
					expectedTestScenario[1].TestCategoryID,
					expectedTestScenario[1].MotherServiceID,
					expectedTestScenario[1].Status,
					expectedTestScenario[1].NumSteps,
					expectedTestScenario[1].MaxTestServiceCount,
					expectedTestScenario[1].Editable,
				))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL`)).
			WithArgs(motherId).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
				"name",
			}).AddRow(
				motherId,
				"m2",
			))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_categories" WHERE "test_categories"."id" IN ($1,$2)`)).WithArgs(3, 6).WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"name",
			"label",
			"has_max_test_service_count",
			"has_num_steps",
		}).AddRow(
			3,
			someTime,
			someTime,
			"peak",
			"peak test",
			true,
			true,
		).AddRow(
			6,
			someTime,
			someTime,
			"spike",
			"spike test",
			true,
			true,
		))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_service_configs" WHERE "test_service_configs"."test_scenario_id" IN ($1,$2)`)).WithArgs(1, 4).WillReturnRows(sqlmock.NewRows([]string{
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
		).AddRow(
			54,
			someTime,
			someTime,
			4,
			54,
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

		result, err := repo.GetByMotherServiceId(context.Background(), motherId)

		require.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario, result)
	})

	t.Run("failed_case_database_connection_error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewTestScenarioRepository(db)
		motherId := uint64(1)

		mock.ExpectQuery(`SELECT \* FROM "test_scenarios" WHERE mother_service_id = .+ AND "test_scenarios"\."deleted_at" IS NULL ORDER BY id DESC`).
			WillReturnError(errors.New("error happened"))

		result, err := repo.GetByMotherServiceId(context.Background(), motherId)

		require.NoError(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenarios)
	})
}

func TestUpdateScenarioAndConfig(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		ctx := context.Background()

		repo := NewTestScenarioRepository(db)

		maxCount := int64(10)

		scenario := &entity.TestScenario{
			ID:                  1,
			Name:                "updated-name",
			MotherServiceID:     2,
			Status:              entity.ScenarioStatusRunning,
			MaxTestServiceCount: &maxCount,
			DeploymentNumber:    3,
			Editable:            false,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:            100,
				MaxDuration:            60,
				RequestDelayDuration:   new(1),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(1),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           1,
				NegativeValueRate:      2,
				RealValueRate:          3,
				ZeroValueRate:          4,
				StringValueRate:        5,
				LongStringValueRate:    6,
				NullValueRate:          7,
				DatabaseName:           "db1",
				DatabaseTableName:      "table1",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "execution_number_multi_agent"=$1,"increase_agent_number"=$2,"max_test_service_count"=$3,"mother_service_id"=$4,"name"=$5,"num_steps"=$6,"status"=$7,"updated_at"=$8 WHERE id = $9 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				scenario.ExecNumMultiAgent,
				scenario.IncreaseAgentNumber,
				scenario.MaxTestServiceCount,
				scenario.MotherServiceID,
				scenario.Name,
				scenario.NumSteps,
				scenario.Status,
				sqlmock.AnyArg(),
				scenario.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_service_configs" SET`,
		)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err = repo.UpdateScenarioAndConfig(ctx, scenario)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case_update_scenario", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		ctx := context.Background()

		repo := NewTestScenarioRepository(db)

		maxCount := int64(10)

		scenario := &entity.TestScenario{
			ID:                  1,
			Name:                "updated-name",
			MotherServiceID:     2,
			Status:              entity.ScenarioStatusRunning,
			MaxTestServiceCount: &maxCount,
			DeploymentNumber:    3,
			Editable:            false,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:            100,
				MaxDuration:            60,
				RequestDelayDuration:   new(1),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(1),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           1,
				NegativeValueRate:      2,
				RealValueRate:          3,
				ZeroValueRate:          4,
				StringValueRate:        5,
				LongStringValueRate:    6,
				NullValueRate:          7,
				DatabaseName:           "db1",
				DatabaseTableName:      "table1",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "execution_number_multi_agent"=$1,"increase_agent_number"=$2,"max_test_service_count"=$3,"mother_service_id"=$4,"name"=$5,"num_steps"=$6,"status"=$7,"updated_at"=$8 WHERE id = $9 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				scenario.ExecNumMultiAgent,
				scenario.IncreaseAgentNumber,
				scenario.MaxTestServiceCount,
				scenario.MotherServiceID,
				scenario.Name,
				scenario.NumSteps,
				scenario.Status,
				sqlmock.AnyArg(),
				scenario.ID,
			).
			WillReturnError(errors.New("something went wrong"))

		mock.ExpectRollback()

		err = repo.UpdateScenarioAndConfig(ctx, scenario)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to update test scenario:")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case_update_config", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		ctx := context.Background()

		repo := NewTestScenarioRepository(db)

		maxCount := int64(10)

		scenario := &entity.TestScenario{
			ID:                  1,
			Name:                "updated-name",
			MotherServiceID:     2,
			Status:              entity.ScenarioStatusRunning,
			MaxTestServiceCount: &maxCount,
			DeploymentNumber:    3,
			Editable:            false,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:            100,
				MaxDuration:            60,
				RequestDelayDuration:   new(1),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(1),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           1,
				NegativeValueRate:      2,
				RealValueRate:          3,
				ZeroValueRate:          4,
				StringValueRate:        5,
				LongStringValueRate:    6,
				NullValueRate:          7,
				DatabaseName:           "db1",
				DatabaseTableName:      "table1",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_scenarios" SET "execution_number_multi_agent"=$1,"increase_agent_number"=$2,"max_test_service_count"=$3,"mother_service_id"=$4,"name"=$5,"num_steps"=$6,"status"=$7,"updated_at"=$8 WHERE id = $9 AND "test_scenarios"."deleted_at" IS NULL`,
		)).
			WithArgs(
				scenario.ExecNumMultiAgent,
				scenario.IncreaseAgentNumber,
				scenario.MaxTestServiceCount,
				scenario.MotherServiceID,
				scenario.Name,
				scenario.NumSteps,
				scenario.Status,
				sqlmock.AnyArg(),
				scenario.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "test_service_configs" SET`,
		)).
			WillReturnError(errors.New("something went wrong"))

		mock.ExpectRollback()

		err = repo.UpdateScenarioAndConfig(ctx, scenario)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to update test service config:")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateScenarioAndConfig(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		ctx := context.Background()

		repo := NewTestScenarioRepository(db)
		now := time.Now()
		databaseName := "test_db"
		databaseTableName := "test_table"

		testScenario := &entity.TestScenario{
			CreatedAt:           now,
			UpdatedAt:           now,
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			Status:              entity.ScenarioStatus(entity.ScenarioStatusReady),
			MaxTestServiceCount: nil,
			StartedAt:           nil,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				TestScenarioID:         uint64(1),
				CreatedAt:              now,
				UpdatedAt:              now,
				MaxRequests:            12,
				MaxDuration:            10,
				DatabaseName:           databaseName,
				DatabaseTableName:      databaseTableName,
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","num_steps","max_test_service_count","deployment_number","started_at","editable","increase_agent_number","execution_number_multi_agent") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`)).
			WithArgs(
				testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				testScenario.NumSteps,
				nil,
				int32(0),
				testScenario.StartedAt,
				true,
				testScenario.IncreaseAgentNumber,
				testScenario.ExecNumMultiAgent).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_service_configs" ("test_scenario_id","max_requests","max_duration","request_delay_duration","random_request_delay_min","random_request_delay_max","fixed_test_number","random_test_number_min","random_test_number_max","bad_value_rate","negative_value_rate","real_value_rate","zero_value_rate","string_value_rate","long_string_value_rate","null_value_rate","created_at","updated_at","database_name","database_table_name","increase_fixed_input","execution_number_multi_fixed_input") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22) RETURNING "id"`)).
			WithArgs(
				testScenario.TestServiceConfig.TestScenarioID,
				testScenario.TestServiceConfig.MaxRequests,
				testScenario.TestServiceConfig.MaxDuration,
				nil, nil, nil, nil, nil, nil,
				0, 0, 0, 0, 0, 0, 0,
				testScenario.TestServiceConfig.CreatedAt,
				testScenario.TestServiceConfig.UpdatedAt,
				testScenario.TestServiceConfig.DatabaseName,
				testScenario.TestServiceConfig.DatabaseTableName,
				testScenario.TestServiceConfig.IncreaseFixedInput,
				testScenario.TestServiceConfig.ExecNumMultiFixedInput,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		err = repo.CreateScenarioAndConfig(ctx, testScenario)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case_create_scenario_error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		ctx := context.Background()

		repo := NewTestScenarioRepository(db)
		now := time.Now()
		databaseName := "test_db"
		databaseTableName := "test_table"

		testScenario := &entity.TestScenario{
			CreatedAt:           now,
			UpdatedAt:           now,
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			Status:              entity.ScenarioStatus(entity.ScenarioStatusReady),
			MaxTestServiceCount: nil,
			StartedAt:           nil,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				TestScenarioID:         uint64(1),
				CreatedAt:              now,
				UpdatedAt:              now,
				MaxRequests:            12,
				MaxDuration:            10,
				DatabaseName:           databaseName,
				DatabaseTableName:      databaseTableName,
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","num_steps","max_test_service_count","deployment_number","started_at","editable","increase_agent_number","execution_number_multi_agent") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`)).
			WithArgs(
				testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				testScenario.NumSteps,
				nil,
				int32(0),
				testScenario.StartedAt,
				true,
				testScenario.IncreaseAgentNumber,
				testScenario.ExecNumMultiAgent).
			WillReturnError(errors.New("something went wrong"))

		mock.ExpectRollback()

		err = repo.CreateScenarioAndConfig(ctx, testScenario)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to create test scenario record")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case_create_config_error", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		ctx := context.Background()

		repo := NewTestScenarioRepository(db)
		now := time.Now()
		databaseName := "test_db"
		databaseTableName := "test_table"

		testScenario := &entity.TestScenario{
			CreatedAt:           now,
			UpdatedAt:           now,
			DeletedAt:           nil,
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			Status:              entity.ScenarioStatus(entity.ScenarioStatusReady),
			MaxTestServiceCount: nil,
			StartedAt:           nil,
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				TestScenarioID:         uint64(1),
				CreatedAt:              now,
				UpdatedAt:              now,
				MaxRequests:            12,
				MaxDuration:            10,
				DatabaseName:           databaseName,
				DatabaseTableName:      databaseTableName,
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_scenarios" ("created_at","updated_at","deleted_at","name","test_category_id","mother_service_id","status","num_steps","max_test_service_count","deployment_number","started_at","editable","increase_agent_number","execution_number_multi_agent") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`)).
			WithArgs(
				testScenario.CreatedAt,
				testScenario.UpdatedAt,
				testScenario.DeletedAt,
				testScenario.Name,
				testScenario.TestCategoryID,
				testScenario.MotherServiceID,
				testScenario.Status,
				testScenario.NumSteps,
				nil,
				int32(0),
				testScenario.StartedAt,
				true,
				testScenario.IncreaseAgentNumber,
				testScenario.ExecNumMultiAgent).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "test_service_configs" ("test_scenario_id","max_requests","max_duration","request_delay_duration","random_request_delay_min","random_request_delay_max","fixed_test_number","random_test_number_min","random_test_number_max","bad_value_rate","negative_value_rate","real_value_rate","zero_value_rate","string_value_rate","long_string_value_rate","null_value_rate","created_at","updated_at","database_name","database_table_name","increase_fixed_input","execution_number_multi_fixed_input") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22) RETURNING "id"`)).
			WithArgs(
				testScenario.TestServiceConfig.TestScenarioID,
				testScenario.TestServiceConfig.MaxRequests,
				testScenario.TestServiceConfig.MaxDuration,
				nil, nil, nil, nil, nil, nil,
				0, 0, 0, 0, 0, 0, 0,
				testScenario.TestServiceConfig.CreatedAt,
				testScenario.TestServiceConfig.UpdatedAt,
				testScenario.TestServiceConfig.DatabaseName,
				testScenario.TestServiceConfig.DatabaseTableName,
				testScenario.TestServiceConfig.IncreaseFixedInput,
				testScenario.TestServiceConfig.ExecNumMultiFixedInput,
			).
			WillReturnError(errors.New("something went wrong"))

		mock.ExpectRollback()

		err = repo.CreateScenarioAndConfig(ctx, testScenario)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to create test service config record")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
