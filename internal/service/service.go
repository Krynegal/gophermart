package service

import (
	"context"
	"github.com/Krynegal/gophermart.git/internal/storage"
	"github.com/Krynegal/gophermart.git/internal/user"
)

type Servicer interface {
	CreateUser(ctx context.Context, user *user.User) error
	AuthenticationUser(ctx context.Context, user *user.User) error
	GenerateToken(user *user.User) (string, error)
}

type Service struct {
	storage storage.Storager
}

func NewService(storage storage.Storager) *Service {
	return &Service{
		storage: storage,
	}
}
