package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrStockInsufficient = errors.New("Insufficient Stock")

type RedisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

func (r *RedisRepository) InitRedisRepo(promoId uint64, stock uint32) error {
	ctx := context.Background()
	key := fmt.Sprintf("flashsale:stock:%d", promoId)

	return r.rdb.Set(ctx, key, stock, 24*time.Hour).Err()
}

func (r *RedisRepository) DeductStock(promoId uint64, quantity uint32) error {
	ctx := context.Background()
	key := fmt.Sprintf("flashsale:stock:%d", promoId)

	// Implement Optimisitc Lock using Watch(), Bad Example
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		stock, err := r.rdb.Get(ctx, key).Int()

		if err != nil {
			return err
		}

		if quantity > uint32(stock) {
			return ErrStockInsufficient
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.DecrBy(ctx, key, int64(quantity))
			return nil
		})
		return err
	}, key)
	return err
}
