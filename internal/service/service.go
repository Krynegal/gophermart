package service

import (
	"context"
	"github.com/Krynegal/gophermart.git/internal/order"
	"github.com/Krynegal/gophermart.git/internal/storage"
	"github.com/Krynegal/gophermart.git/internal/user"
)

type Servicer interface {
	CreateUser(ctx context.Context, user *user.User) error
	AuthenticationUser(ctx context.Context, user *user.User) error
	GenerateToken(user *user.User) (string, error)

	LoadOrder(ctx context.Context, numOrder uint64, userID int) error
	GetUploadedOrders(ctx context.Context, userID int) ([]order.AccrualOrder, error)

	GetBalance(ctx context.Context, userID int) (float32, float32)
	DeductionOfPoints(ctx context.Context, order *order.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]order.WithdrawOrder, error)
}

type Service struct {
	storage storage.Storager
}

func NewService(storage storage.Storager) Servicer {
	return &Service{
		storage: storage,
	}
}
