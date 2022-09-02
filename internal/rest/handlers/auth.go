package handlers

import (
	"context"
	"errors"
	"github.com/Krynegal/gophermart.git/internal/models"
	"github.com/Krynegal/gophermart.git/internal/models/errormodels"
	"net/http"
)

func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong data format"))
		return
	}
	var user models.User
	err := h.readingUserData(w, r, &user)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.Auth.CreateUser(ctx, &user)

	if errors.As(err, &errormodels.ErrLogin{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.writeToken(w, &user)
}

func (h *Handler) authentication(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong data format"))
		return
	}
	var user models.User
	err := h.readingUserData(w, r, &user)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.Auth.AuthenticationUser(ctx, &user)

	if errors.Is(err, errormodels.ErrAuth) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.writeToken(w, &user)
}
