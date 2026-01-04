package repository

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaConsumer *kafka.Reader
	kafkaProducer *kafka.Writer
	kafkaOnce     sync.Once
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func InitKafka(config *KafkaConfig) error {
	var err error
	kafkaOnce.Do(func() {
		// Initialize Consumer
		kafkaConsumer = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        config.Brokers,
			Topic:          config.Topic,
			GroupID:        "gopherbuy-consumer-group",
			MinBytes:       10e3,
			MaxBytes:       10e6,
			CommitInterval: time.Second,
			StartOffset:    kafka.LastOffset,
		})

		kafkaProducer = &kafka.Writer{
			Addr:         kafka.TCP(config.Brokers...),
			Topic:        config.Topic,
			Balancer:     &kafka.LeastBytes{}, // Strategy of Load Balance
			BatchSize:    100,                 // Batch Send
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireAll, // Await for all followers in the ISR(In-Sync Replicas)
			Compression:  kafka.Snappy,     // Compression Strategy
		}

		fmt.Println("Successfully initialize Kafka!")
	})
	return err
}

func GetKafkaConsumer() *kafka.Reader {
	if kafkaConsumer == nil {
		panic("Kafka consumer not initialized")
	}
	return kafkaConsumer
}

func GetKafkaProducer() *kafka.Writer {
	if kafkaProducer == nil {
		panic("Kafka producer not initialized")
	}
	return kafkaProducer
}

func CloseKafka() error {
	var errs []error

	if kafkaConsumer != nil {
		if err := kafkaConsumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("Failed to close consumer: %w", err))
		}
	}

	if kafkaProducer != nil {
		if err := kafkaProducer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("Failed to close producer: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("kafka close error: %v", errs)
	}

	log.Println("Kafka closed")
	return nil
}
