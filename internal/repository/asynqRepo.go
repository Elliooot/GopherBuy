package repository

import (
	"GopherBuy/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const TypeOrderExpired = "order:expired"

type AsynqRepository struct {
	client        *asynq.Client
	orderRepo     *OrderRepository
	redisRepo     *RedisRepository
	flashsaleRepo *FlashSaleRepository
}

func NewAsynqRepository(
	client *asynq.Client,
	orderRepo *OrderRepository,
	redisRepo *RedisRepository,
	flashsaleRepo *FlashSaleRepository,
) *AsynqRepository {
	return &AsynqRepository{
		client:        client,
		orderRepo:     orderRepo,
		redisRepo:     redisRepo,
		flashsaleRepo: flashsaleRepo,
	}
}

type OrderExpiredPayload struct {
	OrderSN  string `json:"order_sn"`
	PromoID  uint64 `json:"promo_id"`
	Quantity uint32 `json:"quantity"`
}

func (r *AsynqRepository) EnqueueOrderExpiration(orderSN string, promoId uint64, quantity uint32) error {
	payload := OrderExpiredPayload{
		OrderSN:  orderSN,
		PromoID:  promoId,
		Quantity: quantity,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeOrderExpired, data)

	info, err := r.client.Enqueue(
		task,
		asynq.ProcessIn(15*time.Minute),
		asynq.MaxRetry(3),
		asynq.TaskID(orderSN),
	)

	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Order expriration task enqueued: OrderSN=%s, ProcessAt=%v\n", orderSN, info.NextProcessAt)
	return nil
}

func (r *AsynqRepository) HandleOrderExpiration(ctx context.Context, task *asynq.Task) error {
	var payload OrderExpiredPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	log.Printf("Processing expired order: %s\n", payload.OrderSN)

	err := r.orderRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.Where("order_sn = ?", payload.OrderSN).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Order not found: %s\n", payload.OrderSN)
				return nil
			}
			return err
		}

		if order.Status != 1 {
			log.Printf("Order already processed: %s (status = %d)\n", payload.OrderSN, order.Status)
			return nil
		}
		// Update Expired Order Status
		if err := tx.Model(&order).Update("status", 3).Error; err != nil {
			return fmt.Errorf("failed to cancel expired order: %w", err)
		}
		// Rollback stock in db
		if err := r.flashsaleRepo.RollBackStock(tx, payload.PromoID, payload.Quantity); err != nil {
			return fmt.Errorf("failed to rollback stock in DB: %w", err)
		}

		log.Printf("Order cancelled and DB stock rolled back: %s\n", payload.OrderSN)
		return nil
	})

	if err != nil {
		return err
	}

	// Rollback stock in Redis
	if err := r.redisRepo.RollBackStock(payload.PromoID, payload.Quantity); err != nil {
		log.Printf("Warning: failed to rollback stock in Redis: %v\n", err)
	} else {
		log.Printf("Redis stock rolled back: %s\n", payload.OrderSN)
	}

	return nil
}
