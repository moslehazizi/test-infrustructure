package usecase

import (
	"context"
	sharedentity "control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/mocks"
	repoMock "control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	dbmock "control-panel-service/pkg/database/postgres/mocks"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewNewConsumer(t *testing.T) {
	factorialRepo := new(repoMock.MotherServiceFactorialResultRepository)
	executorRepo := new(repoMock.MockMotherServiceExecutorResultRepository)
	testScenarioRepo := new(repoMock.MockTestScenario)
	kafkaConsumer := new(mocks.KafkaMock)

	service := NewConsumer(
		factorialRepo,
		executorRepo,
		testScenarioRepo,
		kafkaConsumer,
		func(cfg any) (database.Database, error) { return nil, nil },
	)
	s, ok := service.(*consumer)

	assert.NotNil(t, service)
	assert.NotNil(t, s.factorialRepo)
	assert.NotNil(t, s.executorRepo)
	assert.NotNil(t, s.testScenarioRepo)

	assert.NotNil(t, s.eventConsumer)
	assert.True(t, ok)
}

func mockDBInitializer() database.DBInitializerFn {
	conn := new(dbmock.Connection)
	db, _, err := conn.OpenConnection()
	if err != nil {
		panic(err)
	}

	return func(cfg any) (database.Database, error) {
		return db, nil
	}
}

