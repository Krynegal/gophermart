package handlers

import (
	"context"
	"encoding/json"
	"github.com/Krynegal/gophermart.git/internal/service"
	"io"
	"net/http"
	"strconv"
)

func (r *Router) loadOrders(w http.ResponseWriter, req *http.Request) {
	userID, err := r.getUserIDFromToken(w, req)
	if err != nil {
		return
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internalServerError", http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "incorrect input data", http.StatusBadRequest)
		return
	}
	strBody := string(body)
	numOrder, err := strconv.ParseUint(strBody, 0, 64)
	if err != nil {
		http.Error(w, "can't get numOrder", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = r.Service.LoadOrder(ctx, numOrder, userID)

	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	case service.OrderAlreadyUploadedCurrentUserError:
		http.Error(w, err.Error(), http.StatusOK)
		return
	case service.OrderAlreadyUploadedAnotherUserError:
		http.Error(w, err.Error(), http.StatusConflict)
		return
	case service.CheckLuhnError:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) getUploadedOrders(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	userID, err := r.getUserIDFromToken(w, req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orders, err := r.Service.GetUploadedOrders(ctx, userID)
	if err != nil {
		http.Error(w, "can't get uploaded orders", http.StatusInternalServerError)
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
