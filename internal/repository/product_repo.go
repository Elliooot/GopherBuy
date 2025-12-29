package repository

import (
	"GopherBuy/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

// Canonical style, can use a generics type in the future
func (r *ProductRepository) GetById(productId uint64) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, productId).Error
	return &product, err
}
