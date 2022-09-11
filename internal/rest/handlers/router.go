package handlers

import (
	"github.com/Krynegal/gophermart.git/internal/rest/middlewares"
	"github.com/Krynegal/gophermart.git/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
	Service *service.Service
}

func NewRouter(service *service.Service) *Router {
	router := &Router{
		Router:  mux.NewRouter(),
		Service: service,
	}
	router.InitRoutes()
	return router
}

func (r *Router) InitRoutes() {
	r.Router.HandleFunc("/api/user/register", r.registration).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/user/login", r.authentication).Methods(http.MethodPost)

	r.Router.Handle("/api/user/hello", middlewares.AuthMiddleware(r.Hello))
}
