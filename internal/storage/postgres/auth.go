package postgres

import (
	"context"
	"database/sql"
	"github.com/Krynegal/gophermart.git/internal/models"
	"github.com/Krynegal/gophermart.git/internal/models/errormodels"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (a *AuthPostgres) CreateUser(ctx context.Context, user *models.User) (int, error) {
	stmt, err := a.db.PrepareContext(ctx, "INSERT INTO users(login,password) VALUES ($1,$2) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result := stmt.QueryRowContext(ctx, user.Login, user.Password)
	var ID sql.NullInt64
	_ = result.Scan(&ID)
	if !ID.Valid {
		return -1, errormodels.ErrLogin{Login: user.Login}
	}
	userID := int(ID.Int64)
	return userID, nil
}

func (a *AuthPostgres) GetUserID(ctx context.Context, user *models.User) (int, error) {
	//fmt.Printf("User login: %v, password: %v\n", user.Login, user.Password)
	row := a.db.QueryRowContext(ctx, "SELECT id FROM users WHERE login=$1 AND password=$2", user.Login, user.Password)
	var ID sql.NullInt64
	_ = row.Scan(&ID)
	if !ID.Valid {
		return -1, errormodels.ErrAuth
	}
	userID := int(ID.Int64)
	return userID, nil
}
