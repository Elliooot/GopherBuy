package repository

import (
	"GopherBuy/internal/model"
	"errors"

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

func (r *FlashSaleRepository) DeductStock(promoId uint64, quantity uint32) error {
	result := r.db.Model(&model.FlashSale{}).
		Where("id = ? AND promostock >= ?", promoId, quantity).
		UpdateColumn("promostock", gorm.Expr("promostock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}

	return nil
}
