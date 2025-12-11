package services

import (
	"context"

	"github.com/universeh2h/report/internal/repositories"
)

type ModulService struct {
	repo *repositories.ModulOtomax
}

func NewModulService(repo *repositories.ModulOtomax) *ModulService {
	return &ModulService{
		repo: repo,
	}
}
func (s *ModulService) GetAllModulOtomax(c context.Context, date string) ([]repositories.ModulType, error) {
	data, err := s.repo.GetAllModulOtomax(c, date)
	if err != nil {
		return nil, err
	}
	return data, nil
}
