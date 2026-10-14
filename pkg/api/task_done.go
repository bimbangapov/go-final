package api

import (
	"encoding/json"
	"errors"
	"go-final/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, true, ErrorResponse{Error: "id is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, true, ErrorResponse{Error: "id must be digit"})
		return
	}

	task, err := h.Store.GetTask(r.Context(), idInt)
	if err != nil {
		writeJSON(w, true, ErrorResponse{Error: err.Error()})
		return
	}

	if task.Repeat != "" { // ветка получения след даты
		taskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			writeJSON(w, true, ErrorResponse{Error: err.Error()})
			return
		}

		if DateAfter(taskDate, time.Now().UTC()) || task.Date == time.Now().UTC().Format("20060102") { // Если дата в будущем
			task.Date, err = NextDate(time.Now().UTC(), task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, true, ErrorResponse{Error: err.Error()})
				return
			}
		} else { // если дата сегодня/в прошлом
			err := checkDate(task)
			if err != nil {
				writeJSON(w, true, ErrorResponse{Error: err.Error()})
				return
			}
		}
		err = h.Store.UpdateTask(r.Context(), task)
		if err != nil {
			writeJSON(w, true, ErrorResponse{err.Error()})
			return
		}
		writeJSON(w, false, struct{}{})
		return
	}

	if task.Repeat == "" { // ветка удаления задачи
		err := h.Store.DeleteTask(r.Context(), idInt)
		if errors.Is(err, db.ErrTaskNotFound) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			resp, err := json.Marshal(&ErrorResponse{"task not found"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(resp)
			return
		}
		if err != nil {
			writeJSON(w, true, ErrorResponse{err.Error()})
			return
		}
		writeJSON(w, false, struct{}{})
	}

}
