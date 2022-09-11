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
	var user user.User
	err := r.readingUserData(w, req, &user)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = r.Service.CreateUser(ctx, &user)

	if errors.As(err, &storage.ErrLogin{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	r.writeToken(w, &user)
}

func (r *Router) authentication(w http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong data format"))
		return
	}
	var user user.User
	err := r.readingUserData(w, req, &user)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = r.Service.AuthenticationUser(ctx, &user)

	if errors.Is(err, storage.ErrAuth) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	r.writeToken(w, &user)
}
