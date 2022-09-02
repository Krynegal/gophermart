package handlers

import (
	"github.com/Krynegal/gophermart.git/internal/rest/middlewares"
	"github.com/Krynegal/gophermart.git/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	*mux.Router
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		Router:  mux.NewRouter(),
		Service: service,
	}
	handler.InitRoutes()
	return handler
}

func (h *Handler) InitRoutes() {
	h.Router.HandleFunc("/api/user/register", h.registration).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/user/login", h.authentication).Methods(http.MethodPost)

	h.Router.Handle("/api/user/hello", middlewares.AuthMiddleware(h.Hello))
}