func TestStoreExecuteResult(t *testing.T) {
	t.Run("success_stores_execute_event_to_database", func(t *testing.T) {
		factorialRepo := new(repoMock.MotherServiceFactorialResultRepository)
		executorRepo := new(repoMock.MockMotherServiceExecutorResultRepository)
		testScenarioRepo := new(repoMock.MockTestScenario)

		kafkaConsumer := new(mocks.KafkaMock)
		ctx := context.Background()

		service := NewConsumer(factorialRepo, executorRepo, testScenarioRepo, kafkaConsumer, mockDBInitializer())

		eventData := sharedentity.ExecutorEvent{
			Input:           new(string("5")),
			Output:          "120",
			MotherServiceId: "1",
			TestServiceId:   "2",
			DurationTx:      time.Duration(10),
			StartTxTime:     time.Now().Unix(),
			HttpStatusCode:  200,
			ScenarioId:      5,
		}

		scenario := &sharedentity.TestScenario{
			ID: uint64(eventData.ScenarioId),
		}

		expectedExecutor := &sharedentity.Executor{
			Input:           eventData.Input,
			Output:          eventData.Output,
			MotherServiceId: eventData.MotherServiceId,
			TestServiceId:   eventData.TestServiceId,
			DurationTx:      eventData.DurationTx,
			StartTxTime:     eventData.StartTxTime,
			HttpStatusCode:  eventData.HttpStatusCode,
			ScenarioId:      eventData.ScenarioId,
			Scenario:        scenario,
		}

		testScenarioRepo.On("GetByID", mock.Anything, uint64(eventData.ScenarioId)).Return(scenario, nil)
		executorRepo.On("Create", ctx, expectedExecutor, mock.Anything).Return(nil)

		bts, err := json.Marshal(&eventData)
		assert.NoError(t, err)

		err = service.StoreExecutorResult(ctx, bts)

		assert.NoError(t, err)
		executorRepo.AssertCalled(t, "Create", mock.Anything, expectedExecutor, mock.Anything)
		testScenarioRepo.AssertCalled(t, "GetByID", mock.Anything, uint64(eventData.ScenarioId))
	})

	t.Run("error_saving_database_result", func(t *testing.T) {
		factorialRepo := new(repoMock.MotherServiceFactorialResultRepository)
		executorRepo := new(repoMock.MockMotherServiceExecutorResultRepository)
		kafkaConsumer := new(mocks.KafkaMock)
		testScenarioRepo := new(repoMock.MockTestScenario)
		ctx := context.Background()

		service := NewConsumer(factorialRepo, executorRepo, testScenarioRepo, kafkaConsumer, mockDBInitializer())

		eventData := sharedentity.ExecutorEvent{
			Input:           new(string("5")),
			Output:          "120",
			MotherServiceId: "1",
			TestServiceId:   "2",
			DurationTx:      time.Duration(10),
			StartTxTime:     time.Now().Unix(),
			HttpStatusCode:  200,
			ScenarioId:      5,
		}

		scenario := &sharedentity.TestScenario{
			ID: uint64(eventData.ScenarioId),
		}

		expectedExecutor := &sharedentity.Executor{
			Input:           eventData.Input,
			Output:          eventData.Output,
			MotherServiceId: eventData.MotherServiceId,
			TestServiceId:   eventData.TestServiceId,
			DurationTx:      eventData.DurationTx,
			StartTxTime:     eventData.StartTxTime,
			HttpStatusCode:  eventData.HttpStatusCode,
			ScenarioId:      eventData.ScenarioId,
			Scenario:        scenario,
		}

		testScenarioRepo.On("GetByID", mock.Anything, uint64(eventData.ScenarioId)).Return(scenario, nil)
		executorRepo.On("Create", mock.Anything, expectedExecutor, mock.Anything).Return(errors.New("database error"))

		bts, err := json.Marshal(&eventData)
		assert.NoError(t, err)

		err = service.StoreExecutorResult(ctx, bts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		testScenarioRepo.AssertCalled(t, "GetByID", mock.Anything, uint64(eventData.ScenarioId))

	})

	t.Run("error_on_get_scenario_by_id", func(t *testing.T) {
		factorialRepo := new(repoMock.MotherServiceFactorialResultRepository)
		executorRepo := new(repoMock.MockMotherServiceExecutorResultRepository)
		testScenarioRepo := new(repoMock.MockTestScenario)
		kafkaConsumer := new(mocks.KafkaMock)
		ctx := context.Background()

		service := NewConsumer(factorialRepo, executorRepo, testScenarioRepo, kafkaConsumer, mockDBInitializer())

		eventData := sharedentity.ExecutorEvent{
			Input:           new(string("5")),
			Output:          "120",
			MotherServiceId: "1",
			TestServiceId:   "2",
			DurationTx:      time.Duration(10),
			StartTxTime:     time.Now().Unix(),
			HttpStatusCode:  200,
		}

		bts, err := json.Marshal(&eventData)
		assert.NoError(t, err)

		testScenarioRepo.On("GetByID", mock.Anything, uint64(eventData.ScenarioId)).Return(nil, errors.New("some thing went wrong"))

		err = service.StoreExecutorResult(ctx, bts)
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		testScenarioRepo.AssertCalled(t, "GetByID", mock.Anything, uint64(eventData.ScenarioId))
	})

	t.Run("error_on_unmarshaling_data", func(t *testing.T) {
		factorialRepo := new(repoMock.MotherServiceFactorialResultRepository)
		executorRepo := new(repoMock.MockMotherServiceExecutorResultRepository)
		testScenarioRepo := new(repoMock.MockTestScenario)
		kafkaConsumer := new(mocks.KafkaMock)
		ctx := context.Background()

		service := NewConsumer(factorialRepo, executorRepo, testScenarioRepo, kafkaConsumer, mockDBInitializer())

		err := service.StoreExecutorResult(ctx, []byte("invalid data"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), pkg.ErrFailedToUnmarshalEventData.Error())
	})
}

func TestConsumer(t *testing.T) {
	t.Run("failed_return_error", func(t *testing.T) {
		factorialRepo := new(repoMock.MotherServiceFactorialResultRepository)
		executorRepo := new(repoMock.MockMotherServiceExecutorResultRepository)
		testScenarioRepo := new(repoMock.MockTestScenario)
		kafkaConsumer := new(mocks.KafkaMock)

		service := NewConsumer(factorialRepo, executorRepo, testScenarioRepo, kafkaConsumer, mockDBInitializer())

		topic := "test"
		ch := make(chan []byte)

		kafkaConsumer.On("Consume", mock.Anything, topic, ch).Return(errors.New("something went wrong"))

		err := service.Consume(t.Context(), topic, ch)
		assert.Error(t, err)
		kafkaConsumer.AssertCalled(t, "Consume", mock.Anything, topic, ch)
	})

	t.Run("success", func(t *testing.T) {
		factorialRepo := new(repoMock.MotherServiceFactorialResultRepository)
		executorRepo := new(repoMock.MockMotherServiceExecutorResultRepository)
		testScenarioRepo := new(repoMock.MockTestScenario)
		kafkaConsumer := new(mocks.KafkaMock)

		service := NewConsumer(factorialRepo, executorRepo, testScenarioRepo, kafkaConsumer, mockDBInitializer())

		topic := "test"
		ch := make(chan []byte)

		kafkaConsumer.On("Consume", mock.Anything, topic, ch).Return(nil)

		err := service.Consume(t.Context(), topic, ch)
		assert.NoError(t, err)
		kafkaConsumer.AssertCalled(t, "Consume", mock.Anything, topic, ch)
	})
}
