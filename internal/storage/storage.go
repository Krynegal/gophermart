package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Krynegal/gophermart.git/internal/order"
	"github.com/Krynegal/gophermart.git/internal/user"
)

type ErrLogin struct {
	Login string
}

func (logErr ErrLogin) Error() string {
	return fmt.Sprintf("user with login %s is already exist", logErr.Login)
}

var ErrAuth = errors.New("invalid Login or Password")

type Storager interface {
	CreateUser(ctx context.Context, login, password string) (int, error)
	GetUserID(ctx context.Context, user *user.User) (int, error)

	SaveOrder(ctx context.Context, order *order.AccrualOrder) error
	GetUserIDByNumberOrder(ctx context.Context, number uint64) int
	GetUploadedOrders(ctx context.Context, userID int) ([]order.AccrualOrder, error)

	GetAccruals(ctx context.Context, UserID int) float32
	GetWithdrawals(ctx context.Context, UserID int) float32
	DeductPoints(ctx context.Context, order *order.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]order.WithdrawOrder, error)

	GetOrdersForProcessing(poolSize int) ([]string, error)
	UpdateOrderState(order *order.ProcessingOrder) error
}
