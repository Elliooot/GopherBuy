package worker

import (
	"GopherBuy/api"
	"GopherBuy/internal/repository"
	"GopherBuy/internal/service"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type OrderConsumer struct {
	consumer     *kafka.Reader
	orderService *service.OrderService
}

func NewOrderConsumer(
	consumer *kafka.Reader,
	orderService *service.OrderService,
) *OrderConsumer {
	return &OrderConsumer{
		consumer:     consumer,
		orderService: orderService,
	}
}

func (c *OrderConsumer) StartOrderConsumer(workerCount int) {
	ctx := context.Background()
	jobs := make(chan kafka.Message, 100)

	for w := range workerCount {
		go func(workerID int) {
			for job := range jobs {
				var event repository.OrderCreatedEvent
				if err := json.Unmarshal(job.Value, &event); err != nil {
					log.Printf("Worker %d: Failed to unmarshal message: %v\n", workerID, err)
					continue
				}

				resp, err := c.orderService.CreateFlashOrder(&api.FlashOrderRequest{
					UserId:     event.UserID,
					ProductId:  event.ProductID,
					PromoId:    event.PromoID,
					Quantity:   event.Quantity,
					PromoPrice: event.PromoPrice,
				})

				if err != nil {
					log.Printf("Worker %d: Failed to create order %s: %v\n", workerID, event.OrderSN, err)
					continue
				}

				if err := c.consumer.CommitMessages(ctx, job); err != nil {
					log.Printf("Worker %d: Failed to commit offset: %v\n", workerID, err)
				}

				log.Printf("Worker %d: Successfully processed order %s, response: %s\n", workerID, resp.OrderSn, resp.Msg)
			}
		}(w)
	}

	for {
		msg, err := c.consumer.FetchMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %v\n", err)
			continue
		}
		jobs <- msg
	}
}
