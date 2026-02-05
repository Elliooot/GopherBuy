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

	products := []model.Product{
		{ID: 1, Name: "Test iPhone", Price: 1000},
		{ID: 101, Name: "iPhone 15", Price: 100},
		{ID: 102, Name: "MacBook Pro", Price: 2000},
	}

	for _, p := range products {
		if err := repository.DB.Create(&p).Error; err != nil {
			log.Fatalf("Failed to create product %d: %v", p.ID, err)
		}
		log.Printf("Created Product ID: %d - %s", p.ID, p.Name)
	}

	flashSales := []model.FlashSale{
		{
			ID:          1,
			ProductID:   1,
			PromoPrice:  1000,
			PromoStock:  17,
			MaxPurchase: 1,
			StartTime:   time.Now().UTC().Add(-1 * time.Hour),
			EndTime:     time.Now().UTC().Add(24 * time.Hour),
		},
		{
			ID:          2,
			ProductID:   101,
			PromoPrice:  100,
			PromoStock:  100000,
			MaxPurchase: 1,
			StartTime:   time.Now().UTC().Add(-1 * time.Hour),
			EndTime:     time.Now().UTC().Add(24 * time.Hour),
		},
		{
			ID:          3,
			ProductID:   102,
			PromoPrice:  2000,
			PromoStock:  5000,
			MaxPurchase: 2,
			StartTime:   time.Now().UTC().Add(-1 * time.Hour),
			EndTime:     time.Now().UTC().Add(24 * time.Hour),
		},
	}
	// repository.DB.Create(&flashSale)

	redisRepo := repository.NewRedisRepository(repository.GetRedis())

	for _, fs := range flashSales {
		if err := repository.DB.Create(&fs).Error; err != nil {
			log.Fatalf("Failed to create flashsale %d: %v", fs.ID, err)
		}

		err := redisRepo.SyncFlashSaleFromDB(repository.FlashSaleDTO{
			PromoID:     fs.ID,
			PromoStock:  fs.PromoStock,
			MaxPurchase: fs.MaxPurchase,
			StartTime:   fs.StartTime,
			EndTime:     fs.EndTime,
		})

		if err != nil {
			log.Fatalf("Failed to warm up Redis for flashsale %d: %v", fs.ID, err)
		}

		log.Printf("Created FlashSale ID: %d for Product: %d, Stock: %d", fs.ID, fs.ProductID, fs.PromoStock)
	}

	log.Println("All test data set")
}
