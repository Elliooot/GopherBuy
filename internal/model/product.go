package model

import "time"

type Product struct {
	ID        uint64 `gorm:"primarykey"`
	Name      string
	Stock     int32
	Price     int32
	StartTime time.Time
	EndTime   time.Time
	Version   int
}
