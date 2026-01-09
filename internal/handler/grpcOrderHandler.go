package handler

import (
	"GopherBuy/api"
	"GopherBuy/internal/service"
	"context"
)

type GrpcOrderHandler struct {
	api.UnimplementedOrderServiceServer
	orderService *service.OrderService
}

func NewGrpcOrderHandler(orderService *service.OrderService) *GrpcOrderHandler {
	return &GrpcOrderHandler{
		orderService: orderService,
	}
}

func (h *GrpcOrderHandler) CreateOrder(ctx context.Context, req *api.OrderRequest) (*api.CreateOrderResponse, error) {
	return h.orderService.CreateOrder(req)
}

func (h *GrpcOrderHandler) CreateFlashOrder(ctx context.Context, req *api.FlashOrderRequest) (*api.CreateOrderResponse, error) {
	return h.orderService.HandleFlashSaleRequest(req)
}

func (h *GrpcOrderHandler) GetOrder(ctx context.Context, req *api.GetOrderRequest) (*api.GetOrderResponse, error) {
	return h.orderService.GetOrder(req)
}
