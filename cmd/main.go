package main

import (
	"Go-SQL-http_net/internal/middleware"
	"Go-SQL-http_net/pkg/logger"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	logger.Init()
	mw := middleware.NewMiddleware(1)

	http.HandleFunc("/departments/", mw.HandlerCreateDepartment) //POST

	slog.Info("Server start at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
