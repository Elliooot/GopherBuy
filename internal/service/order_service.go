package service

import (
	"GopherBuy/api"
	"GopherBuy/internal/model"
	"GopherBuy/internal/repository"
	"GopherBuy/pkg/utils"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	// dependency injection
	productRepo   *repository.ProductRepository
	orderRepo     *repository.OrderRepository
	flashsaleRepo *repository.FlashSaleRepository
	redisRepo     *repository.RedisRepository
	kafkaRepo     *repository.KafkaRepository
}

// gRPC methods implementation
func (s *OrderService) CreateOrder(req *api.OrderRequest) (*api.OrderResponse, error) {
	product, err := s.productRepo.GetById(req.ProductId)
	if err != nil {
		return nil, err
	}

	if req.Quantity > product.MaxPurchase {
		return &api.OrderResponse{
			Status: 400,
			Msg:    "Exceeded purchase limit",
		}, nil
	}

	order := &model.Order{
		OrderSN:   utils.GenerateOrderSN(req.UserId),
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Amount:    float32(req.Quantity) * product.Price,
		Status:    1,
		CreatedAt: time.Now(),
	}

	// err := s.orderRepo.Create(order)
	// if err != nil {}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return &api.OrderResponse{
		Status:  201,
		Msg:     "Order Created Successfully!",
		OrderSn: order.OrderSN,
	}, nil
}

func (s *OrderService) HandleFlashSaleRequest(req *api.FlashOrderRequest) (*api.OrderResponse, error) {
	flashsale, err := s.redisRepo.GetFlashSale(req.PromoId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if now.Before(flashsale.StartTime) || now.After(flashsale.EndTime) {
		return &api.OrderResponse{
			Status: 400,
			Msg:    "Not in the flashsale duration",
		}, nil
	}

	if req.Quantity > flashsale.MaxPurchase {
		return &api.OrderResponse{
			Status: 400,
			Msg:    "Exceeded purchase limit",
		}, nil
	}

	// Redis
	if err := s.redisRepo.DeductStock(req.PromoId, req.Quantity); err != nil {
		return nil, err
	}

	// Kafka Producer Publish Msg
	event := repository.OrderCreatedEvent{
		OrderSN:    utils.GenerateOrderSN(req.UserId),
		UserID:     req.UserId,
		ProductID:  req.ProductId,
		PromoID:    req.PromoId,
		PromoPrice: req.PromoPrice,
		Quantity:   req.Quantity,
		CreatedAt:  time.Now(),
	}

	if err := s.kafkaRepo.PublishOrderCreated(event); err != nil {
		// if error happened at Kafka, Redis need to rollback stock
		s.redisRepo.RollBackStock(req.PromoId, req.Quantity)
		return nil, err
	}

	return &api.OrderResponse{Status: 202, OrderSn: event.OrderSN, Msg: "Order In Queue"}, nil
}

func (s *OrderService) CreateFlashOrder(req *api.FlashOrderRequest) (*api.OrderResponse, error) {
	var order *model.Order // Define a nil pointer
	// var order = &model.Order{} // Define an empty structure, a bit more waste memory

	err := s.orderRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create Order
		order = &model.Order{
			OrderSN:   utils.GenerateOrderSN(req.UserId),
			UserID:    req.UserId,
			ProductID: req.ProductId,
			Quantity:  req.Quantity,
			Amount:    float32(req.Quantity) * req.PromoPrice,
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiredAt: time.Now().Add(15 * time.Minute),
		}

		// Handle the situation that the order has been created (Error at worker pool or kafka server)
		if err := tx.Create(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				log.Printf("This message has already been handled")
				return nil
			}
			return err
		}

		// Deduct Stock
		if err := s.flashsaleRepo.DeductStock(tx, req.PromoId, req.Quantity); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &api.OrderResponse{
		Status:  201,
		Msg:     "Order Created Successfully!",
		OrderSn: order.OrderSN,
	}, nil
}
