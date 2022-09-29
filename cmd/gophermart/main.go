package main

import (
	"github.com/Krynegal/gophermart.git/internal/client"
	"github.com/Krynegal/gophermart.git/internal/configs"
	"github.com/Krynegal/gophermart.git/internal/rest/handlers"
	"github.com/Krynegal/gophermart.git/internal/server"
	"github.com/Krynegal/gophermart.git/internal/service"
	"github.com/Krynegal/gophermart.git/internal/storage/postgres"
	"log"
	"net/http"
)

func main() {
	cfg := configs.NewConfigs()
	cfg.ParseFlags()

	storage, err := postgres.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Can't crate the storage: ", err)
	}

	cl := client.NewAccrualProcessor(storage, cfg.AccrualSysAddr, 10)
	go cl.Run()

	svc := service.NewService(storage, cfg)
	router := handlers.NewRouter(svc)

	srv := server.NewServer(cfg, router)

	log.Printf("server run on address: %s", cfg.RunAddr)
	log.Fatal(http.ListenAndServe(cfg.RunAddr, srv.HTTPServer.Handler))
}
