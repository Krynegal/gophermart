package postgres

import (
	"context"
	"fmt"
	"github.com/Krynegal/gophermart.git/internal/order"
	"time"
)

func (db *DB) GetAccruals(ctx context.Context, UserID int) float32 {
	row := db.db.QueryRowContext(ctx, "SELECT SUM(amount) FROM public.accruals WHERE user_id=$1", UserID)
	var accruals float32
	_ = row.Scan(&accruals)

	return accruals
}

func (db *DB) GetWithdrawals(ctx context.Context, UserID int) float32 {
	row := db.db.QueryRowContext(ctx, "SELECT SUM(amount) FROM public.withdrawals WHERE user_id=$1", UserID)
	var withdrawals float32
	_ = row.Scan(&withdrawals)

	return withdrawals
}

func (db *DB) DeductPoints(ctx context.Context, order *order.WithdrawOrder) (err error) {
	order.ProcessedAt = time.Now()

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			txError := tx.Rollback()
			if txError != nil {
				err = fmt.Errorf("balance DeductPoints rollback error %s: %s", txError.Error(), err.Error())
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO public.orders(order_num,user_id) VALUES ($1,$2)", order.Order, order.UserID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO public.withdrawals(order_num,user_id,amount,processed_at) VALUES ($1,$2,$3,$4)",
		order.Order, order.UserID, order.Sum, order.ProcessedAt)

	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) GetWithdrawalOfPoints(ctx context.Context, userID int) ([]order.WithdrawOrder, error) {
	rows, err := db.db.QueryContext(ctx, "SELECT order_num,amount,processed_at FROM public.withdrawals WHERE user_id =$1 ORDER BY processed_at", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []order.WithdrawOrder
	for rows.Next() {
		var order order.WithdrawOrder
		err = rows.Scan(&order.Order, &order.Sum, &order.ProcessedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
