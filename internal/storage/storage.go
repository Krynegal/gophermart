package storage

import (
	"context"
	"errors"
	"fmt"
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
	CreateUser(ctx context.Context, user *user.User) (int, error)
	GetUserID(ctx context.Context, user *user.User) (int, error)
}
