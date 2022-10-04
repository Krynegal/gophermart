package handlers

import (
	"context"
	"errors"
	"github.com/Krynegal/gophermart.git/internal/storage"
	"github.com/Krynegal/gophermart.git/internal/user"
	"net/http"
)

func (r *Router) registration(w http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong data format"))
		return
	}
	var u user.User
	err := r.readingUserData(w, req, &u)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var userID int
	userID, err = r.Service.CreateUser(ctx, u.Login, u.Password)

	if errors.As(err, &storage.ErrLogin{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	r.writeToken(w, userID)
}

func (r *Router) authentication(w http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong data format"))
		return
	}
	var u user.User
	err := r.readingUserData(w, req, &u)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = r.Service.AuthenticationUser(ctx, &u)

	if errors.Is(err, storage.ErrAuth) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	r.writeToken(w, u.ID)
}
