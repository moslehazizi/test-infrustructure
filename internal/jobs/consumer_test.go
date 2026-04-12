package jobs

import (
	"context"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestRunExecutorConsumer(t *testing.T) {
	t.Run("failed_case_RunExecutorConsumer_kafka_error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second*3600)
		defer cancel()

		factorialCha := make(chan []byte)
		factorialTopic := "factorial"

		mockConsumerUscase := new(mocks.MockConsumerUsecase)
		c := &consumerJob{mockConsumerUscase}

		mockConsumerUscase.On("Consume", mock.Anything, "factorial", factorialCha).Return(fmt.Errorf("failed to consume data"))

		c.RunExecutorConsumer(ctx, factorialTopic, factorialCha)

		mockConsumerUscase.AssertCalled(t, "Consume", mock.Anything, "factorial", factorialCha)
	})

	t.Run("failed_case_kafka_error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second*3600)
		defer cancel()

		executorCha := make(chan []byte)
		executorTopic := "executor"

		mockConsumerUscase := new(mocks.MockConsumerUsecase)
		c := &consumerJob{mockConsumerUscase}

		mockConsumerUscase.On("Consume", mock.Anything, "executor", executorCha).Return(fmt.Errorf("failed to consume data"))

		c.RunExecutorConsumer(ctx, executorTopic, executorCha)

		mockConsumerUscase.AssertCalled(t, "Consume", mock.Anything, "executor", executorCha)
	})

	t.Run("failed_case_database_error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
		defer cancel()

		chanel := make(chan []byte)
		topik := "executor"

		mockConsumerUscase := new(mocks.MockConsumerUsecase)
		c := &consumerJob{mockConsumerUscase}

		go func() {
			time.Sleep(100 * time.Millisecond)
			chanel <- []byte("sample")
		}()

		mockConsumerUscase.On("Consume", mock.Anything, "executor", chanel).Return(nil)
		mockConsumerUscase.On("StoreExecutorResult", mock.Anything, []byte("sample")).Return(fmt.Errorf("%w: %w", pkg.ErrFailedToUnmarshalEventData, errors.New("failed to unmarshal data")))

		c.RunExecutorConsumer(ctx, topik, chanel)

		mockConsumerUscase.AssertCalled(t, "StoreExecutorResult", mock.Anything, []byte("sample"))
		mockConsumerUscase.AssertCalled(t, "Consume", mock.Anything, "executor", chanel)
	})

	t.Run("success_case", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
		defer cancel()

		factorialCha := make(chan []byte)
		factorialTopic := "executor"

		mockConsumerUscase := new(mocks.MockConsumerUsecase)
		c := &consumerJob{mockConsumerUscase}

		go func() {
			time.Sleep(100 * time.Millisecond)
			factorialCha <- []byte("sample")
		}()

		mockConsumerUscase.On("Consume", mock.Anything, "executor", factorialCha).Return(nil)
		mockConsumerUscase.On("StoreExecutorResult", mock.Anything, []byte("sample")).Return(nil)

		c.RunExecutorConsumer(ctx, factorialTopic, factorialCha)

		mockConsumerUscase.AssertCalled(t, "StoreExecutorResult", mock.Anything, []byte("sample"))
		mockConsumerUscase.AssertCalled(t, "Consume", mock.Anything, "executor", factorialCha)
	})

}
