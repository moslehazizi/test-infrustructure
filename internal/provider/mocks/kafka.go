package mocks

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type KafkaMock struct {
	mock.Mock
}

func (kafkaMock *KafkaMock) SendEvent(ctx context.Context, payload []byte, topic string) error {
	args := kafkaMock.Called(ctx, payload, topic)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (kafkaMock *KafkaMock) Consume(ctx context.Context, topic string, ch chan []byte) error {
	args := kafkaMock.Called(ctx, topic, ch)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (kafkaMock *KafkaMock) Close() error {
	args := kafkaMock.Called()

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}
