package handlers

import (
	"github.com/Krynegal/gophermart.git/internal/rest/middlewares"
	"github.com/Krynegal/gophermart.git/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
	Service service.Servicer
}

func NewRouter(service service.Servicer) *Router {
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

	r.Router.Handle("/api/user/orders", middlewares.AuthMiddleware(r.loadOrders)).Methods(http.MethodPost)
	r.Router.Handle("/api/user/orders", middlewares.AuthMiddleware(r.getUploadedOrders)).Methods(http.MethodGet)

	r.Router.Handle("/api/user/balance", middlewares.AuthMiddleware(r.getCurrentBalance)).Methods(http.MethodGet)
	r.Router.Handle("/api/user/balance/withdraw", middlewares.AuthMiddleware(r.deductionOfPoints)).Methods(http.MethodPost)
	r.Router.Handle("/api/user/withdrawals", middlewares.AuthMiddleware(r.getWithdrawalOfPoints)).Methods(http.MethodGet)
}
