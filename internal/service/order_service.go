package service

import (
	"GopherBuy/api"
	"GopherBuy/internal/model"
	"GopherBuy/internal/repository"
	"GopherBuy/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	// dependency injection
	productRepo   *repository.ProductRepository
	orderRepo     *repository.OrderRepository
	flashsaleRepo *repository.FlashSaleRepository
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

func (s *OrderService) CreateFlashOrder(req *api.FlashOrderRequest) (*api.OrderResponse, error) {
	flashsale, err := s.flashsaleRepo.GetById(req.PromoId)
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

	var order *model.Order // Define a nil pointer
	// var order = &model.Order{} // Define an empty structure, a bit more waste memory

	err = s.orderRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		// Deduct Stock
		if err := s.flashsaleRepo.DeductStock(req.PromoId, req.Quantity); err != nil {
			return err
		}

		// Create Order
		order = &model.Order{
			OrderSN:   utils.GenerateOrderSN(req.UserId),
			UserID:    req.UserId,
			ProductID: req.ProductId,
			Quantity:  req.Quantity,
			Amount:    float32(req.Quantity) * flashsale.PromoPrice,
			Status:    1,
			CreatedAt: time.Now(),
		}

		return tx.Create(order).Error
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
