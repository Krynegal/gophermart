package errormodels

import (
	"errors"
	"fmt"
)

type ErrLogin struct {
	Login string
}

func (logErr ErrLogin) Error() string {
	return fmt.Sprintf("user with login %s is already exist", logErr.Login)
}

var ErrAuth = errors.New("invalid Login or Password")
