package main

import (
	"flag"
	"github.com/Krynegal/gophermart.git/internal/configs"
	"github.com/Krynegal/gophermart.git/internal/rest/handlers"
	"github.com/Krynegal/gophermart.git/internal/server"
	"github.com/Krynegal/gophermart.git/internal/service"
	"github.com/Krynegal/gophermart.git/internal/storage"
	"github.com/Krynegal/gophermart.git/internal/storage/postgres"
	"log"
	"net/http"
)

func main() {
	cfg := configs.NewConfigs()

	flag.StringVar(&cfg.RunAddr, "a", cfg.RunAddr, "Run server address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database URI")
	flag.StringVar(&cfg.AccrualSysAddr, "r", cfg.AccrualSysAddr, "Accrual system address")
	flag.Parse()

	db, err := postgres.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Can't crate the storage: ", err)
	}

	repository := storage.NewRepository(db.DB)
	svc := service.NewService(repository)
	handler := handlers.NewHandler(svc)

	srv := server.NewServer(cfg, handler)

	log.Printf("server run on address: %s", cfg.RunAddr)
	log.Fatal(http.ListenAndServe(cfg.RunAddr, srv.HTTPServer.Handler))
}
