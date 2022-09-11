package postgres

import (
	"context"
	"database/sql"
	"github.com/Krynegal/gophermart.git/internal/storage"
	"github.com/Krynegal/gophermart.git/internal/user"
)

func (db *DB) CreateUser(ctx context.Context, user *user.User) (int, error) {
	stmt, err := db.db.PrepareContext(ctx, "INSERT INTO users(login,password) VALUES ($1,$2) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result := stmt.QueryRowContext(ctx, user.Login, user.Password)
	var ID sql.NullInt64
	_ = result.Scan(&ID)
	if !ID.Valid {
		return -1, storage.ErrLogin{Login: user.Login}
	}
	userID := int(ID.Int64)
	return userID, nil
}

func (db *DB) GetUserID(ctx context.Context, user *user.User) (int, error) {
	//fmt.Printf("User login: %v, password: %v\n", user.Login, user.Password)
	row := db.db.QueryRowContext(ctx, "SELECT id FROM users WHERE login=$1 AND password=$2", user.Login, user.Password)
	var ID sql.NullInt64
	_ = row.Scan(&ID)
	if !ID.Valid {
		return -1, storage.ErrAuth
	}
	userID := int(ID.Int64)
	return userID, nil
}
