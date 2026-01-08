package repository

import (
	"fmt"
	"log"
	"sync"

	"github.com/hibiken/asynq"
)

var (
	asynqClient *asynq.Client
	asynqServer *asynq.Server
	asynqOnce   sync.Once
)

type AsynqConfig struct {
	Addr     string
	Password string
	DB       int
}

func InitAsynq(config *AsynqConfig) error {
	var err error
	asynqOnce.Do(func() {
		redisConnOpt := asynq.RedisClientOpt{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		}

		asynqClient = asynq.NewClient(redisConnOpt)

		asynqServer = asynq.NewServer(redisConnOpt, asynq.Config{
			Concurrency: 10,
		})

		fmt.Println("Successfully connect to Asynq!")
	})
	return err
}

func GetAsynqClient() *asynq.Client {
	if asynqClient == nil {
		panic("Asynq client not initialized, call InitAsynq() first")
	}
	return asynqClient
}

func GetAsynqServer() *asynq.Server {
	if asynqServer == nil {
		panic("Asynq server not initialized, call InitAsynq() first")
	}
	return asynqServer
}

func CloseAsynq() error {
	var err error

	if asynqClient != nil {
		if err = asynqClient.Close(); err != nil {
			err = fmt.Errorf("failed to close asynq client: %w", err)
		}
	}

	if asynqServer != nil {
		asynqServer.Shutdown()
	}

	if err != nil {
		return err
	}

	log.Println("Asynq closed")
	return nil
}
