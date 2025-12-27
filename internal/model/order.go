package model

import "time"

type Order struct {
	ID        uint64 `gorm:"primarykey"`
	OrderSN   string `gorm:"uniqueIndex;type:varchar(64)"`
	UserID    uint64 `gorm:"index"`
	ProductID uint64 `gorm:"index"`
	Amount    float32
	Status    int32 // 1: Unpaid, 2: Paid, 3: Finished
	PaymentID string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}
