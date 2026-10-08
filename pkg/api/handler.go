package api

import (
	"database/sql"
	"encoding/json"
	"go-final/pkg/db"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	Logger     *log.Logger
	FileServer http.Handler
	Store      *db.SchedulerStore
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(logger *log.Logger, handler http.Handler, database *sql.DB) (*Handler, error) {

	store, err := db.NewSchedulerStore(database)
	if err != nil {
		return &Handler{}, err
	}

	return &Handler{
		Logger:     logger,
		FileServer: handler,
		Store:      store,
	}, nil
}

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if r.Header.Get("Sec-Fetch-Dest") == "document" {
		h.Logger.Printf("Входящий запрос: метод=%s URL=%s", r.Method, r.RequestURI)
		log.Printf("Входящий запрос: метод=%s URL=%s\n", r.Method, r.RequestURI)
		//}
		next.ServeHTTP(w, r)
	})
}

func isValidInput(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func writeJSON(w http.ResponseWriter, isError bool, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if isError {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func DateAfter(t1, t2 time.Time) bool {
	return t1.Format("20060102") > t2.Format("20060102")
}

func DateBefore(t1, t2 time.Time) bool {
	return t1.Format("20060102") < t2.Format("20060102")
}

func checkDate(task *db.Task) error {
	now := time.Now().UTC()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		if DateBefore(t, now) {
			task.Date = now.Format("20060102")
		}
	}

	if task.Repeat != "" && DateBefore(t, now) {
		task.Date, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	return nil
}
