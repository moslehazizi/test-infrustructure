package postgres

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewDatabaseMetadataRepository(db database.Database, cfg *config.Config) *databaseMetadata {
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
	repoCTX, span := tracer.Start(ctx, "get-all-databases-repository")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "select"))

	var databases []string

	err := postgres.QueryBuilder(repoCTX, repo.db).
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

func (repo *databaseMetadata) GetTablesByDBName(ctx context.Context, dbName string) (*entity.TablesByType, error) {
	tracer := otel.Tracer("database-metadata-repository")
	repoCTX, span := tracer.Start(ctx, "get-tables-by-db-name-repository")
	defer span.End()

	span.SetAttributes(
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

	// First get all table names
	var tables []string
	tableQuery := `
		SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname = 'public'
		ORDER BY tablename;
	`

	err = postgres.QueryBuilder(repoCTX, db).
		Raw(tableQuery).
		Scan(&tables).Error
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "query_error"))
		return nil, fmt.Errorf("failed to load tables for database %s: %w", dbName, err)
	}

	result := &entity.TablesByType{
		MotherTables: make([]string, 0),
		TestTables:   make([]string, 0),
	}

	for _, tableName := range tables {
		columnQuery := `
			SELECT column_name
			FROM information_schema.columns
			WHERE table_schema = 'public' 
			AND table_name = $1
			ORDER BY ordinal_position;
		`

		var columns []string
		err = postgres.QueryBuilder(repoCTX, db).
			Raw(columnQuery, tableName).
			Scan(&columns).Error
		if err != nil {
			continue
		}

		// Since all fields that exist in the mother service table also exist in the test service table,
		// it is critical to first check whether the table is of test type or not. If it is not of test type,
		// then check whether it is of mother type or not.
		if isTestScenarioTable(columns) {
			result.TestTables = append(result.TestTables, tableName)
		} else if isMotherTable(columns) {
			result.MotherTables = append(result.MotherTables, tableName)
		}
	}

	return result, nil
}

func isTestScenarioTable(columns []string) bool {
	requiredColumns := []string{
		"mother_service_id",
		"test_service_id",
		"start_tx_time",
		"step_num",
		"execution_id",
		"scenario_id",
		"duration_tx",
		"delay_before_tx",
		"http_status_code",
	}

	columnSet := make(map[string]bool)
	for _, col := range columns {
		columnSet[col] = true
	}

	for _, required := range requiredColumns {
		if !columnSet[required] {
			return false
		}
	}

	return true
}

func isMotherTable(columns []string) bool {
	requiredColumns := []string{
		"event_id",
		"input",
		"output",
	}

	columnSet := make(map[string]bool)
	for _, col := range columns {
		columnSet[col] = true
	}

	for _, req := range requiredColumns {
		if !columnSet[req] {
			return false
		}
	}

	return true
}
