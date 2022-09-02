package storage

import (
	"context"
	"database/sql"
	"github.com/Krynegal/gophermart.git/internal/models"
	"github.com/Krynegal/gophermart.git/internal/storage/postgres"
)

type Auth interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUserID(ctx context.Context, user *models.User) (int, error)
}

type AccrualOrder interface {
}

type WithdrawOrder interface {
}

type Repository struct {
	Auth     Auth
	Accrual  AccrualOrder
	Withdraw WithdrawOrder
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Auth:     postgres.NewAuthPostgres(db),
		Accrual:  nil,
		Withdraw: nil,
	}
}
