package repository

import (
	"GopherBuy/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func (r *OrderRepository) Create(order *model.Order) error {
	result := r.db.Create(order)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
