package mocks

import (
	"fmt"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	mock.Mock
}

func (c *Connection) OpenConnection() (*gorm.DB, sqlmock.Sqlmock, error) {
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

	return gormDB, sqlMock, nil
}
