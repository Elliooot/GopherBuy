package repository

import (
	"GopherBuy/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	*baseRepository[model.Order]
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		baseRepository: NewRepository[model.Order](db),
	}
}

func (r *OrderRepository) Create(order *model.Order) error {
	result := r.db.Create(order)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
