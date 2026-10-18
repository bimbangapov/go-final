package api

import (
	"go-final/pkg/config"
	"net/http"
	"time"
)

func (h *Handler) NextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.Printf("method not allowed")
		h.writeJSON(w, 405, ErrorResponse{Error: "method not allowed"})

		return
	}

	now := r.URL.Query().Get("now")
	if now == "" {
		now = time.Now().Format(config.DatePattern)
	}
	nowTime, err := time.Parse(config.DatePattern, now)
	if err != nil {
		h.Logger.Printf("now/date must by YYYYMMDD")
		h.writeJSON(w, 400, ErrorResponse{Error: "now/date must by YYYYMMDD"})
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		h.Logger.Printf("date is required")
		h.writeJSON(w, 400, ErrorResponse{Error: "date is required"})
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		h.Logger.Printf("repeat is required")
		h.writeJSON(w, 400, ErrorResponse{Error: "repeat is required"})
		return
	}

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}
	h.Logger.Printf("Входящий запрос GET /api/nextdate")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(nextDate)); err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}
}
