package repository

import (
	"GopherBuy/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type FlashSaleRepository struct {
	*baseRepository[model.FlashSale]
}

func NewFlashSaleRepository(db *gorm.DB) *FlashSaleRepository {
	return &FlashSaleRepository{
		baseRepository: NewRepository[model.FlashSale](db),
	}
}

func (r *FlashSaleRepository) Create(fs *model.FlashSale) error {
	result := r.db.Create(fs)
	if result.Error != nil {
		return result.Error
	}
	return nil
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

func (r *FlashSaleRepository) RollBackStock(tx *gorm.DB, promoId uint64, quantity uint32) error {
	result := r.db.Where("id = ?", promoId).
		Update("promo_stock", gorm.Expr("promo_stock + ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *FlashSaleRepository) GetUpcomingFlashSales(now time.Time, upcoming time.Time) ([]model.FlashSale, error) {
	var flashSales []model.FlashSale
	result := r.db.Where("start_time BETWEEN ? AND ?", now, upcoming).Find(&flashSales)
	if result.Error != nil {
		return nil, result.Error
	}
	return flashSales, nil
}
