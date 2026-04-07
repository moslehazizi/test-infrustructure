package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	sharedentity "control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/mocks"
	repoMock "control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewNewCounsumer(t *testing.T) {
	factorialRepo := new(repoMock.FactorialRepository)
	executorRepo := new(repoMock.ExecutorRepository)
	kafkaConsumer := new(mocks.KafkaMock)

	service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)
	s, ok := service.(*consumer)

	assert.NotNil(t, service)
	assert.NotNil(t, s.factorialRepo)
	assert.NotNil(t, s.executorRepo)
	assert.NotNil(t, s.eventConsumer)
	assert.True(t, ok)
}

func TestStoreExecuteResult(t *testing.T) {
	t.Run("success_stores_factorial_event_to_database", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)
		ctx := context.Background()

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		input := "5"
		eventData := sharedentity.ExecutorEvent{
			Input:           &input,
			Output:          "120",
			MotherServiceId: "1",
			TestServiceId:   "2",
			DurationTx:      time.Duration(10),
			StartTxTime:     time.Now().Unix(),
			HttpStatusCode:  200,
		}

		expectedExecutor := &sharedentity.Executor{
			Input:           eventData.Input,
			Output:          eventData.Output,
			MotherServiceId: eventData.MotherServiceId,
			TestServiceId:   eventData.TestServiceId,
			DurationTx:      eventData.DurationTx,
			StartTxTime:     eventData.StartTxTime,
			HttpStatusCode:  eventData.HttpStatusCode,
		}

		executorRepo.On("Create", ctx, expectedExecutor).Return(nil)

		bts, _ := json.Marshal(&eventData)
		err := service.StoreExecutorResult(ctx, bts)

		assert.NoError(t, err)
		executorRepo.AssertCalled(t, "Create", mock.Anything, expectedExecutor)
	})

	t.Run("error_on_unmarshalling_data", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)
		ctx := context.Background()

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		err := service.StoreExecutorResult(ctx, []byte("invalid data"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), pkg.ErrFailedToUnmarshalEventData.Error())
	})

	t.Run("error_saving_database_result", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)
		ctx := context.Background()

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		input := "5"
		eventData := sharedentity.ExecutorEvent{
			Input:           &input,
			Output:          "120",
			MotherServiceId: "1",
			TestServiceId:   "2",
			DurationTx:      time.Duration(10),
			StartTxTime:     time.Now().Unix(),
			HttpStatusCode:  200,
		}

		bts, _ := json.Marshal(&eventData)

		expectedExecutor := &sharedentity.Executor{
			Input:           eventData.Input,
			Output:          eventData.Output,
			MotherServiceId: eventData.MotherServiceId,
			TestServiceId:   eventData.TestServiceId,
			DurationTx:      eventData.DurationTx,
			StartTxTime:     eventData.StartTxTime,
			HttpStatusCode:  eventData.HttpStatusCode,
		}

		executorRepo.On("Create", mock.Anything, expectedExecutor).Return(errors.New("database error"))

		err := service.StoreExecutorResult(ctx, bts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestStoreFactorialResult(t *testing.T) {
	t.Run("success - stores factorial event to database", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		eventData := entity.FactorialEvent{
			Input:  big.NewInt(5),
			Output: big.NewInt(120),
		}
		bts, _ := json.Marshal(&eventData)

		expectedFactorial := &entity.Factorial{
			Input:  eventData.Input.String(),
			Output: eventData.Output.String(),
		}

		factorialRepo.On("Create", mock.Anything, mock.MatchedBy(func(factorial *entity.Factorial) bool {
			return factorial.Input == expectedFactorial.Input &&
				factorial.Output == expectedFactorial.Output
		})).Return(nil)

		err := service.StoreFactorialResult(context.Background(), bts)

		assert.NoError(t, err)
		factorialRepo.AssertCalled(t, "Create", mock.Anything, mock.MatchedBy(func(factorial *entity.Factorial) bool {
			return factorial.Input == expectedFactorial.Input &&
				factorial.Output == expectedFactorial.Output
		}))
	})

	t.Run("error on unmarshalling data", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		err := service.StoreFactorialResult(context.Background(), []byte("invalid data"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), pkg.ErrFailedToUnmarshalEventData.Error())
	})

	t.Run("error saving database result", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		eventData := entity.FactorialEvent{
			Input:  big.NewInt(5),
			Output: big.NewInt(120),
		}
		bts, _ := json.Marshal(&eventData)

		expectedFactorial := &entity.Factorial{
			Input:  eventData.Input.String(),
			Output: eventData.Output.String(),
		}

		factorialRepo.On("Create", mock.Anything, mock.MatchedBy(func(factorial *entity.Factorial) bool {
			return factorial.Input == expectedFactorial.Input &&
				factorial.Output == expectedFactorial.Output
		})).Return(errors.New("database error"))

		err := service.StoreFactorialResult(context.Background(), bts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestConsumer(t *testing.T) {
	t.Run("failed_return_error", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		topic := "test"
		ch := make(chan []byte)

		kafkaConsumer.On("Consume", mock.Anything, topic, ch).Return(errors.New("something went wrong"))

		err := service.Consume(t.Context(), topic, ch)
		assert.Error(t, err)
		kafkaConsumer.AssertCalled(t, "Consume", mock.Anything, topic, ch)
	})

	t.Run("success", func(t *testing.T) {
		factorialRepo := new(repoMock.FactorialRepository)
		executorRepo := new(repoMock.ExecutorRepository)
		kafkaConsumer := new(mocks.KafkaMock)

		service := NewCounsumer(factorialRepo, executorRepo, kafkaConsumer)

		topic := "test"
		ch := make(chan []byte)

		kafkaConsumer.On("Consume", mock.Anything, topic, ch).Return(nil)

		err := service.Consume(t.Context(), topic, ch)
		assert.NoError(t, err)
		kafkaConsumer.AssertCalled(t, "Consume", mock.Anything, topic, ch)
	})
}
