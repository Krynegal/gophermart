package service

import (
	"context"
	"github.com/Krynegal/gophermart.git/internal/order"
)

type NotEnoughPoints struct {
}

func (n NotEnoughPoints) Error() string {
	return "not enough points on the account"
}

func (s *Service) GetBalance(ctx context.Context, userID int) (float32, float32) {
	accruals := s.storage.GetAccruals(ctx, userID)
	withdrawn := s.storage.GetWithdrawals(ctx, userID)
	return accruals, withdrawn
}

func (s *Service) DeductionOfPoints(ctx context.Context, order *order.WithdrawOrder) error {
	accruals, withdrawn := s.GetBalance(ctx, order.UserID)

	if order.Sum > accruals-withdrawn {
		return NotEnoughPoints{}
	}

	err := s.storage.DeductPoints(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetWithdrawalOfPoints(ctx context.Context, userID int) ([]order.WithdrawOrder, error) {
	orders, err := s.storage.GetWithdrawalOfPoints(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
