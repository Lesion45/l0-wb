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
	Close() error
}

// KafkaConsumer is an implementation of the Consumer interface for Kafka.
type KafkaConsumer struct {
	log     *zap.Logger
	reader  *kafka.Reader
	service service.Order
}

// NewKafkaConsumer return a new instance of KafkaConsumer with the given configuration.
func NewKafkaConsumer(log *zap.Logger, services *service.Services, brokers []string, groupID string, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		log: log,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
		service: services.Order,
	}
}

// Listen starts the message consumption loop.
func (k *KafkaConsumer) Listen(ctx context.Context) error {
	const op = "broker.KafkaConsumer.Listen"

	k.log.With(
		zap.String("op", op),
	)

	for {
		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				k.log.Info("Consumer context canceled")

				return nil
			}

			k.log.Error("Failed to receive message from Kafka",
				zap.Error(err),
			)

			continue
		}

		k.log.Info("Message received",
			zap.String("Topic", msg.Topic),
			zap.Int("Partition", msg.Partition),
			zap.Int64("Offset", msg.Offset),
			zap.String("Timestamp", msg.Time.Format(time.RFC3339)),
		)

		var order Order

		err = json.Unmarshal(msg.Value, &order)
		if err != nil {
			k.log.Error("Failed to unmarshal data")

			continue
		}

		if err := validator.New().Struct(order); err != nil {
			validateErr := err.(validator.ValidationErrors)

			k.log.Error("invalid data",
				zap.Error(validateErr),
			)

			continue
		}

		err = k.service.SaveOrder(ctx, order.OrderUID, msg.Value)
		if err != nil {
			k.log.Error("unexpected error",
				zap.Error(err),
			)
		}
	}
}

// Close shuts down the Kafka reader and releases resources.
func (k *KafkaConsumer) Close() error {
	return k.reader.Close()
}
