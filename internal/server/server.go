package server

import (
	"github.com/Krynegal/gophermart.git/internal/configs"
	"net/http"
)

type Server struct {
	HTTPServer *http.Server
}

func NewServer(cfg *configs.Config, router http.Handler) *Server {
	return &Server{
		HTTPServer: &http.Server{
			Addr:    cfg.RunAddr,
			Handler: router,
		},
	}
}
