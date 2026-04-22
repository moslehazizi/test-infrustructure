package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/logger"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type consumer struct {
	factorialRepo    repository.MotherServiceFactorialResultRepository
	executorRepo     repository.TestServiceExecutorResultRepository
	testScenarioRepo repository.TestScenarioRepository
	eventConsumer    provider.EventConsumer
	dbInitializer    database.DBInitializerFn
}

type Consumer interface {
	Consume(ctx context.Context, topic string, ch chan []byte) error
	StoreExecutorResult(ctx context.Context, msg []byte) error
}

func NewConsumer(
	factorialRepo repository.MotherServiceFactorialResultRepository,
	executorRepo repository.TestServiceExecutorResultRepository,
	testScenarioRepo repository.TestScenarioRepository,
	eventConsumer provider.EventConsumer,
	dbInitializer database.DBInitializerFn,

) Consumer {
	return &consumer{
		factorialRepo:    factorialRepo,
		executorRepo:     executorRepo,
		testScenarioRepo: testScenarioRepo,
		eventConsumer:    eventConsumer,
		dbInitializer:    dbInitializer,
	}
}

func (c *consumer) StoreExecutorResult(ctx context.Context, msg []byte) error {
	logger.WithContext(ctx).Debug("storing executor result",
		zap.Int("message_size", len(msg)),
		zap.String(logger.FieldOperation, "store_executor_result"),
	)

	data := entity.ExecutorEvent{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		logger.WithContext(ctx).Error("failed to unmarshal executor event message",
			zap.Error(err),
			zap.String(logger.FieldOperation, "store_executor_result"),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToUnmarshalEventData, err)
	}

	// load test scenario by testScenarioID
	scenario, err := c.testScenarioRepo.GetByID(ctx, data.ScenarioId)
	if err != nil {
		logger.WithContext(ctx).Error("failed to get scenario by id",
			zap.Error(err),
			zap.String(logger.FieldOperation, "store_executor_result"),
		)

		return pkg.ErrFailedToGetTestScenario
	}

	// load scenario by id from database
	execute := &entity.Executor{}
	execute.FromExecutorEvent(&data, scenario)

	err = c.executorRepo.Create(ctx, execute, c.dbInitializer)
	if err != nil {
		logger.WithContext(ctx).Error("failed to store execute result",
			zap.Error(err),
			zap.Any("input", execute.Input),
			zap.Any("MotherServiceId", execute.MotherServiceId),
			zap.Any("Duration", execute.DurationTx),
			zap.String(logger.FieldOperation, "store_execute_result"),
		)

		return fmt.Errorf("failed to store executor result: %w", err)
	}

	return nil
}

func (c *consumer) Consume(ctx context.Context, topic string, ch chan []byte) error {
	logger.WithContext(ctx).Info("starting to consume messages",
		zap.String(logger.FieldTopic, topic),
		zap.String(logger.FieldOperation, "consume_messages"),
	)
	err := c.eventConsumer.Consume(ctx, topic, ch)
	if err != nil {
		logger.WithContext(ctx).Error("failed to consume messages",
			zap.Error(err),
			zap.String(logger.FieldTopic, topic),
			zap.String(logger.FieldOperation, "consume_messages"),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToConsumeData, err)
	}

	logger.WithContext(ctx).Info("stopped consuming messages",
		zap.String(logger.FieldTopic, topic),
		zap.String(logger.FieldOperation, "consume_messages"),
	)

	return nil
}
