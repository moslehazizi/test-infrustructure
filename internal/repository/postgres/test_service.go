package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewTestServiceRepository(db database.Database) repository.TestServiceRepository {
	return &testService{
		db: db,
	}
}

type testService struct {
	db database.Database
}

func (repo *testService) GetRunningByScenario(ctx context.Context, scenarioID uint64, limit int64) ([]entity.TestService, error) {
	tracer := otel.Tracer("test-service-repository")
	_, span := tracer.Start(ctx, "GetRunningByScenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	qry := postgres.QueryBuilder(ctx, repo.db)
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
	_, span := tracer.Start(ctx, "GetCountAllRunningByScenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	var cnt int64
	err := postgres.QueryBuilder(ctx, repo.db).
		Model(&entity.TestService{}).
		Where("status", entity.ScenarioStatusRunning).
		Count(&cnt).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return 0, fmt.Errorf("failed to get test services count: %w", err)
	}

	return cnt, nil
}
