package postgres

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewDatabaseMetadataRepository(db database.Database, cfg *config.Config) repository.DatabaseMetadata {
	return &databaseMetadata{
		db,
		cfg,
	}
}

type databaseMetadata struct {
	db  database.Database
	cfg *config.Config
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
	ctx, span := tracer.Start(ctx, "get_tables_by_db_name")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(
		attribute.String("request_id", requestID),
		attribute.String("database.name", dbName),
		attribute.String("database.operation", "select_tables"),
	)

	db, err := postgres.New(&postgres.DatabaseConfig{
		Host:               repo.cfg.Postgres.Host,
		Port:               repo.cfg.Postgres.Port,
		User:               repo.cfg.Postgres.User,
		Password:           repo.cfg.Postgres.Password,
		Database:           dbName,
		MaxOpenConnections: repo.cfg.Postgres.MaxOpenConnections,
		LogLevel:           postgres.Silent,
	})
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "connection_error"))

		return nil, fmt.Errorf("failed to connect to database %s: %w", dbName, err)
	}
	defer db.Close()

	var tables []string

	query := `
		SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname = 'public'
		ORDER BY tablename;
	`

	err = postgres.QueryBuilder(ctx, db).
		Raw(query).
		Scan(&tables).Error

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "query_error"))

		return nil, fmt.Errorf("failed to load tables for database %s: %w", dbName, err)
	}

	if tables == nil {
		tables = make([]string, 0)
	}

	return tables, nil
}
