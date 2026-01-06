package repository

import (
	"GopherBuy/internal/model"
	"errors"

	"gorm.io/gorm"
)

type FlashSaleRepository struct {
	*baseRepository[model.FlashSale]
}

func NewFlashSaleRepository(db *gorm.DB) *FlashSaleRepository {
	return &FlashSaleRepository{
		baseRepository: (*baseRepository[model.FlashSale])(NewRepository[model.FlashSale](db)),
	}
}

func (r *FlashSaleRepository) DeductStock(tx *gorm.DB, promoId uint64, quantity uint32) error {
	result := tx.Model(&model.FlashSale{}).
		Where("id = ? AND promo_stock >= ?", promoId, quantity).
		UpdateColumn("promo_stock", gorm.Expr("promo_stock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}

	return nil
}
