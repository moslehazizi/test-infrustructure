package mocks

import (
	"control-panel-service/pkg/database"
	pg "control-panel-service/pkg/database/postgres"
	"fmt"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	mock.Mock
}

func (c *Connection) OpenConnection() (database.Database, sqlmock.Sqlmock, error) {
	sqlDB, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("create sqlmock connection: %w", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open gorm mock postgres connection: %w", err)
	}

	dd := &pg.Database{
		DB: gormDB,
	}

	return dd, sqlMock, nil
}
