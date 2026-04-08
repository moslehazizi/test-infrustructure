package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
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
	"gorm.io/gorm"
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

	span.SetAttributes(attribute.String("test_config.id", strconv.FormatUint(id, 10)))

	return nil
}

func (repo *testServiceConfig) GetByID(ctx context.Context, id uint64) (*entity.TestServiceConfig, error) {
	tracer := otel.Tracer("test-service-config-repository")
	_, span := tracer.Start(ctx, "get_test_service_config_by_id")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("service.id", strconv.FormatUint(id, 10)))

	var testSvcConfig entity.TestServiceConfig
	err := postgres.QueryBuilder(ctx, repo.db).First(&testSvcConfig, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return nil, pkg.ErrTestServiceConfigNotFound
		}

		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, fmt.Errorf("failed to get test service config record: %w", err)
	}

	span.SetAttributes(attribute.String("service.id", strconv.FormatUint(testSvcConfig.ID, 10)))

	return &testSvcConfig, nil
}

func (repo *testServiceConfig) UpdateByScenarioID(
	ctx context.Context,
	scenarioID uint64,
	cfg *entity.TestServiceConfig,
) error {
	tracer := otel.Tracer("test-service-config-repository")
	_, span := tracer.Start(ctx, "update_test_service_config")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))
	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(scenarioID, 10)))

	err := postgres.QueryBuilder(ctx, repo.db).
		Model(&entity.TestServiceConfig{}).
		Where("test_scenario_id = ?", scenarioID).
		Updates(map[string]any{
			"max_requests":                       cfg.MaxRequests,
			"max_duration":                       cfg.MaxDuration,
			"request_delay_duration":             cfg.RequestDelayDuration,
			"random_request_delay_min":           cfg.RandomRequestDelayMin,
			"random_request_delay_max":           cfg.RandomRequestDelayMax,
			"fixed_test_number":                  cfg.FixedTestNumber,
			"random_test_number_min":             cfg.RandomTestNumberMin,
			"random_test_number_max":             cfg.RandomTestNumberMax,
			"bad_value_rate":                     cfg.BadValueRate,
			"negative_value_rate":                cfg.NegativeValueRate,
			"real_value_rate":                    cfg.RealValueRate,
			"zero_value_rate":                    cfg.ZeroValueRate,
			"string_value_rate":                  cfg.StringValueRate,
			"long_string_value_rate":             cfg.LongStringValueRate,
			"null_value_rate":                    cfg.NullValueRate,
			"database_name":                      cfg.DatabaseName,
			"database_table_name":                cfg.DatabaseTableName,
			"updated_at":                         time.Now(),
			"increase_fixed_input":               cfg.IncreaseFixedInput,
			"execution_number_multi_fixed_input": cfg.ExecNumMultiFixedInput,
		}).Error

	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", err.Error()),
		)
		return fmt.Errorf("failed to update test service config: %w", err)
	}

	return nil
}
