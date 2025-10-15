package services

import (
	"context"

	"github.com/universeh2h/report/internal/model"
	"github.com/universeh2h/report/internal/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductServices(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) Analytics(c context.Context, req model.PaginationParams) (*model.AnalyticsResult, error) {
	return s.repo.TotalTransaksi(c, req)
}
func (s *ProductService) TransaksiReseller(c context.Context, kodeReseller string, startDate, endDate string) ([]model.TransaksiDetails, error) {
	return s.repo.TransaksiReseller(c, kodeReseller, startDate, endDate)
}

func (s *ProductService) GetTrxTercuan(c context.Context, startDate string, endDate string, kodeReseller string) ([]repositories.TrxTerCuan, error) {
	return s.repo.GetTrxTercuan(c, startDate, endDate, kodeReseller)
}

func (s *ProductService) GetProductTrxTerbanyak(c context.Context, startDate string, endDate string, kodeReseller string) ([]repositories.ProductResponse, error) {
	return s.repo.GetProductTrxTerbanyak(c, startDate, endDate, kodeReseller)
}

func (s *ProductService) GetTotalProfit(c context.Context, startDate string, endDate string) ([]model.TopProductsBestSeller, error) {
	return s.repo.Report(c, startDate, endDate, true)
}
