package services

import (
	"context"

	"github.com/universeh2h/report/internal/repositories"
)

type TransactionsService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionsService(repo *repositories.TransactionRepository) *TransactionsService {
	return &TransactionsService{
		repo: repo,
	}
}

func (s *TransactionsService) GetTransactions(c context.Context) ([]repositories.TransactionData, error) {
	return s.repo.GetTransactions(c)
}
