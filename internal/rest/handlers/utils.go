package handlers

import (
	"encoding/json"
	"github.com/Krynegal/gophermart.git/internal/user"
	"io"
	"net/http"
)

func (r *Router) readingUserData(w http.ResponseWriter, req *http.Request, user *user.User) error {
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
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

func (r *Router) writeToken(w http.ResponseWriter, user *user.User) {
	token, err := r.Service.GenerateToken(user)
	if err != nil {
		http.Error(w, "Can't get token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
}
