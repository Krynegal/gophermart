package handlers

import (
	"encoding/json"
	"github.com/Krynegal/gophermart.git/internal/models"
	"io"
	"net/http"
)

func (h *Handler) readingUserData(w http.ResponseWriter, r *http.Request, user *models.User) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "empty login or password", http.StatusBadRequest)
		return err
	}
	return nil
}

func (h *Handler) writeToken(w http.ResponseWriter, user *models.User) {
	token, err := h.Service.Auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Can't get token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
}
