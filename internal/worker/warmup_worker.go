package worker

import (
	"GopherBuy/internal/repository"
	"context"
	"log"
	"time"
)

type WarmUpWorker struct {
	flashsaleRepo *repository.FlashSaleRepository
	redisRepo     *repository.RedisRepository
}

func NewWarmUpWorker(fsRepo *repository.FlashSaleRepository, rRepo *repository.RedisRepository) *WarmUpWorker {
	return &WarmUpWorker{
		flashsaleRepo: fsRepo,
		redisRepo:     rRepo,
	}
}

func (w *WarmUpWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Warmup Worker Started...")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.scanAndWarmUp()
		}
	}
}

func (w *WarmUpWorker) scanAndWarmUp() {
	now := time.Now()
	upcoming := time.Now().Add(10 * time.Minute)

	flashsales, err := w.flashsaleRepo.GetUpcomingFlashSales(now, upcoming)
	if err != nil {
		log.Printf("Failed to get upcoming flashsales: %v", err)
		return
	}

	for _, fs := range flashsales {
		if fs.IsWarmed == true {
			continue
		}

		err := w.redisRepo.SyncFlashSaleFromDB(repository.FlashSaleDTO{
			PromoID:     fs.ID,
			PromoStock:  fs.PromoStock,
			MaxPurchase: fs.MaxPurchase,
			StartTime:   fs.StartTime,
			EndTime:     fs.EndTime,
		})

		if err != nil {
			log.Printf("Failed to warm up PromoID: %d: %v", fs.ID, err)
		}

		fs.IsWarmed = true
	}
}
