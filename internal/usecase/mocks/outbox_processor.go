package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockOutboxProcessor struct {
	mock.Mock
}

func (m *MockOutboxProcessor) ProcessPending(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}
