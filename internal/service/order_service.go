package service

import (
	"GopherBuy/api"
	"GopherBuy/internal/model"
	"GopherBuy/internal/repository"
	"GopherBuy/pkg/utils"
)

type OrderService struct {
	// dependency injection
	productRepo *repository.ProductRepository
	orderRepo   *repository.OrderRepository
}

// gRPC methods implementation
func (s *OrderService) CreateOrder(req *api.OrderRequest) (*model.Order, error) {
	price, err := s.productRepo.GetPriceById(req.ProductId)
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		OrderSN:   utils.GenerateOrderSN(req.UserId),
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Amount:    float32(req.Quantity) * price,
		Status:    1,
	}

	// err := s.orderRepo.Create(order)
	// if err != nil {}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

// func (s *OrderService) CreateFlashOrder(req *api.FlashOrderRequest) (*api.OrderResponse, error) {

// }
