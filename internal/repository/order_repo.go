package repository

import (
	"GopherBuy/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func (r *OrderRepository) GetById(orderId uint64) (*model.Order, error) {
	var order model.Order
	err := r.db.First(&order, orderId).Error
	return &order, err
}

func (r *OrderRepository) Create(order *model.Order) error {
	result := r.db.Create(order)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
