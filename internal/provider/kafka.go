package provider

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/pkg/logger"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.uber.org/zap"
)

const (
	readMessageTimeout             = 15 * time.Second
	resetMessageTimeout            = 5 * time.Second
	consecutiveTimeoutsToReconnect = 3
)

type EventProducer interface {
	SendEvent(ctx context.Context, payload []byte, topic string) error
	Close() error
}

type kafkaEventProducer struct {
	writer *kafka.Writer
}

func NewKafkaEventProducer(ctx context.Context, cfg *config.Config) (EventProducer, error) {
	zap.L().Info("initializing Kafka event producer",
		zap.String("host", cfg.Kafka.Host),
		zap.Int("port", cfg.Kafka.Port),
	)

	var transport *kafka.Transport
	var dialer *kafka.Dialer
	if cfg.Kafka.Username != "" && cfg.Kafka.Password != "" {
		zap.L().Info("configuring SASL authentication for Kafka producer")
		mechanism := plain.Mechanism{
			Username: cfg.Kafka.Username,
			Password: cfg.Kafka.Password,
		}

		dialer = &kafka.Dialer{
			Timeout:       cfg.Kafka.DialerTimeout,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
		zap.L().Info("connecting to Kafka broker",
			zap.String("broker", brokerAddr),
		)
		_, err := dialer.DialContext(ctx, "tcp", brokerAddr)
		if err != nil {
			zap.L().Error("failed to connect to Kafka broker",
				zap.String("broker", brokerAddr),
				zap.Error(err),
			)

			return nil, fmt.Errorf("could not connect to kafka: %w", err)
		}
		zap.L().Info("successfully connected to Kafka broker",
			zap.String("broker", brokerAddr),
		)

		transport = &kafka.Transport{
			SASL: mechanism,
		}
	} else {
		//nolint
		dialer = &kafka.Dialer{
			Timeout:   cfg.Kafka.DialerTimeout,
			DualStack: true,
		}
		zap.L().Info("no authentication configured for Kafka producer")
	}

	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Balancer:     &kafka.Hash{},
		BatchTimeout: cfg.Kafka.BatchTimeout,
		BatchSize:    cfg.Kafka.BatchSize,
		BatchBytes:   int64(cfg.Kafka.BatchBytes),
	}

	if transport != nil {
		kafkaWriter.Transport = transport
	}

	zap.L().Info("Kafka event producer initialized successfully")

	return &kafkaEventProducer{
		writer: kafkaWriter,
	}, nil
}

func (eventProducer *kafkaEventProducer) SendEvent(ctx context.Context, payload []byte, topic string) error {
	requestID := logger.GetRequestID(ctx)
	zap.L().Info("sending event to Kafka topic",
		zap.String(logger.FieldTopic, topic),
		zap.Int("payload_size", len(payload)),
		zap.String(logger.FieldRequestID, requestID),
	)

	err := eventProducer.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Value: payload,
		},
	)

	if err != nil {
		zap.L().Error("failed to send event to Kafka topic",
			zap.String(logger.FieldTopic, topic),
			zap.String(logger.FieldRequestID, requestID),
			zap.Error(err),
		)

		return fmt.Errorf("failed to write messages: %w", err)
	}

	zap.L().Info("successfully sent event to Kafka topic",
		zap.String(logger.FieldTopic, topic),
		zap.String(logger.FieldRequestID, requestID),
	)

	return nil
}

func (eventProducer *kafkaEventProducer) Close() error {
	zap.L().Info("closing Kafka producer")
	err := eventProducer.writer.Close()
	if err != nil {
		zap.L().Error("error closing Kafka producer", zap.Error(err))

		return err
	}
	zap.L().Info("Kafka producer closed successfully")

	return nil
}

// -------- EventConsumer ---------.
type EventConsumer interface {
	Consume(ctx context.Context, topic string, ch chan []byte) error
	Close() error
}

//nolint:unused // kafkaConsumer is kept for future use
type kafkaConsumer struct {
	reader *kafka.Reader
	cfg    *config.Config
}

func NewKafkaEventConsumer(ctx context.Context, cfg *config.Config, topic string) (EventConsumer, error) {
	zap.L().Info("initializing Kafka event consumer",
		zap.String("host", cfg.Kafka.Host),
		zap.Int("port", cfg.Kafka.Port),
		zap.String(logger.FieldTopic, topic),
	)

	var dialer *kafka.Dialer
	if cfg.Kafka.Username != "" && cfg.Kafka.Password != "" {
		zap.L().Info("configuring SASL authentication for Kafka consumer")
		mechanism := plain.Mechanism{
			Username: cfg.Kafka.Username,
			Password: cfg.Kafka.Password,
		}

		dialer = &kafka.Dialer{
			Timeout:       cfg.Kafka.DialerTimeout,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
		zap.L().Info("connecting to Kafka broker",
			zap.String("broker", brokerAddr),
		)
		_, err := dialer.DialContext(ctx, "tcp", brokerAddr)
		if err != nil {
			zap.L().Error("failed to connect to Kafka broker",
				zap.String("broker", brokerAddr),
				zap.Error(err),
			)

			return nil, fmt.Errorf("could not connect to kafka: %w", err)
		}
		zap.L().Info("successfully connected to Kafka broker",
			zap.String("broker", brokerAddr),
		)
	} else {
		zap.L().Info("no authentication configured for Kafka consumer")
	}

	brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
	readerCfg := kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    topic,
		MaxBytes: cfg.Kafka.MaxBytes,
		GroupID:  cfg.Kafka.ConsumerGroup,
	}
	if dialer != nil {
		readerCfg.Dialer = dialer
	}
	kafkaReader := kafka.NewReader(readerCfg)
	zap.L().Info("Kafka event consumer initialized successfully",
		zap.String(logger.FieldTopic, topic),
	)

	return &kafkaConsumer{
		reader: kafkaReader,
		cfg:    cfg,
	}, nil
}

