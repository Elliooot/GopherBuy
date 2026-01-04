package main

import (
	"GopherBuy/internal/repository"
	"log"
)

func main() {
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	if err := repository.InitRedis(); err != nil {
		log.Fatalf("Failed to init Redis: %v", err)
	}

	kafkaConfig := &repository.KafkaConfig{
		Brokers: []string{"localhost:9092", "localhost:9094"},
		Topic:   "order_create",
	}

	if err := repository.InitKafka(kafkaConfig); err != nil {
		log.Fatalf("Failed to init Kafka: %v", err)
	}
}
