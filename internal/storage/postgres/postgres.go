package postgres

import (
	"database/sql"
	"github.com/Krynegal/gophermart.git/internal/configs"
	"github.com/Krynegal/gophermart.git/internal/storage"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func NewDatabase(config *configs.Config) (storage.Storager, error) {
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
		db: db,
	}, nil
}
