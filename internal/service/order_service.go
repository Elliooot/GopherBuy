package service

import (
	"GopherBuy/internal/model"
	"GopherBuy/internal/repository"
	"GopherBuy/pkg/utils"
)

type OrderService struct {
	// dependency injection
	orderRepo *repository.OrderRepository
}

// gRPC methods implementation
func (s *OrderService) CreateOrder(userID uint64, productID uint64, amount float32) (*model.Order, error) {
	order := &model.Order{
		OrderSN:   utils.GenerateOrderSN(userID),
		UserID:    userID,
		ProductID: productID,
		Amount:    amount,
		Status:    1,
	}

	// err := s.orderRepo.Create(order)
	// if err != nil {}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}
