package broker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
	"wb-internship-l0/internal/service"
)

// Consumer defines an interface for various consumer implementations.
type Consumer interface {
	Listen(ctx context.Context) error
	Shutdown() error
}

// KafkaConsumer is an implementation of the Consumer interface for Kafka.
type KafkaConsumer struct {
	log     *zap.Logger
	reader  *kafka.Reader
	service service.Order
}

// NewKafkaConsumer return a new instance of KafkaConsumer with the given configuration.
func NewKafkaConsumer(log *zap.Logger, services *service.Services, brokers []string, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		log: log,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
		service: services.Order,
	}
}

// Listen starts the message consumption loop.
func (k *KafkaConsumer) Listen(ctx context.Context) error {
	const op = "broker.KafkaConsumer.Listen"

	k.log.Info("Kafka reader is running")

	for {
		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				k.log.Info("Consumer context canceled")

				return err
			}

			continue
		}

		k.log.Info("Message received from Kafka",
			zap.String("Topic", msg.Topic),
			zap.Int("Partition", msg.Partition),
			zap.Int64("Offset", msg.Offset),
			zap.String("Timestamp", msg.Time.Format(time.RFC3339)),
		)

		var order Order

		err = json.Unmarshal(msg.Value, &order)
		if err != nil {
			k.log.Error("Failed to unmarshal data",
				zap.String("op", op),
				zap.Error(err),
			)

			continue
		}

		if err := validator.New().Struct(order); err != nil {
			validateErr := err.(validator.ValidationErrors)

			k.log.Error("Invalid data",
				zap.String("op", op),
				zap.Error(validateErr),
			)

			continue
		}

		err = k.service.SaveOrder(ctx, order.OrderUID, msg.Value)
		if err != nil {
			k.log.Error("Unexpected error",
				zap.String("op", op),
				zap.Error(err),
			)

			continue
		}

		if err := k.reader.CommitMessages(ctx, msg); err != nil {
			k.log.Error("Failed to commit message to Kafka",
				zap.String("op", op),
				zap.Error(err),
			)
		}
	}

	return nil
}

// Shutdown shuts down the Kafka reader and releases resources.
func (k *KafkaConsumer) Shutdown() {
	const op = "broker.KafkaConsumer.Shutdown"

	k.log.Info("Closing Kafka reader...")
	if err := k.reader.Close(); err != nil {
		k.log.Error("Failed to close Kafka reader",
			zap.String("op", op),
			zap.Error(err),
		)
		return
	} else {
		k.log.Info("Kafka reader closed")
	}
	return
}
