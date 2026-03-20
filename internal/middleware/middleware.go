package middleware

import (
	"Go-SQL-http_net/internal/handlers"
	"Go-SQL-http_net/internal/storage/postgre"
	"log/slog"
	"net/http"
)

type Middleware struct {
	Storage handlers.StorageHandler
}

func NewMiddleware(storage *postgre.Storage) *Middleware {
	return &Middleware{Storage: handlers.StorageHandler{Storage: storage}}
}

func (mw *Middleware) HandlerCreateDepartment(w http.ResponseWriter, r *http.Request) {
	if msg, err := mw.Storage.HandlerCreateDepartment(w, r); err != nil {
		slog.Error("HandlerCreateDepartment", msg, err)
	} else {
		slog.Info("HandlerCreateDepartment", "Message", msg)
	}
}
