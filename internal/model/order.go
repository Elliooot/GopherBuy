package model

import "time"

type Order struct {
	ID        uint64 `gorm:"primarykey"`
	OrderSN   string `gorm:"uniqueIndex;type:varchar(64)"`
	UserID    uint64 `gorm:"index"`
	ProductID uint64 `gorm:"index"`
	Quantity  uint32
	Amount    float32
	Status    int32 // 1: Pending, 2: Finished, 3: Cancelled
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}
