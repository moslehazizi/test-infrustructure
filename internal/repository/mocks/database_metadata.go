package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockDatabaseMetadata struct {
	mock.Mock
}

func (m *MockDatabaseMetadata) GetAll(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	var result []string
	if args.Get(0) != nil {
		result = args.Get(0).([]string)
	}

	return result, args.Error(1)
}

func (m *MockDatabaseMetadata) GetTablesByDBName(ctx context.Context, dbName string) (*entity.TablesByType, error) {
	args := m.Called(ctx, dbName)

	var result *entity.TablesByType
	if args.Get(0) != nil {
		result = args.Get(0).(*entity.TablesByType)
	}

	return result, args.Error(1)
}
