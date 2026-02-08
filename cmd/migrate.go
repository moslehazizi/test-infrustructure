package cmd

import (
	"control-panel-service/config"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"fmt"
	"os"

	"errors"
	"net/url"

	"control-panel-service/migrations"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
				_, _ = os.Stderr.WriteString("failed to load config: " + err.Error() + "\n")

				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize logger
			if err := initLogger(&cfg); err != nil {
				return err
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
				_, _ = os.Stderr.WriteString("failed to load config: " + err.Error() + "\n")

				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize logger
			if err := initLogger(&cfg); err != nil {
				return err
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
				_, _ = os.Stderr.WriteString("failed to load config: " + err.Error() + "\n")

				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize logger
			if err := initLogger(&cfg); err != nil {
				return err
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
				_, _ = os.Stderr.WriteString("failed to load config: " + err.Error() + "\n")

				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize logger
			if err := initLogger(&cfg); err != nil {
				return err
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
	zap.L().Info("checking migration status")
	migrations, err := dbmateDB(cfg).FindMigrations()
	if err != nil {
		zap.L().Error("failed to load migrations", zap.Error(err))

		return fmt.Errorf("failed to load migrations: %w", err)
	}
	for _, m := range migrations {
		if m.Applied {
			zap.L().Info("migration applied",
				zap.String("version", m.Version),
				zap.String("file_path", m.FilePath),
			)
		} else {
			zap.L().Info("migration pending",
				zap.String("version", m.Version),
				zap.String("file_path", m.FilePath),
			)
		}
	}

	return nil
}

func migrate(cfg config.Config) error {
	zap.L().Info("applying migrations")
	err := dbmateDB(cfg).CreateAndMigrate()
	if err != nil {
		zap.L().Error("failed to apply migrations", zap.Error(err))

		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	zap.L().Info("migrations applied successfully")

	return nil
}

func migrateRollback(cfg config.Config) error {
	zap.L().Info("rolling back migration")
	err := dbmateDB(cfg).Rollback()
	if err != nil {
		zap.L().Error("failed to rollback migration", zap.Error(err))

		return fmt.Errorf("failed to rollback migration: %w", err)
	}
	zap.L().Info("migration rolled back successfully")

	return nil
}

func makeMigration(name string, cfg config.Config) error {
	zap.L().Info("creating new migration",
		zap.String("name", name),
	)
	db := dbmateDB(cfg)
	db.MigrationsDir = []string{"migrations"}

	err := db.NewMigration(name)
	if err != nil {
		zap.L().Error("failed to create database migration",
			zap.String("name", name),
			zap.Error(err),
		)

		return fmt.Errorf("failed to create database migration: %w", err)
	}

	zap.L().Info("new migration created successfully",
		zap.String("name", name),
	)

	return nil
}

func initLogger(cfg *config.Config) error {
	loggerConfig := &logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Output: cfg.Logger.Output,
	}
	zapLogger, err := logger.New(loggerConfig, logger.DefaultServiceName)
	if err != nil {
		_, _ = os.Stderr.WriteString("could not initialize logger: " + err.Error() + "\n")

		return fmt.Errorf("could not initialize logger: %w", err)
	}
	zap.ReplaceGlobals(zapLogger)
	defer func() {
		if syncErr := zap.L().Sync(); syncErr != nil {
			_ = syncErr
		}
	}()

	return nil
}

func init() {
	rootCmd.AddCommand(migrateCmd())
}
