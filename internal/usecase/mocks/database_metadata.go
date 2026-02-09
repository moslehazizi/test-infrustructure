package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDatabaseMetadataUsecase struct {
	mock.Mock
}

func (m *MockDatabaseMetadataUsecase) GetAll(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	var result []string
	if args.Get(0) != nil {
		result = args.Get(0).([]string)
	}

	return result, args.Error(1)
}

func (m *MockDatabaseMetadataUsecase) GetTablesByDBName(ctx context.Context, dbName string) ([]string, error) {
	args := m.Called(ctx, dbName)

	var result []string
	if args.Get(0) != nil {
		result = args.Get(0).([]string)
	}

	return result, args.Error(1)
}
