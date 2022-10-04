package handlers

import (
	"context"
	"encoding/json"
	"github.com/Krynegal/gophermart.git/internal/order"
	"github.com/Krynegal/gophermart.git/internal/service"
	"io"
	"net/http"
)

type balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

//getCurrentBalance GET /api/user/balance - получение текущего баланса пользователя
func (r *Router) getCurrentBalance(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := r.getUserIDFromToken(w, req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	accruals, withdraws := r.Service.GetBalance(ctx, userID)

	b := balance{Current: accruals - withdraws, Withdrawn: withdraws}

	output, err := json.Marshal(b)
	if err != nil {
		http.Error(w, "internalServerError", http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

//deductionOfPoints POST /api/user/balance/withdraw - запрос на списание средств
func (r *Router) deductionOfPoints(w http.ResponseWriter, req *http.Request) {
	userID, err := r.getUserIDFromToken(w, req)
	if err != nil {
		return
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "incorrect input data", http.StatusInternalServerError)
		return
	}

	var order *order.WithdrawOrder
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, "internalServerError", http.StatusInternalServerError)
		return
	}

	order.UserID = userID

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = r.Service.DeductionOfPoints(ctx, order)

	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusOK)
	case service.NotEnoughPoints:
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	default:
		http.Error(w, "internalServerError", http.StatusInternalServerError)
	}
}

//getWithdrawalOfPoints GET /api/user/balance/withdrawals - получение информации о выводе средств
func (r *Router) getWithdrawalOfPoints(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := r.getUserIDFromToken(w, req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orders, err := r.Service.GetWithdrawalOfPoints(ctx, userID)
	if err != nil {
		http.Error(w, "internalServerError", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	output, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "internalServerError", http.StatusInternalServerError)
		return
	}
	w.Write(output)
}
