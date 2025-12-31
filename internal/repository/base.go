package repository

import "gorm.io/gorm"

type baseRepository[T any] struct {
	db *gorm.DB
}

func (r *baseRepository[T]) GetDB() *gorm.DB {
	return r.db
}

func (r *baseRepository[T]) GetById(id uint64) (*T, error) {
	var t T
	err := r.db.First(&t, id).Error
	return &t, err
}

// Generic Constructor, cannot be implemented as a method
func NewRepository[T any](db *gorm.DB) *baseRepository[T] {
	return &baseRepository[T]{db: db}
}
