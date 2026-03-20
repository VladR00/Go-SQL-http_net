package main

import (
	"Go-SQL-http_net/internal/middleware"
	"Go-SQL-http_net/internal/storage/postgre"
	"Go-SQL-http_net/pkg/config"
	"Go-SQL-http_net/pkg/logger"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	logger.Init()
	cfg, err := config.MustLoadConfig()
	if err != nil {
		slog.Error("Load config", "Error", err)
		log.Fatal(err)
	}
	db, err := postgre.ConnectPostgreSQL(cfg.PostgreSQL)
	if err != nil {
		slog.Error("Load PostgreSQL", "Error", err)
		log.Fatal(err)
	}

	storage := postgre.NewPostgreSQL(db)
	mw := middleware.NewMiddleware(storage)

	http.HandleFunc("/departments/", mw.HandlerCreateDepartment) //POST

	slog.Info("Server start", "Port", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil))
}
