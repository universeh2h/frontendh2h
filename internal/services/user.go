package services

import (
	"context"

	"github.com/universeh2h/report/internal/model"
	"github.com/universeh2h/report/internal/repositories"
	"github.com/universeh2h/report/pkg/config"
)

type UserService struct {
	repo *repositories.AuthRepository
}

func NewUserService(repo *repositories.AuthRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Login(c context.Context, req model.Login) (*model.User, string, error) {
	user, err := s.repo.Login(c, req)
	if err != nil {
		return nil, "", err
	}
	token, err := config.GenerateJWT(user.Username)

	return user, token, err
}
