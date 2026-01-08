package main

import (
	"GopherBuy/internal/repository"
	"GopherBuy/internal/service"
	"GopherBuy/internal/worker"
	"log"

	"github.com/hibiken/asynq"
)

type Dependencies struct {
	OrderService  *service.OrderService
	AsynqRepo     *repository.AsynqRepository
	orderConsumer *worker.OrderConsumer
}

func main() {
	// Init
	initInfrastructure()
	deps := initDependencies()

	go deps.orderConsumer.StartOrderConsumer(10)

	go startAsynqWorker(deps.AsynqRepo)
}

func initInfrastructure() {
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	if err := repository.InitRedis(); err != nil {
		log.Fatalf("Failed to init Redis: %v", err)
	}

	kafkaConfig := &repository.KafkaConfig{
		Brokers: []string{"localhost:9092", "localhost:9094"},
		Topic:   "order_create",
	}

	if err := repository.InitKafka(kafkaConfig); err != nil {
		log.Fatalf("Failed to init Kafka: %v", err)
	}

	asynqConfig := &repository.AsynqConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	}

	if err := repository.InitAsynq(asynqConfig); err != nil {
		log.Fatalf("Failed to init Asynq: %v", err)
	}
}

func initDependencies() *Dependencies {
	redis := repository.GetRedis()
	kafkaProducer := repository.GetKafkaProducer()
	kafkaConsumer := repository.GetKafkaConsumer()
	asynqClient := repository.GetAsynqClient()

	productRepo := repository.NewProductRepository(repository.DB)
	orderRepo := repository.NewOrderRepository(repository.DB)
	flashsaleRepo := repository.NewFlashSaleRepository(repository.DB)

	redisRepo := repository.NewRedisRepository(redis)
	kafkaRepo := repository.NewKafkaRepository(kafkaProducer, kafkaProducer.Topic)

	asynqRepo := repository.NewAsynqRepository(asynqClient, orderRepo, redisRepo, flashsaleRepo)

	orderService := service.NewOrderService(productRepo, orderRepo, flashsaleRepo, redisRepo, kafkaRepo, asynqRepo)

	orderConsumer := worker.NewOrderConsumer(kafkaConsumer, orderService)

	return &Dependencies{
		OrderService:  orderService,
		AsynqRepo:     asynqRepo,
		orderConsumer: orderConsumer,
	}
}

func startAsynqWorker(asynqRepo *repository.AsynqRepository) {
	server := repository.GetAsynqServer()

	mux := asynq.NewServeMux()
	mux.HandleFunc(repository.TypeOrderExpired, asynqRepo.HandleOrderExpiration)

	log.Println("Starting Asynq Worker...")
	if err := server.Run(mux); err != nil {
		log.Fatalf("Failed to start Asynq worker: %v", err)
	}
}
