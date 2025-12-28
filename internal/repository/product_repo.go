package repository

import (
	"GopherBuy/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func (r *ProductRepository) GetPriceById(productId uint64) (float32, error) {
	var product model.Product
	err := r.db.First(&product, productId).Error
	if err != nil {
		return 0, err
	}
	return product.Price, nil
}
