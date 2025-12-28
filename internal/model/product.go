package model

type Product struct {
	ID          uint64 `gorm:"primarykey"`
	Name        string
	Stock       uint32
	Price       float32
	MaxPurchase uint32
}
