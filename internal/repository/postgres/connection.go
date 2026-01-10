package postgres

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/pkg"
	"fmt"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func OpenConnection(ctx context.Context, cfg config.Postgres) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(pkg.GetConnectionString(
		pkg.DatabaseConfig{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			Password: cfg.Password,
			Database: cfg.Database,
			SSLMode:  cfg.SSLMode,
		},
	)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open gorm postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB from gorm: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnection)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}
