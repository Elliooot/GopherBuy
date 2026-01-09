package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrStockInsufficient = errors.New("Insufficient Stock")

type RedisRepository struct {
	rdb *redis.Client
}

type FlashSaleDTO struct {
	PromoId     uint64
	PromoStock  uint32
	MaxPurchase uint32
	StartTime   time.Time
	EndTime     time.Time
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

func (r *RedisRepository) SyncFlashSaleFromDB(flashSale FlashSaleDTO) error {
	ctx := context.Background()
	key := fmt.Sprintf("flashsale:stock:%d", flashSale.PromoId)

	err := r.rdb.HSet(ctx, key, map[string]interface{}{
		"promo_stock":  flashSale.PromoStock,
		"max_purchase": flashSale.MaxPurchase,
		"start_time":   flashSale.StartTime.Unix(),
		"end_time":     flashSale.EndTime.Unix(),
	}).Err()

	if err != nil {
		return err
	}

	// Set expiring time
	ttl := time.Until(flashSale.EndTime) + time.Hour
	return r.rdb.Expire(ctx, key, ttl).Err()
}

func (r *RedisRepository) GetFlashSale(promoId uint64) (*FlashSaleDTO, error) {
	ctx := context.Background()
	key := fmt.Sprintf("flashsale:stock:%d", promoId)

	result, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		fmt.Printf("Failed to get flashsale")
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("flashsale not found")
	}

	promoStock, _ := strconv.ParseUint(result["promo_stock"], 10, 32)
	maxPurchase, _ := strconv.ParseUint(result["max_purchase"], 10, 32)
	startTime, _ := strconv.ParseInt(result["start_time"], 10, 64)
	endTime, _ := strconv.ParseInt(result["end_time"], 10, 64)

	return &FlashSaleDTO{
		PromoId:     promoId,
		PromoStock:  uint32(promoStock),
		MaxPurchase: uint32(maxPurchase),
		StartTime:   time.Unix(startTime, 0),
		EndTime:     time.Unix(endTime, 0),
	}, nil
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

func (r *RedisRepository) RollBackStock(promoId uint64, quantity uint32) error {
	ctx := context.Background()
	key := fmt.Sprintf("flashsale:stock:%d", promoId)

	result, err := r.rdb.IncrBy(ctx, key, int64(quantity)).Result()

	if err != nil {
		return err
	}

	fmt.Printf("Successfully rollback stock to %v\n", result)
	return nil
}
