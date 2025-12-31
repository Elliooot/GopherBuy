package repository

import (
	"GopherBuy/internal/model"
	"errors"

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

func (r *ProductRepository) DeductStock(productId uint64, quantity uint32) error {
	result := r.db.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", productId, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}

	return nil
}
