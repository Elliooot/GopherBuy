package repository

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisOnce sync.Once

func InitRedis() {
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         "localhost:6379", // Maybe read from .env?
			Password:     "",               // Maybe read from .env?
			DB:           0,
			PoolSize:     100,
			MinIdleConns: 10,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
		})

		ctx := context.Background()
		_, err := redisClient.Ping(ctx).Result()

		if err != nil {
			log.Fatalf("Cannot connect to Redis: %v\n", err)
		}
		fmt.Println("Successfully connect to Redis!")
	})
}

func GetRedis() *redis.Client {
	if redisClient == nil {
		panic("Redis not initialized, call InitRedis() first")
	}
	return redisClient
}
