package service

import (
	"GopherBuy/api"
	"GopherBuy/internal/model"
	"GopherBuy/internal/repository"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCreateFlashOrder_Concurrency(t *testing.T) {
	repository.InitDB()

	// start
	repository.DB.Exec("DELETE FROM orders")
	repository.DB.Exec("DELETE FROM flash_sales")
	repository.DB.Exec("DELETE FROM products")

	product := model.Product{ID: 1, Name: "Test iPhone", Price: 1000}
	repository.DB.Create(&product)

	flashSale := model.FlashSale{
		ID:          1,
		ProductID:   1,
		PromoPrice:  1,
		PromoStock:  3,
		MaxPurchase: 1,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(1 * time.Hour),
	}
	repository.DB.Create(&flashSale)
	// end

	svc := &OrderService{
		flashsaleRepo: repository.NewFlashSaleRepository(repository.DB),
		orderRepo:     repository.NewOrderRepository(repository.DB),
	}

	promoID := uint64(1)
	repository.DB.Model(&model.FlashSale{}).Where("id = ?", promoID).Update("promo_stock", 3)

	var wg sync.WaitGroup
	userCount := 10
	results := make(chan string, userCount)

	for i := 1; i <= userCount; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			req := &api.FlashOrderRequest{
				UserId:    uint64(userID),
				ProductId: 1,
				PromoId:   promoID,
				Quantity:  1,
			}
			resp, err := svc.CreateFlashOrder(req)
			if err != nil {
				results <- "Error"
			} else if resp.Status == 201 {
				results <- "Success"
			} else {
				results <- "Failed"
			}
		}(i)
	}

	wg.Wait()
	close(results)

	successCount := 0
	for res := range results {
		if res == "Success" {
			successCount++
		}
	}

	fmt.Printf("Successful order count: %d\n", successCount)
	if successCount > 3 {
		t.Errorf("Oversold! Sold %d orders in total.", successCount)
	}
}
