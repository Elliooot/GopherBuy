package repository

import (
	"GopherBuy/internal/model"

	"gorm.io/gorm"
)

type FlashSaleRepository struct {
	db *gorm.DB
}

func (r *FlashSaleRepository) GetById(flashsaleId uint64) (*model.FlashSale, error) {
	var flashsale model.FlashSale
	err := r.db.First(&flashsale, flashsaleId).Error
	return &flashsale, err
}
