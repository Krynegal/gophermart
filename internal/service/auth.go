package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/Krynegal/gophermart.git/internal/user"
	"github.com/dgrijalva/jwt-go"
)

const secretKey = "be55d1079e6c6167118ac91318fe"

func (s *Service) CreateUser(ctx context.Context, user *user.User) error {
	user.Password = s.generatePasswordHash(user.Password)
	userID, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (s *Service) AuthenticationUser(ctx context.Context, user *user.User) error {
	user.Password = s.generatePasswordHash(user.Password)
	userID, err := s.storage.GetUserID(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (s *Service) GenerateToken(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})
	tokenString, err := token.SignedString([]byte("KSFjH$53KSFjH6745u#uEQQjF349%835hFpzA"))
	return tokenString, err
}

func (s *Service) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(secretKey)))
}
