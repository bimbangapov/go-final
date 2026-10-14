package api

import (
	"net/http"
	"time"
)

func (h *Handler) NextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		h.Logger.Printf("method not allowed")
		return
	}

	now := r.URL.Query().Get("now")
	if now == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("now is required"))
		h.Logger.Printf("now is required")
		return
	}
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("now/date must by YYYYMMDD"))
		h.Logger.Printf("now/date must by YYYYMMDD")
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date is required"))
		h.Logger.Printf("date is required")
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repeat is required"))
		h.Logger.Printf("repeat is required")
		return
	}

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	h.Logger.Printf("Входящий запрос GET /api/nextdate")

	w.Header().Set("Content-Type", "application/txt")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