func newKafkaReader(ctx context.Context, cfg *config.Config, topic string) (*kafka.Reader, error) {
	var dialer *kafka.Dialer
	if cfg.Kafka.Username != "" && cfg.Kafka.Password != "" {
		mechanism := plain.Mechanism{
			Username: cfg.Kafka.Username,
			Password: cfg.Kafka.Password,
		}
		//nolint
		dialer = &kafka.Dialer{
			Timeout:       cfg.Kafka.DialerTimeout,
			DualStack:     true,
			SASLMechanism: mechanism,
		}
		brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
		if _, err := dialer.DialContext(ctx, "tcp", brokerAddr); err != nil {
			return nil, fmt.Errorf("could not connect to kafka: %w", err)
		}
	} else {
		//nolint
		dialer = &kafka.Dialer{
			Timeout:   cfg.Kafka.DialerTimeout,
			DualStack: true,
		}
	}
	brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
	readerCfg := kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    topic,
		MaxBytes: cfg.Kafka.MaxBytes,
		GroupID:  cfg.Kafka.ConsumerGroup,
	}
	//nolint
	if dialer != nil {
		readerCfg.Dialer = dialer
	}

	return kafka.NewReader(readerCfg), nil
}

//nolint:unused // Consume method is kept for future use
func (eventConsumer *kafkaConsumer) Consume(ctx context.Context, topic string, ch chan []byte) error {
	logger.WithContext(ctx).Info("starting to consume messages",
		zap.String(logger.FieldTopic, topic),
		zap.String(logger.FieldOperation, "consume_messages"),
	)
	consecutiveTimeouts := 0
	for {
		if ctx.Err() != nil {
			logger.WithContext(ctx).Info("consumption context cancelled",
				zap.String(logger.FieldTopic, topic),
			)
			//nolint
			return nil
		}
		readCtx, readCancel := context.WithTimeout(ctx, readMessageTimeout)
		m, err := eventConsumer.reader.ReadMessage(readCtx)
		readCancel()
		//nolint
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				consecutiveTimeouts++
				logger.WithContext(ctx).Debug("no message within timeout, retrying read",
					zap.String(logger.FieldTopic, topic),
					zap.Int("consecutive_timeouts", consecutiveTimeouts),
				)
				if consecutiveTimeouts >= consecutiveTimeoutsToReconnect {
					logger.WithContext(ctx).Info("recreating Kafka reader after consecutive timeouts (simulating restart)",
						zap.String(logger.FieldTopic, topic),
						zap.Int("consecutive_timeouts", consecutiveTimeouts),
					)
					_ = eventConsumer.reader.Close()
					for {
						newReader, newErr := newKafkaReader(ctx, eventConsumer.cfg, topic)
						if newErr == nil {
							eventConsumer.reader = newReader
							consecutiveTimeouts = 0
							break
						}
						logger.WithContext(ctx).Error("failed to recreate Kafka reader, retrying",
							zap.Error(newErr),
							zap.String(logger.FieldTopic, topic),
						)
						if ctx.Err() != nil {
							//nolint
							return nil
						}
						select {
						case <-ctx.Done():
							return nil
						case <-time.After(resetMessageTimeout):
						}
					}
				}
			} else {
				consecutiveTimeouts = 0
				logger.WithContext(ctx).Error("error reading message from kafka",
					zap.Error(err),
					zap.String(logger.FieldTopic, topic),
				)
			}
			continue
		}
		consecutiveTimeouts = 0
		logger.WithContext(ctx).Debug("received message from kafka",
			zap.String(logger.FieldTopic, m.Topic),
			zap.Int(logger.FieldPartition, m.Partition),
			zap.Int64(logger.FieldOffset, m.Offset),
			zap.Int("message_size", len(m.Value)),
		)
		ch <- m.Value
	}
}

func (eventConsumer *kafkaConsumer) Close() error {
	zap.L().Info("closing Kafka consumer")
	err := eventConsumer.reader.Close()
	if err != nil {
		zap.L().Error("error closing Kafka consumer", zap.Error(err))

		return err
	}
	zap.L().Info("Kafka consumer closed successfully")

	return nil
}
