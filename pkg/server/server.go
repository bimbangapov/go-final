package server

import (
	"database/sql"
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

func NewServer(log *log.Logger, port int, db *sql.DB) (*Server, error) {
	fileServer := http.FileServer(http.Dir("./web"))
	h, err := api.NewHandler(log, fileServer, db)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(h.LoggingMiddleware)

	r.Handle("/*", h.FileServer)
	r.Get("/api/nextdate", h.NextDayHandler)
	r.Post("/api/task", h.AddTaskHandler)
	r.Get("/api/tasks", h.TasksHadler)

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
	}, nil
}
