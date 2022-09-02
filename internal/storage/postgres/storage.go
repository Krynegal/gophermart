package postgres

import (
	"database/sql"
	"github.com/Krynegal/gophermart.git/internal/configs"
	_ "github.com/lib/pq"
)

type DB struct {
	DB *sql.DB
}

func NewDatabase(config *configs.Config) (*DB, error) {
	db, err := sql.Open("postgres", config.DatabaseURI)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	tables := []string{
		UsersTable, OrdersTable, AccrualsTable, WithdrawalsTable,
	}
	for _, table := range tables {
		if _, err = db.Exec(table); err != nil {
			return nil, err
		}
	}
	return &DB{
		DB: db,
	}, nil
}
