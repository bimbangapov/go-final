package server

import (
	"fmt"
	"go-final/pkg/api"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Logger *log.Logger
	Server *http.Server
}

func NewServer(log *log.Logger, port int) *Server {
	h := &api.Handler{
		Logger:     log,
		FileServer: http.FileServer(http.Dir("./web")),
	}

	r := chi.NewRouter()
	r.Use(h.LoggingMiddleware)

	r.Handle("/*", h.FileServer)

	var server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ErrorLog:     log,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: log,
		Server: server,
	}
}
