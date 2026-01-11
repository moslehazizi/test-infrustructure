package provider

import (
	"context"
	"control-panel-service/config"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type EventProducer interface {
	SendEvent(ctx context.Context, payload []byte, topic string) error
	Close() error
}

type kafkaEventProducer struct {
	writer *kafka.Writer
}

func NewKafkaEventProducer(ctx context.Context, cfg *config.Config) (EventProducer, error) {
	log.Printf("initializing Kafka event producer: %s:%d", cfg.Kafka.Host, cfg.Kafka.Port)

	var transport *kafka.Transport
	if cfg.Kafka.Username != "" && cfg.Kafka.Password != "" {
		log.Println("configuring SASL authentication for Kafka producer")
		mechanism := plain.Mechanism{
			Username: cfg.Kafka.Username,
			Password: cfg.Kafka.Password,
		}

		dialer := &kafka.Dialer{
			Timeout:       cfg.Kafka.DialerTimeout,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
		log.Printf("connecting to Kafka broker: %s", brokerAddr)
		_, err := dialer.DialContext(ctx, "tcp", brokerAddr)
		if err != nil {
			log.Printf("failed to connect to Kafka broker %s: %v", brokerAddr, err)

			return nil, fmt.Errorf("could not connect to kafka: %w", err)
		}
		log.Printf("successfully connected to Kafka broker: %s", brokerAddr)

		transport = &kafka.Transport{
			SASL: mechanism,
		}
	} else {
		log.Println("no authentication configured for Kafka producer")
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

	log.Println("Kafka event producer initialized successfully")

	return &kafkaEventProducer{
		writer: kafkaWriter,
	}, nil
}

func (eventProducer *kafkaEventProducer) SendEvent(ctx context.Context, payload []byte, topic string) error {
	log.Printf("sending event to topic '%s', payload size: %d bytes", topic, len(payload))

	err := eventProducer.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Value: payload,
		},
	)

	if err != nil {
		log.Printf("failed to send event to topic '%s': %v", topic, err)

		return fmt.Errorf("failed to write messages: %w", err)
	}

	log.Printf("successfully sent event to topic '%s'", topic)

	return nil
}

func (eventProducer *kafkaEventProducer) Close() error {
	log.Println("closing Kafka producer")
	err := eventProducer.writer.Close()
	if err != nil {
		log.Printf("error closing Kafka producer: %v", err)

		return fmt.Errorf("%w", err)
	}
	log.Println("Kafka producer closed successfully")

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
	log.Printf("initializing Kafka event consumer: %s:%d, topic: %s", cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.ProvisioningTopic)

	var dialer *kafka.Dialer
	if cfg.Kafka.Username != "" && cfg.Kafka.Password != "" {
		log.Println("configuring SASL authentication for Kafka consumer")
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
		log.Printf("connecting to Kafka broker: %s", brokerAddr)
		_, err := dialer.DialContext(ctx, "tcp", brokerAddr)
		if err != nil {
			log.Printf("failed to connect to Kafka broker %s: %v", brokerAddr, err)

			return nil, fmt.Errorf("could not connect to kafka: %w", err)
		}
		log.Printf("successfully connected to Kafka broker: %s", brokerAddr)
	} else {
		log.Println("no authentication configured for Kafka consumer")
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
	log.Printf("Kafka event consumer initialized successfully for topic: %s", cfg.Kafka.ProvisioningTopic)

	return &kafkaConsumer{
		reader: kafkaReader,
	}, nil
}

//nolint:unused // Consume method is kept for future use
func (eventConsumer *kafkaConsumer) Consume(ctx context.Context, topic string, ch chan []byte) error {
	log.Printf("starting to consume messages from topic: %s", topic)
	for {
		select {
		case <-ctx.Done():
			log.Printf("consumption context cancelled for topic: %s", topic)

			return nil
		default:
			m, err := eventConsumer.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("error reading message from kafka topic '%s': %v", topic, err)

				continue
			}
			log.Printf("received message from topic '%s', partition: %d, offset: %d, size: %d bytes",
				m.Topic, m.Partition, m.Offset, len(m.Value))
			ch <- m.Value
		}
	}
}

func (eventConsumer *kafkaConsumer) Close() error {
	log.Println("closing Kafka consumer")
	err := eventConsumer.reader.Close()
	if err != nil {
		log.Printf("error closing Kafka consumer: %v", err)

		return fmt.Errorf("%w", err)
	}
	log.Println("Kafka consumer closed successfully")

	return nil
}
