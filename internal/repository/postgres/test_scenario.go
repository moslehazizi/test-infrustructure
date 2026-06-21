package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewTestScenarioRepository(db database.Database) *testScenario {
	return &testScenario{
		db: db,
	}
}

type testScenario struct {
	db database.Database
}

func (repo *testScenario) Create(ctx context.Context, testSci *entity.TestScenario) (uint64, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "create_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "insert"), attribute.String("test_scenario.name", testSci.Name))

	err := postgres.QueryBuilder(ctx, repo.db).
		Omit(clause.Associations).
		Create(testSci).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return 0, fmt.Errorf("failed to create test scenario record: %w", err)
	}

	id := testSci.ID

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	return id, nil
}

func (repo *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "get_test_scenario_by_id")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	var testScenario entity.TestScenario
	err := postgres.QueryBuilder(ctx, repo.db).
		Preload("TestCategory").
		Preload("MotherService").
		Preload("TestServiceConfig").
		First(&testScenario, id).Error
	_ = err
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return nil, pkg.ErrTestScenarioNotFound
		}

		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(testScenario.ID, 10)))

	return &testScenario, nil
}

func (repo *testScenario) GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, int64, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "get_paginated_test_scenarios")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("pagination.page", strconv.Itoa(pagRequest.Page)), attribute.String("pagination.per_page", strconv.Itoa(pagRequest.PerPage)))

	var testScenarios []*entity.TestScenario
	if pagRequest.Page < 0 || pagRequest.PerPage < 0 {
		span.SetAttributes(attribute.String("error.type", "invalid_pagination"))

		return nil, 0, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenarios, pkg.ErrNegativePageOrPerPageNotAllowed)
	}

	query := postgres.QueryBuilder(ctx, repo.db).
		Preload("TestCategory").
		Preload("MotherService").
		Order("id DESC")

	var count int64

	if err := postgres.QueryBuilder(ctx, repo.db).Model(&entity.TestScenario{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get test scenario records total count: %w", err)
	}

	if pagRequest.Page > 0 && pagRequest.PerPage > 0 {
		offset := (pagRequest.Page - 1) * pagRequest.PerPage
		query = query.Limit(pagRequest.PerPage)
		if offset > 0 {
			query = query.Offset(offset)
		}
	}

	err := query.Find(&testScenarios).Error

	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, 0, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	return testScenarios, count, nil
}

func (repo *testScenario) SetStatus(ctx context.Context, id uint64, status entity.ScenarioStatus, editable bool) error {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "set_test_scenario_status")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "update"), attribute.String("test_scenario.id", strconv.FormatUint(id, 10)), attribute.String("test_scenario.status", string(status)))

	err := postgres.QueryBuilder(ctx, repo.db).
		Omit(clause.Associations).
		Model(&entity.TestScenario{}).
		Where("id", id).
		Updates(map[string]any{
			"status":     status,
			"editable":   editable,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("failed to update test scenario status: %w", err)
	}

	return nil
}

func (repo *testScenario) GetByStatus(ctx context.Context, status entity.ScenarioStatus) ([]*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "get_by_status_test_scenarios")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("status", string(status)))

	var testScenarios []*entity.TestScenario

	query := postgres.QueryBuilder(ctx, repo.db).
		Where("status = ?", status).
		Preload("TestCategory").
		Preload("MotherService").
		Preload("TestServiceConfig").
		Order("id DESC")

	err := query.Find(&testScenarios).Error

	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenariosByStatus, err)
	}

	return testScenarios, nil
}

func (repo *testScenario) GetDeploymentNumberByScenarioID(
	ctx context.Context,
	id uint64,
) (int32, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "get_deployment_number_by_scenario_id")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))
	span.SetAttributes(
		attribute.String("database.operation", "select"),
		attribute.String("test_scenario.id", strconv.FormatUint(id, 10)),
	)

	var deploymentNumber int32

	err := postgres.QueryBuilder(ctx, repo.db).
		Model(&entity.TestScenario{}).
		Select("deployment_number").
		Where("id = ?", id).
		Scan(&deploymentNumber).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return 0, pkg.ErrTestScenarioNotFound
		}

		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", err.Error()),
		)

		return 0, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	return deploymentNumber, nil
}

func (repo *testScenario) UpdateDeploymentNumber(
	ctx context.Context,
	id uint64,
	newDeploymentNumber int32,
) error {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "update_deployment_number")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))
	span.SetAttributes(
		attribute.String("database.operation", "update"),
		attribute.String("test_scenario.id", strconv.FormatUint(id, 10)),
	)

	result := postgres.QueryBuilder(ctx, repo.db).
		Model(&entity.TestScenario{}).
		Where("id = ?", id).
		Where(`"test_scenarios"."deleted_at" IS NULL`).
		Update("deployment_number", newDeploymentNumber)

	if result.Error != nil {
		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", result.Error.Error()),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToUpdateTestScenario, result.Error)
	}

	if result.RowsAffected == 0 {
		span.SetAttributes(attribute.String("error.type", "not_found"))

		return pkg.ErrTestScenarioNotFound
	}

	return nil
}

func (repo *testScenario) Update(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "update_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))
	span.SetAttributes(
		attribute.String("database.operation", "update"),
		attribute.String("test_scenario.id", strconv.FormatUint(scenario.ID, 10)),
	)

	err := postgres.QueryBuilder(ctx, repo.db).
		Omit(
			"TestCategory",
			"TestCategoryID",
			clause.Associations,
		).
		Model(&entity.TestScenario{}).
		Where("id = ?", scenario.ID).
		Updates(map[string]any{
			"name":                         scenario.Name,
			"mother_service_id":            scenario.MotherServiceID,
			"max_test_service_count":       scenario.MaxTestServiceCount,
			"num_steps":                    scenario.NumSteps,
			"increase_agent_number":        scenario.IncreaseAgentNumber,
			"execution_number_multi_agent": scenario.ExecNumMultiAgent,
			"updated_at":                   time.Now(),
			"status":                       scenario.Status,
		}).Error

	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", err.Error()),
		)

		return fmt.Errorf("failed to update test scenario: %w", err)
	}

	return nil
}

