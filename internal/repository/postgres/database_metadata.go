package postgres

import (
	"context"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewDatabaseMetadataRepository(db database.Database) repository.DatabaseMetadata {
	return &databaseMetadata{
		db,
	}
}

type databaseMetadata struct {
	db database.Database
}

func (repo *databaseMetadata) GetAll(ctx context.Context) ([]string, error) {
	tracer := otel.Tracer("database-metadata-repository")
	_, span := tracer.Start(ctx, "get_all_databases")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"))

	var databases []string

	err := postgres.QueryBuilder(ctx, repo.db).
		Raw(`
			SELECT datname
			FROM pg_database
			WHERE datistemplate = false
			ORDER BY datname
		`).
		Scan(&databases).Error

	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", err.Error()),
		)

		return nil, fmt.Errorf("failed to load databases from db: %w", err)
	}

	if databases == nil {
		databases = make([]string, 0)
	}

	return databases, nil
}

func (repo *databaseMetadata) GetTablesByDBName(ctx context.Context, dbName string) ([]string, error) {
	tracer := otel.Tracer("database-metadata-repository")
	_, span := tracer.Start(ctx, "get_tables_by_db_name")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(
		attribute.String("request_id", requestID),
		attribute.String("database.name", dbName),
		attribute.String("database.operation", "select"),
	)

	var tables []string

	err := postgres.QueryBuilder(ctx, repo.db).
		Raw(`
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = 'public'
			ORDER BY table_name
		`).
		Scan(&tables).Error

	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "database_error"),
			attribute.String("error.message", err.Error()),
		)
		return nil, fmt.Errorf("failed to load tables for database %s: %w", dbName, err)
	}

	if tables == nil {
		tables = make([]string, 0)
	}

	return tables, nil
}
