package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/Krynegal/gophermart.git/internal/models"
	"github.com/dgrijalva/jwt-go"
)

const secretKey = "be55d1079e6c6167118ac91318fe"

type AuthRepo interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUserID(ctx context.Context, user *models.User) (int, error)
}

type AuthService struct {
	repo AuthRepo
}

func NewAuthService(repo AuthRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (auth *AuthService) CreateUser(ctx context.Context, user *models.User) error {
	user.Password = auth.generatePasswordHash(user.Password)
	userID, err := auth.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (auth *AuthService) AuthenticationUser(ctx context.Context, user *models.User) error {
	user.Password = auth.generatePasswordHash(user.Password)
	userID, err := auth.repo.GetUserID(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (auth *AuthService) GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})
	tokenString, err := token.SignedString([]byte("KSFjH$53KSFjH6745u#uEQQjF349%835hFpzA"))
	return tokenString, err
}

func (auth *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(secretKey)))
}
