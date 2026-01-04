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

func (r *RedisRepository) GetPromoStock(promoId uint64) uint32 {
	ctx := context.Background()
	key := fmt.Sprintf("flashsale:stock:%d", promoId)

	result, err := r.rdb.Get(ctx, key).Int()
	if err != nil {
		fmt.Printf("Problem of getting stock")
	}

	return uint32(result)
}

func (r *RedisRepository) DeductStock(promoId uint64, quantity uint32) error {
	ctx := context.Background()
	key := fmt.Sprintf("flashsale:stock:%d", promoId)

	// Embeded lua script
	script := `
		local stock = tonumber(redis.call("GET", KEYS[1]) or '0')
		
		if tonumber(ARGV[1]) > stock then
			return -1
		end
		
		redis.call("DECRBY", KEYS[1], ARGV[1])
		return 0
	`

	result, err := r.rdb.Eval(ctx, script, []string{key}, quantity).Int()

	if err != nil {
		return err
	}

	if result == -1 {
		return ErrStockInsufficient
	}

	return nil
}
