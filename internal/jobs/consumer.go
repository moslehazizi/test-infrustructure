package jobs

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository/postgres"
	"control-panel-service/internal/usecase"
	psql "control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

type consumerJob struct {
	consumerUseCase usecase.Consumer
}

func Serve(ctx context.Context, cfg *config.Config) error {
	db, err := psql.New(&psql.DatabaseConfig{
		Host:               cfg.Postgres.Host,
		Port:               cfg.Postgres.Port,
		User:               cfg.Postgres.User,
		Password:           cfg.Postgres.Password,
		Database:           cfg.Postgres.Database,
		MaxOpenConnections: cfg.Postgres.MaxOpenConnections,
		LogLevel:           psql.Silent,
	})
	if err != nil {
		return fmt.Errorf("could not open postgres connection: %w", err)
	}
	factRepo := postgres.NewMotherServiceFactorialResultRepository(cfg)
	execRepo := postgres.NewTestServiceExecutorResultRepository(cfg)
	motherServiceRepo := postgres.NewMotherServiceRepository(db)
	testScenarioRepo := postgres.NewTestScenarioRepository(db)

	eventConsumer, err := provider.NewKafkaEventConsumer(ctx, cfg)
	if err != nil {
		return fmt.Errorf("could not create kafka event consumer: %w", err)
	}
	defer func() {
		err := eventConsumer.Close()
		if err != nil {
			zap.L().Error("could not close kafka event consumer",
				zap.Error(err),
				zap.String(logger.FieldOperation, "close_kafka_consumer"),
			)
		}
	}()

	consumer := usecase.NewConsumer(factRepo, execRepo, motherServiceRepo, testScenarioRepo, eventConsumer, psql.DBInitializerFn)
	consumerJob := &consumerJob{consumer}

	consumerFactChannel := make(chan []byte)
	go consumerJob.RunFactorialConsumer(ctx, cfg.Kubernetese.MotherServiceKafkaDatabaseTopic, consumerFactChannel)

	consumerExecChannel := make(chan []byte)
	go consumerJob.RunExecutorConsumer(ctx, cfg.Kubernetese.TestServiceKafkaDatabaseTopic, consumerExecChannel)

	<-ctx.Done()

	return nil
}

func (c *consumerJob) RunFactorialConsumer(ctx context.Context, topic string, ch chan []byte) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := c.consumerUseCase.Consume(ctx, topic, ch)
		if err != nil {
			logger.WithContext(ctx).Error("could not consume factorial result",
				zap.Error(err),
				zap.String(logger.FieldOperation, "consume_factorial_result"),
			)
			cancel()

			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-ch:
			err := c.consumerUseCase.StoreFactorialResult(ctx, data)
			if err != nil {
				logger.WithContext(ctx).Error("could not store factorial result",
					zap.Error(err),
					zap.String(logger.FieldOperation, "store_factorial_result"),
				)

				continue
			}
		}
	}
}

func (c *consumerJob) RunExecutorConsumer(ctx context.Context, topic string, ch chan []byte) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := c.consumerUseCase.Consume(ctx, topic, ch)
		if err != nil {
			logger.WithContext(ctx).Error("could not consume executor result",
				zap.Error(err),
				zap.String(logger.FieldOperation, "consume_executor_result"),
			)
			cancel()

			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-ch:
			err := c.consumerUseCase.StoreExecutorResult(ctx, data)
			if err != nil {
				logger.WithContext(ctx).Error("could not store execute result",
					zap.Error(err),
					zap.String(logger.FieldOperation, "store_execute_result"),
				)

				continue
			}
		}
	}
}
