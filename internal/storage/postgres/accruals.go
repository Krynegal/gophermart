package postgres

import (
	"context"
	"fmt"
	"github.com/Krynegal/gophermart.git/internal/order"
	"github.com/Krynegal/gophermart.git/internal/status"
	"time"
)

func (db *DB) SaveOrder(ctx context.Context, order *order.AccrualOrder) (err error) {
	order.UploadedAt = time.Now()

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			txError := tx.Rollback()
			if txError != nil {
				err = fmt.Errorf("accruals SaveOrder rollback error %s: %s", txError.Error(), err.Error())
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO orders(order_num,user_id) VALUES ($1,$2)", order.Number, order.UserID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO accruals(order_num,user_id,status,uploaded_at) VALUES ($1,$2,$3,$4)",
		order.Number, order.UserID, order.Status.String(), order.UploadedAt)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) GetUserIDByNumberOrder(ctx context.Context, number uint64) int {
	row := db.db.QueryRowContext(ctx, "SELECT user_id FROM accruals WHERE order_num=$1", number)
	var userID int
	_ = row.Scan(&userID)

	return userID
}

func (db *DB) GetUploadedOrders(ctx context.Context, userID int) ([]order.AccrualOrder, error) {
	rows, err := db.db.QueryContext(ctx, "SELECT order_num,status,amount,uploaded_at FROM accruals WHERE user_id =$1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []order.AccrualOrder
	for rows.Next() {
		var ord order.AccrualOrder
		var stat string
		err = rows.Scan(&ord.Number, &stat, &ord.Accrual, &ord.UploadedAt)
		if err != nil {
			return nil, err
		}
		ord.Status, err = status.GetStatus(stat)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (db *DB) GetOrdersForProcessing(pool int) ([]string, error) {
	var orders []string
	rows, err := db.db.Query(
		"SELECT order_num FROM accruals WHERE status IN ($1, $2) ORDER BY uploaded_at LIMIT $3", "NEW", "PROCESSING", pool,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	for rows.Next() {
		var orderID string
		if err = rows.Scan(&orderID); err != nil {
			return orders, err
		}
		orders = append(orders, orderID)
	}
	err = rows.Err()
	return orders, err
}

func (db *DB) UpdateOrderState(order *order.ProcessingOrder) error {
	_, err := db.db.Exec(
		"UPDATE accruals SET status=$1, amount=$2 WHERE order_num = ($3)",
		order.Status, order.Accrual, order.Order,
	)
	return err
}
