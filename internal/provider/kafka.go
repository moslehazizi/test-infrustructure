package provider

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/pkg/logger"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.uber.org/zap"
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
	dialer := &kafka.Dialer{}
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
}

func NewKafkaEventConsumer(ctx context.Context, cfg *config.Config) (EventConsumer, error) {
	zap.L().Info("initializing Kafka event consumer",
		zap.String("host", cfg.Kafka.Host),
		zap.Int("port", cfg.Kafka.Port),
		zap.String(logger.FieldTopic, cfg.Kafka.ProvisioningTopic),
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
		Topic:    cfg.Kafka.ProvisioningTopic,
		MaxBytes: cfg.Kafka.MaxBytes,
		GroupID:  cfg.Kafka.ConsumerGroup,
	}
	if dialer != nil {
		readerCfg.Dialer = dialer
	}
	kafkaReader := kafka.NewReader(readerCfg)
	zap.L().Info("Kafka event consumer initialized successfully",
		zap.String(logger.FieldTopic, cfg.Kafka.ProvisioningTopic),
	)

	return &kafkaConsumer{
		reader: kafkaReader,
	}, nil
}

//nolint:unused // Consume method is kept for future use
func (eventConsumer *kafkaConsumer) Consume(ctx context.Context, topic string, ch chan []byte) error {
	zap.L().Info("starting to consume messages from Kafka topic",
		zap.String(logger.FieldTopic, topic),
	)
	for {
		select {
		case <-ctx.Done():
			zap.L().Info("consumption context cancelled",
				zap.String(logger.FieldTopic, topic),
			)

			return nil
		default:
			m, err := eventConsumer.reader.ReadMessage(ctx)
			if err != nil {
				zap.L().Error("error reading message from Kafka topic",
					zap.String(logger.FieldTopic, topic),
					zap.Error(err),
				)

				continue
			}
			zap.L().Info("received message from Kafka topic",
				zap.String(logger.FieldTopic, m.Topic),
				zap.Int(logger.FieldPartition, m.Partition),
				zap.Int64(logger.FieldOffset, m.Offset),
				zap.Int("message_size", len(m.Value)),
			)
			ch <- m.Value
		}
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
