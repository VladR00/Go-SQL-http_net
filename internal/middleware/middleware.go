package middleware

import (
	"Go-SQL-http_net/internal/handlers"
	"Go-SQL-http_net/internal/storage/postgre"
	"log/slog"
	"net/http"
)

type Middleware struct {
	handler handlers.StorageHandler
}

func NewMiddleware(storage *postgre.Storage) *Middleware {
	return &Middleware{
		handler: handlers.StorageHandler{Storage: storage},
	}
}

// Handler is the main middleware handler that routes all requests
func (mw *Middleware) Handler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Incoming request", "method", r.Method, "path", r.URL.Path)

	// Call the storage handler
	mw.handler.Handler(w, r)
}