func (repo *testScenario) GetByMotherServiceId(ctx context.Context, motherServiceId uint64) ([]*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "get_by_mother_service_id_test_scenarios")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("mother_service_id", strconv.FormatUint(motherServiceId, 10)))

	var testScenarios []*entity.TestScenario

	query := postgres.QueryBuilder(ctx, repo.db).
		Where("mother_service_id = ?", motherServiceId).
		Preload("TestCategory").
		Preload("MotherService").
		Preload("TestServiceConfig").
		Order("id DESC")

	err := query.Find(&testScenarios).Error

	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	return testScenarios, nil
}

func (repo *testScenario) UpdateScenarioAndConfig(ctx context.Context, scenario *entity.TestScenario) (e error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "update_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))
	span.SetAttributes(
		attribute.String("database.operation", "update"),
		attribute.String("test_scenario.id", strconv.FormatUint(scenario.ID, 10)),
	)

	tx := repo.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)

	defer func() {
		if e != nil {
			_ = tx.Rollback()
			span.SetAttributes(attribute.String("transaction.status", "rolled_back"))
			zap.L().Error("test scenario update failed, transaction rolled back",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", scenario.Name),
				zap.Error(e),
			)
		}
	}()

	err := postgres.QueryBuilder(dbCtx, tx).
		Omit(
			"TestCategory",
			"TestCategoryID",
			clause.Associations,
		).
		Model(&entity.TestScenario{}).
		Where("id = ?", scenario.ID).
		Updates(map[string]any{
			"name":                         scenario.Name,
			"mother_service_id":            scenario.MotherServiceID,
			"max_test_service_count":       scenario.MaxTestServiceCount,
			"num_steps":                    scenario.NumSteps,
			"increase_agent_number":        scenario.IncreaseAgentNumber,
			"execution_number_multi_agent": scenario.ExecNumMultiAgent,
			"updated_at":                   time.Now(),
			"status":                       scenario.Status,
		}).Error

	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", err.Error()),
		)

		return fmt.Errorf("failed to update test scenario: %w", err)
	}

	err = postgres.QueryBuilder(dbCtx, tx).
		Model(&entity.TestServiceConfig{}).
		Where("test_scenario_id = ?", scenario.ID).
		Updates(map[string]any{
			"max_requests":                       scenario.TestServiceConfig.MaxRequests,
			"max_duration":                       scenario.TestServiceConfig.MaxDuration,
			"request_delay_duration":             scenario.TestServiceConfig.RequestDelayDuration,
			"random_request_delay_min":           scenario.TestServiceConfig.RandomRequestDelayMin,
			"random_request_delay_max":           scenario.TestServiceConfig.RandomRequestDelayMax,
			"fixed_test_number":                  scenario.TestServiceConfig.FixedTestNumber,
			"random_test_number_min":             scenario.TestServiceConfig.RandomTestNumberMin,
			"random_test_number_max":             scenario.TestServiceConfig.RandomTestNumberMax,
			"bad_value_rate":                     scenario.TestServiceConfig.BadValueRate,
			"negative_value_rate":                scenario.TestServiceConfig.NegativeValueRate,
			"real_value_rate":                    scenario.TestServiceConfig.RealValueRate,
			"zero_value_rate":                    scenario.TestServiceConfig.ZeroValueRate,
			"string_value_rate":                  scenario.TestServiceConfig.StringValueRate,
			"long_string_value_rate":             scenario.TestServiceConfig.LongStringValueRate,
			"null_value_rate":                    scenario.TestServiceConfig.NullValueRate,
			"database_name":                      scenario.TestServiceConfig.DatabaseName,
			"database_table_name":                scenario.TestServiceConfig.DatabaseTableName,
			"updated_at":                         time.Now(),
			"increase_fixed_input":               scenario.TestServiceConfig.IncreaseFixedInput,
			"execution_number_multi_fixed_input": scenario.TestServiceConfig.ExecNumMultiFixedInput,
		}).Error

	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", err.Error()),
		)

		return fmt.Errorf("failed to update test service config: %w", err)
	}

	_ = tx.Commit()
	span.SetAttributes(attribute.String("transaction.status", "committed"))

	return nil
}

func (repo *testScenario) CreateScenarioAndConfig(ctx context.Context, scenario *entity.TestScenario) (e error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "create_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "insert"), attribute.String("test_scenario.name", scenario.Name))

	tx := repo.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
			span.SetAttributes(attribute.String("transaction.status", "rolled_back"))
			zap.L().Error("test scenario creation failed, transaction rolled back",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", scenario.Name),
				zap.Error(e),
			)
		}
	}()

	err := postgres.QueryBuilder(dbCtx, tx).
		Omit(clause.Associations).
		Create(scenario).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error_create_scenario"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("failed to create test scenario record: %w", err)
	}

	scenario.TestServiceConfig.TestScenarioID = scenario.ID

	err = postgres.QueryBuilder(dbCtx, tx).Create(scenario.TestServiceConfig).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error_create_scenario_config"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("failed to create test service config record: %w", err)
	}

	_ = tx.Commit()

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(scenario.ID, 10)))

	return nil
}
