package api

import (
	"encoding/json"
	"go-final/pkg/config"
	"log"
	"net/http"
	"strings"
	"time"

	"go-final/pkg/db"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	Logger     *log.Logger
	FileServer http.Handler
	Store      *db.SchedulerStore
	password   string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(logger *log.Logger, handler http.Handler, store *db.SchedulerStore, password string) (*Handler, error) {

	return &Handler{
		Logger:     logger,
		FileServer: handler,
		Store:      store,
		password:   password,
	}, nil
}

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.Logger.Printf("Входящий запрос: метод=%s URL=%s", r.Method, r.RequestURI)
		log.Printf("Входящий запрос: метод=%s URL=%s\n", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := h.password
		if len(pass) > 0 {
			var tokenString string
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			token, err := jwt.Parse(
				tokenString,
				func(token *jwt.Token) (any, error) {
					return []byte(h.password), nil
				},
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
				jwt.WithExpirationRequired(),
			)

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func isValidInput(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func (h *Handler) writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch code {
	case 200:
		w.WriteHeader(http.StatusOK)
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 401:
		w.WriteHeader(http.StatusUnauthorized)
	case 405:
		w.WriteHeader(http.StatusMethodNotAllowed)
	case 500:
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(data)
	if err != nil {
		h.Logger.Println(err)
		return
	}

	if _, err := w.Write(resp); err != nil {
		log.Printf("write response: %v", err)
		return
	}
}

func DateAfter(t1, t2 time.Time) bool {
	return t1.Format(config.DatePattern) > t2.Format(config.DatePattern)
}

func DateBeforeOrEqual(t1, t2 time.Time) bool {
	return t1.Format(config.DatePattern) <= t2.Format(config.DatePattern)
}

func DateBefore(t1, t2 time.Time) bool {
	return t1.Format(config.DatePattern) < t2.Format(config.DatePattern)
}

func checkDate(task *db.Task) error {
	now := time.Now().UTC()
	if task.Date == "" {
		task.Date = now.Format(config.DatePattern)
	}

	t, err := time.Parse(config.DatePattern, task.Date)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		if DateBefore(t, now) {
			task.Date = now.Format(config.DatePattern)
		}
	}

	if task.Repeat != "" {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if DateBefore(t, now) {
			task.Date = nextDate
		}
	}

	return nil
}
