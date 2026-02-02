package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewTestScenarioRepository(db database.Database) repository.TestScenarioRepository {
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

	span.SetAttributes(attribute.String("database.operation", "insert"), attribute.String("test_scenario.name", testSci.Name))

	err := postgres.QueryBuilder(ctx, repo.db).
		Omit(clause.Associations).
		Create(testSci).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))
		return 0, fmt.Errorf("failed to create test scenario record: %w", err)
	}

	id := testSci.ID

	span.SetAttributes(attribute.String("test_scenario.id", fmt.Sprintf("%d", id)))

	return id, nil
}

func (repo *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "get_test_scenario_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("test_scenario.id", fmt.Sprintf("%d", id)))

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

	span.SetAttributes(attribute.String("test_scenario.id", fmt.Sprintf("%d", testScenario.ID)))
	return &testScenario, nil
}

func (repo *testScenario) GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "get_paginated_test_scenarios")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("pagination.page", fmt.Sprintf("%d", pagRequest.Page)), attribute.String("pagination.per_page", fmt.Sprintf("%d", pagRequest.PerPage)))

	var testScenarios []*entity.TestScenario
	if pagRequest.Page < 0 || pagRequest.PerPage < 0 {
		span.SetAttributes(attribute.String("error.type", "invalid_pagination"))
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenarios, pkg.ErrNegativePageOrPerPageNotAllowed)
	}

	query := postgres.QueryBuilder(ctx, repo.db).
		Preload("TestCategory").
		Preload("MotherService").
		Order("id DESC")

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
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	return testScenarios, nil
}

func (repo *testScenario) SetStatus(ctx context.Context, id uint64, status entity.ScenarioStatus) error {
	tracer := otel.Tracer("test-scenario-repository")
	_, span := tracer.Start(ctx, "set_test_scenario_status")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "update"), attribute.String("test_scenario.id", fmt.Sprintf("%d", id)), attribute.String("test_scenario.status", fmt.Sprintf("%s", status)))

	err := postgres.QueryBuilder(ctx, repo.db).
		Omit(clause.Associations).
		Model(&entity.TestScenario{}).
		Where("id", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("failed to update test scenario status: %w", err)
	}

	return nil
}
