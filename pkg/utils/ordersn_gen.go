package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderSN(userID uint64) string {
	now := time.Now().Format("20060102150405")
	suffix := fmt.Sprintf("%04d", userID%10000)
	randomNum := rand.Intn(9000) + 1000
	return fmt.Sprintf("%s%s%d", now, suffix, randomNum)
}
