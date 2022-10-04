package service

import (
	"context"
	"github.com/Krynegal/gophermart.git/internal/order"
	"github.com/Krynegal/gophermart.git/internal/status"
)

type OrderAlreadyUploadedCurrentUserError struct{}

func (o OrderAlreadyUploadedCurrentUserError) Error() string {
	return "the order has already been uploaded by the current user"
}

type OrderAlreadyUploadedAnotherUserError struct{}

func (o OrderAlreadyUploadedAnotherUserError) Error() string {
	return "the order has already been uploaded by another user"
}

type CheckLuhnError struct{}

func (c CheckLuhnError) Error() string {
	return "invalid order number format"
}

func (s *Service) LoadOrder(ctx context.Context, numOrder uint64, userID int) error {

	if !s.CheckLuhn(numOrder) {
		return CheckLuhnError{}
	}

	ord := order.AccrualOrder{
		Number: numOrder,
		UserID: userID,
		Status: status.StatusNEW,
	}

	userIDinDB := s.storage.GetUserIDByNumberOrder(ctx, ord.Number)
	if userIDinDB == 0 {
		return s.storage.SaveOrder(ctx, &ord)
	}
	if userIDinDB == ord.UserID {
		return OrderAlreadyUploadedCurrentUserError{}
	}
	return OrderAlreadyUploadedAnotherUserError{}
}

func (s *Service) CheckLuhn(number uint64) bool {
	var sum uint64

	for i := 0; number > 0; i++ {
		cur := number % 10
		if i%2 == 0 {
			sum += cur
			number = number / 10
			continue
		}
		cur = cur * 2
		if cur > 9 {
			cur = cur - 9
		}
		sum += cur
		number = number / 10
	}

	return sum%10 == 0
}

func (s *Service) GetUploadedOrders(ctx context.Context, userID int) ([]order.AccrualOrder, error) {
	orders, err := s.storage.GetUploadedOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
