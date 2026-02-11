package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"control-panel-service/pkg/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host               string
	Port               int
	User               string
	Password           string `json:"-"`
	Database           string
	SSLMode            string
	LogLevel           LogLevel
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
	ConnMaxIdleTime    time.Duration
}

// LogLevel log level.
type LogLevel int

const (
	// Silent silent log level.
	Silent LogLevel = iota + 1
	// Error error log level.
	Error
	// Warn warn log level.
	Warn
	// Info info log level.
	Info
)

type Database struct {
	url      string
	cfg      *DatabaseConfig
	DB       *gorm.DB
	logLevel LogLevel
}

// New return new database connection instance.
func New(config *DatabaseConfig) (database.Database, error) {
	d := &Database{
		url:      GetConnectionString(config),
		cfg:      config,
		logLevel: config.LogLevel,
	}

	err := d.Init()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Init initialize database connection.
func (d *Database) Init() error {
	orm, err := gorm.Open(postgres.New(postgres.Config{
		DSN: d.url,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	orm.Logger = orm.Logger.LogMode(logger.LogLevel(d.logLevel))
	sql, err := orm.DB()
	if err != nil {
		return fmt.Errorf("failed to get orm DB instance to initialize database: %w", err)
	}
	// Connection pool settings
	sql.SetMaxOpenConns(d.cfg.MaxOpenConnections)
	sql.SetMaxIdleConns(d.cfg.MaxIdleConnections)
	sql.SetConnMaxLifetime(d.cfg.ConnMaxLifetime)
	sql.SetConnMaxIdleTime(d.cfg.ConnMaxIdleTime)

	d.DB = orm

	return nil
}

// Ping ping database to make sure connection is alive.
func (d *Database) Ping() error {
	sql, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	err = sql.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close close database connection.
func (d *Database) Close() error {
	sql, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to ping database instance to close: %w", err)
	}

	err = sql.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

func (d *Database) Client() database.Database {
	return d
}

// WithContext get new connection instance using given context.
func (d *Database) WithContext(ctx context.Context) database.Database {
	tx := ctx.Value(database.ContextKeyDBTx)
	if tx == nil {
		return &Database{
			url: d.url,
			cfg: d.cfg,
			DB:  d.DB.WithContext(ctx),
		}
	}

	dbTx, ok := tx.(*Database)
	if !ok {
		return &Database{
			url: d.url,
			cfg: d.cfg,
			DB:  d.DB.WithContext(ctx),
		}
	}

	return dbTx
}

// Begin start new database transaction.
func (d *Database) Begin() database.Database {
	return &Database{
		DB: d.DB.Begin(),
	}
}

// Commit commit the database transaction.
func (d *Database) Commit() error {
	return d.DB.Commit().Error
}

// Rollback rollback the ongoing transaction.
func (d *Database) Rollback() error {
	return d.DB.Rollback().Error
}

func QueryBuilder(ctx context.Context, db database.Database) *gorm.DB {
	//nolint
	return db.WithContext(ctx).(*Database).DB
}

// GetConnectionString return postgres connection string from configs.
func GetConnectionString(cfg *DatabaseConfig) string {
	userpass := url.UserPassword(cfg.User, cfg.Password)
	conURL := url.URL{
		Scheme: "postgres",
		User:   userpass,
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:   cfg.Database,
	}

	qs := url.Values{}
	qs.Add("sslmode", cfg.SSLMode)
	conURL.RawQuery = qs.Encode()

	return conURL.String()
}
