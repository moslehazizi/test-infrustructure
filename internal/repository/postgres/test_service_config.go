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

func NewTestServiceConfigRepository(db database.Database) repository.TestServiceConfigRepository {
	return &testServiceConfig{
		db: db,
	}
}

type testServiceConfig struct {
	db database.Database
}

func (repo *testServiceConfig) Create(ctx context.Context, testSvcCfg *entity.TestServiceConfig) error {
	tracer := otel.Tracer("test-service-config-repository")
	_, span := tracer.Start(ctx, "create_test_service_config")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	err := postgres.QueryBuilder(ctx, repo.db).Create(testSvcCfg).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("failed to create test service config record: %w", err)
	}

	id := testSvcCfg.ID

	span.SetAttributes(attribute.String("test_config.id", fmt.Sprintf("%d", id)))

	return nil
}
