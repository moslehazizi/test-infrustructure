package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockConsumerUsecase struct {
	mock.Mock
}

func (m *MockConsumerUsecase) Consume(ctx context.Context, topic string, ch chan []byte) error {
	args := m.Called(ctx, topic, ch)
	return args.Error(0)
}
func (m *MockConsumerUsecase) StoreExecutorResult(ctx context.Context, msg []byte) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
