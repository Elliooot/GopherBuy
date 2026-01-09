package main

import (
	"log"
	"os"
	"time"

	"GopherBuy/internal/model"
	"GopherBuy/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("REDIS_HOST", "localhost")

	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	if err := repository.InitRedis(); err != nil {
		log.Fatalf("Failed to init Redis: %v", err)
	}

	repository.DB.Exec("TRUNCATE TABLE flash_sales, orders, products RESTART IDENTITY")

	product := model.Product{ID: 1, Name: "Test iPhone", Price: 1000}
	repository.DB.Create(&product)

	flashSale := model.FlashSale{
		ID:          1,
		ProductID:   1,
		PromoPrice:  1,
		PromoStock:  17,
		MaxPurchase: 1,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
	}
	repository.DB.Create(&flashSale)

	redisRepo := repository.NewRedisRepository(repository.GetRedis())
	err := redisRepo.SyncFlashSaleFromDB(repository.FlashSaleDTO{
		PromoID:     flashSale.ID,
		PromoStock:  flashSale.PromoStock,
		MaxPurchase: flashSale.MaxPurchase,
		StartTime:   flashSale.StartTime,
		EndTime:     flashSale.EndTime,
	})

	if err != nil {
		log.Fatalf("Failed to warm up Redis: %v\n", err)
	}
	log.Println("Finished data population and warmup")
}
