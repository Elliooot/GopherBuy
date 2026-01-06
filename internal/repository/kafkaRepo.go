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
	OrderSN    string    `json:"order_sn"`
	UserID     uint64    `json:"user_id"`
	ProductID  uint64    `json:"product_id"`
	PromoID    uint64    `json:"promo_id"`
	PromoPrice float32   `json:"promo_price"`
	Quantity   uint32    `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
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
