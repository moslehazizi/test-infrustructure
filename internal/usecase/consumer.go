package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	sharedentity "control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type consumer struct {
	factorialRepo repository.FactorialRepository
	executorRepo  repository.ExecutorRepository
	eventConsumer provider.EventConsumer
}

func NewCounsumer(factorialRepo repository.FactorialRepository, executorRepo repository.ExecutorRepository, eventConsumer provider.EventConsumer) interfaces.Consumer {
	return &consumer{
		factorialRepo: factorialRepo,
		executorRepo:  executorRepo,
		eventConsumer: eventConsumer,
	}
}

func (mse *consumer) StoreExecutorResult(ctx context.Context, msg []byte) error {
	logger.WithContext(ctx).Debug("storing executor result",
		zap.Int("message_size", len(msg)),
		zap.String(logger.FieldOperation, "store_executor_result"),
	)

	data := sharedentity.ExecutorEvent{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		logger.WithContext(ctx).Error("failed to unmarshal executor event message",
			zap.Error(err),
			zap.String(logger.FieldOperation, "store_executor_result"),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToUnmarshalEventData, err)
	}

	execute := &sharedentity.Executor{
		Input:           data.Input,
		Output:          data.Output,
		MotherServiceId: data.MotherServiceId,
		TestServiceId:   data.TestServiceId,
		StepNum:         data.StepNum,
		ExecutionId:     data.ExecutionId,
		ScenarioId:      data.ScenarioId,
		StepIncrement:   data.StepIncrement,
		StartTxTime:     data.StartTxTime,
		DurationTx:      data.DurationTx,
		DelayBeforeTx:   data.DelayBeforeTx,
		HttpStatusCode:  data.HttpStatusCode,
	}

	err = mse.executorRepo.Create(ctx, execute)
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

func (service *consumer) StoreFactorialResult(ctx context.Context, msg []byte) error {
	logger.WithContext(ctx).Debug("storing factorial result",
		zap.Int("message_size", len(msg)),
		zap.String(logger.FieldOperation, "store_factorial_result"),
	)

	data := entity.FactorialEvent{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		logger.WithContext(ctx).Error("failed to unmarshal factorial event message",
			zap.Error(err),
			zap.String(logger.FieldOperation, "store_factorial_result"),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToUnmarshalEventData, err)
	}

	factorial := &entity.Factorial{
		Input:  data.Input.String(),
		Output: data.Output.String(),
	}

	err = service.factorialRepo.Create(ctx, factorial)
	if err != nil {
		logger.WithContext(ctx).Error("failed to store factorial result",
			zap.Error(err),
			zap.String("input", factorial.Input),
			zap.String(logger.FieldOperation, "store_factorial_result"),
		)

		return fmt.Errorf("failed to store factorial result: %w", err)
	}

	return nil
}

func (mse *consumer) Consume(ctx context.Context, topic string, ch chan []byte) error {
	logger.WithContext(ctx).Info("starting to consume messages",
		zap.String(logger.FieldTopic, topic),
		zap.String(logger.FieldOperation, "consume_messages"),
	)
	err := mse.eventConsumer.Consume(ctx, topic, ch)
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
