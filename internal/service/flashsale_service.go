package service

import (
	"GopherBuy/api"
	"GopherBuy/internal/model"
	"GopherBuy/internal/repository"
	"time"
)

type FlashSaleService struct {
	flashsaleRepo *repository.FlashSaleRepository
}

func NewFlashSaleService() *FlashSaleService {
	return &FlashSaleService{}
}

func (s *FlashSaleService) CreateFlashSale(req *api.CreateFlashSaleRequest) (*api.CreateFlashSaleResponse, error) {
	newFS := &model.FlashSale{
		ProductID:   req.ProductId,
		PromoPrice:  req.PromoPrice,
		PromoStock:  req.PromoStock,
		MaxPurchase: req.MaxPurchase,
		StartTime:   time.Unix(req.StartTime, 0),
		EndTime:     time.Unix(req.EndTime, 0),
	}

	err := s.flashsaleRepo.Create(newFS)
	if err != nil {
		return nil, err
	}

	return &api.CreateFlashSaleResponse{
		Msg: "Successfully create flashsale!",
	}, nil
}
