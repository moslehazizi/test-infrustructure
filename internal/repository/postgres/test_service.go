package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewTestServiceRepository(db database.Database) *testService {
	return &testService{
		db: db,
	}
}

type testService struct {
	db database.Database
}

func (repo *testService) GetRunningByScenario(ctx context.Context, scenarioID uint64, limit int64) ([]entity.TestService, error) {
	tracer := otel.Tracer("test-service-repository")
	repoCTX, span := tracer.Start(ctx, "GetRunningByScenario")
	defer span.End()

	qry := postgres.QueryBuilder(repoCTX, repo.db)
	if limit >= 0 {
		qry = qry.Limit(int(limit))
	}
	var items []entity.TestService
	err := qry.
		Where("status", entity.ScenarioStatusRunning).
		Find(&items).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, fmt.Errorf("failed to get test services records: %w", err)
	}

	return items, nil
}

func (repo *testService) GetCountAllRunningByScenario(ctx context.Context, scenarioID uint64) (int64, error) {
	tracer := otel.Tracer("test-service-repository")
	repoCTX, span := tracer.Start(ctx, "get-count-all-running-by-scenario-repository")
	defer span.End()

	var cnt int64
	err := postgres.QueryBuilder(repoCTX, repo.db).
		Model(&entity.TestService{}).
		Where("status", entity.ScenarioStatusRunning).
		Count(&cnt).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return 0, fmt.Errorf("failed to get test services count: %w", err)
	}

	return cnt, nil
}
