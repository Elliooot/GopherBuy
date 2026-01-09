package main

import (
	"GopherBuy/api"
	"GopherBuy/internal/handler"
	"GopherBuy/internal/repository"
	"GopherBuy/internal/service"
	"GopherBuy/internal/worker"
	"context"
	"log"
	"net"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type Dependencies struct {
	OrderService     *service.OrderService
	AsynqRepo        *repository.AsynqRepository
	OrderConsumer    *worker.OrderConsumer
	WarmUpWorker     *worker.WarmUpWorker
	GrpcOrderHandler *handler.GrpcOrderHandler
}

func main() {
	// Init
	initInfrastructure()
	deps := initDependencies()

	go deps.OrderConsumer.StartOrderConsumer(10)
	go deps.WarmUpWorker.Start(context.Background())

	go startAsynqWorker(deps.AsynqRepo)

	startGRPCServerNonBlocking(deps.GrpcOrderHandler)
}

func initInfrastructure() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	if err := repository.InitRedis(); err != nil {
		log.Fatalf("Failed to init Redis: %v", err)
	}

	kafkaConfig := &repository.KafkaConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   os.Getenv("KAFKA_TOPIC"),
	}

	if err := repository.InitKafka(kafkaConfig); err != nil {
		log.Fatalf("Failed to init Kafka: %v", err)
	}

	asynqConfig := &repository.AsynqConfig{
		Addr:     os.Getenv("ASYNQ_REDIS_HOST") + ":" + os.Getenv("ASYNQ_REDIS_PORT"),
		Password: os.Getenv("ASYNQ_REDIS_PASSWORD"),
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
	warmUpWorker := worker.NewWarmUpWorker(flashsaleRepo, redisRepo)

	grpcOrderhandler := handler.NewGrpcOrderHandler(orderService)

	return &Dependencies{
		OrderService:     orderService,
		AsynqRepo:        asynqRepo,
		OrderConsumer:    orderConsumer,
		WarmUpWorker:     warmUpWorker,
		GrpcOrderHandler: grpcOrderhandler,
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

func startGRPCServerNonBlocking(orderHandler *handler.GrpcOrderHandler) *grpc.Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterOrderServiceServer(grpcServer, orderHandler)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
		log.Println("gRPC Server started on port 50051")
	}()

	return grpcServer
}
