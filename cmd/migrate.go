package cmd

import (
	"control-panel-service/config"
	"control-panel-service/pkg/database/postgres"
	"fmt"

	"errors"
	"log"
	"net/url"

	"control-panel-service/migrations"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var ErrMigrationFileNameRequired = errors.New("migration name is required")

func dbmateDB(cfg config.Config) *dbmate.DB {
	connStr := postgres.GetConnectionString(&postgres.DatabaseConfig{
		Host:               cfg.Postgres.Host,
		Port:               cfg.Postgres.Port,
		Database:           cfg.Postgres.Database,
		User:               cfg.Postgres.User,
		SSLMode:            cfg.Postgres.SSLMode,
		Password:           cfg.Postgres.Password,
		MaxOpenConnections: cfg.Postgres.MaxOpenConnections,
		MaxIdleConnections: cfg.Postgres.MaxIdleConnections,
		ConnMaxLifetime:    cfg.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime:    cfg.Postgres.ConnMaxIdleTime,
	})

	u, _ := url.Parse(connStr)
	dbConn := dbmate.New(u)
	dbConn.FS = migrations.Migrations
	dbConn.MigrationsDir = []string{"./"}
	dbConn.AutoDumpSchema = false

	return dbConn
}

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "handle database migration actions",
	}

	migrateMake := &cobra.Command{
		Use:   "make",
		Short: "create new migrate",
		RunE: func(_ *cobra.Command, args []string) error {
			_ = godotenv.Load(envFile)

			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			return makeMigration(args[0], cfg)
		},
	}

	migrateUp := &cobra.Command{
		Use:   "up",
		Short: "migrate the database",
		RunE: func(_ *cobra.Command, _ []string) error {
			_ = godotenv.Load(envFile)

			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			return migrate(cfg)
		},
	}

	migrateDown := &cobra.Command{
		Use:   "down",
		Short: "rollback database migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			_ = godotenv.Load(envFile)

			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			return migrateRollback(cfg)
		},
	}

	migrateStatus := &cobra.Command{
		Use:   "status",
		Short: "get migration status",
		RunE: func(_ *cobra.Command, _ []string) error {
			_ = godotenv.Load(envFile)

			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			return migrateStatus(cfg)
		},
	}

	cmd.AddCommand(migrateMake)
	cmd.AddCommand(migrateUp)
	cmd.AddCommand(migrateDown)
	cmd.AddCommand(migrateStatus)

	return cmd
}

func migrateStatus(cfg config.Config) error {
	log.Println("Migrations:")
	migrations, err := dbmateDB(cfg).FindMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}
	for _, m := range migrations {
		if m.Applied {
			log.Println("[✅]", m.Version, m.FilePath)
		} else {
			log.Println("[❌]", m.Version, m.FilePath)
		}
	}

	return nil
}

func migrate(cfg config.Config) error {
	log.Println("Applying Migrations:")
	err := dbmateDB(cfg).CreateAndMigrate()
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func migrateRollback(cfg config.Config) error {
	err := dbmateDB(cfg).Rollback()
	if err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

func makeMigration(name string, cfg config.Config) error {
	db := dbmateDB(cfg)
	db.MigrationsDir = []string{"migrations"}

	err := db.NewMigration(name)
	if err != nil {
		return fmt.Errorf("failed to create database migration: %w", err)
	}

	log.Println("new migration created: ", name)

	return nil
}

func init() {
	rootCmd.AddCommand(migrateCmd())
}
