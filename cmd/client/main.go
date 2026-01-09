package main

import (
	"GopherBuy/api"
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	client := api.NewOrderServiceClient(conn)

	totalRequests := 100
	promoID := 1
	var successCount int64
	var failCount int64

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 1; i <= totalRequests; i++ {
		wg.Add(1)
		go func(userID uint64) {
			defer wg.Done()

			resp, err := client.CreateFlashOrder(context.Background(), &api.FlashOrderRequest{
				UserId:     userID,
				ProductId:  1,
				PromoId:    uint64(promoID),
				Quantity:   1,
				PromoPrice: 100,
			})

			if err != nil {
				atomic.AddInt64(&failCount, 1)
				if failCount == 1 {
					log.Printf("First Error: %v", err)
				}
			} else {
				if resp.Status == 202 {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&failCount, 1)
				}
			}
		}(uint64(i))
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("-----------------------------------")
	fmt.Printf("Time usage: %v\n", duration)
	fmt.Printf("Total request: %d\n", totalRequests)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Println("-----------------------------------")
}
