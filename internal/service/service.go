package service

import (
	"context"
	"github.com/Krynegal/gophermart.git/internal/configs"
	"github.com/Krynegal/gophermart.git/internal/order"
	"github.com/Krynegal/gophermart.git/internal/storage"
	"github.com/Krynegal/gophermart.git/internal/user"
)

type Servicer interface {
	CreateUser(ctx context.Context, login, password string) (int, error)
	AuthenticationUser(ctx context.Context, user *user.User) error
	GenerateToken(uid int) (string, error)

	LoadOrder(ctx context.Context, numOrder uint64, userID int) error
	GetUploadedOrders(ctx context.Context, userID int) ([]order.AccrualOrder, error)

	GetBalance(ctx context.Context, userID int) (float32, float32)
	DeductionOfPoints(ctx context.Context, order *order.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]order.WithdrawOrder, error)
}

type Service struct {
	storage storage.Storager
	config  *configs.Config
}

func NewService(storage storage.Storager, config *configs.Config) Servicer {
	return &Service{
		storage: storage,
		config:  config,
	}
}
