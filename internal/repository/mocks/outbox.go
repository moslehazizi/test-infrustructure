package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockOutbox struct {
	mock.Mock
}

func (m *MockOutbox) Create(ctx context.Context, item *entity.Outbox) (uint64, error) {
	args := m.Called(ctx, item)

	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockOutbox) ClaimPending(ctx context.Context, limit int) ([]*entity.Outbox, error) {
	args := m.Called(ctx, limit)

	var result []*entity.Outbox
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.Outbox)
	}

	return result, args.Error(1)
}

func (m *MockOutbox) MarkCompleted(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockOutbox) Reschedule(ctx context.Context, id uint64, availableAt time.Time, lastErr string) error {
	args := m.Called(ctx, id, availableAt, lastErr)

	return args.Error(0)
}

func (m *MockOutbox) MarkFailedPermanently(ctx context.Context, id uint64, lastErr string) error {
	args := m.Called(ctx, id, lastErr)

	return args.Error(0)
}
