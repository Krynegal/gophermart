package service

import (
	"context"
	"github.com/Krynegal/gophermart.git/internal/models"
	"github.com/Krynegal/gophermart.git/internal/storage"
)

type Auth interface {
	CreateUser(ctx context.Context, user *models.User) error
	AuthenticationUser(ctx context.Context, user *models.User) error
	GenerateToken(user *models.User) (string, error)
}

type AccrualOrder interface {
}

type WithdrawOrder interface {
}

type Service struct {
	Auth     Auth
	Accrual  AccrualOrder
	Withdraw WithdrawOrder
}

func NewService(repo *storage.Repository) *Service {
	return &Service{
		Auth:     NewAuthService(repo.Auth),
		Accrual:  nil,
		Withdraw: nil,
	}
}
