package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaRepository struct {
	producer *kafka.Writer
}

func NewKafkaRepository(producer *kafka.Writer, topic string) *KafkaRepository {
	return &KafkaRepository{
		producer: producer,
	}
}

type OrderCreatedEvent struct {
	OrderSN   string    `json:"order_sn"`
	UserID    uint64    `json:"user_id"`
	ProductID uint64    `json:"product_id"`
	Quantity  uint32    `json:"quantity"`
	Amount    uint32    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *KafkaRepository) PublishOrderCreated(event OrderCreatedEvent) error {
	// Encoding event
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = r.producer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(event.OrderSN),
		Value: payload,
		Time:  time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}
