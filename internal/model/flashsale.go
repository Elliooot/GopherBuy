package model

import "time"

type FlashSale struct {
	ID          uint64 `gorm:"primarykey"`
	ProductID   uint64 `gorm:"index"`
	PromoPrice  float32
	PromoStock  uint32
	MaxPurchase uint32
	StartTime   time.Time
	EndTime     time.Time
	IsWarmed    bool `gorm:"default:false"`
}
