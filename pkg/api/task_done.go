package api

import (
	"encoding/json"
	"errors"
	"go-final/pkg/config"
	"net/http"
	"strconv"
	"time"

	"go-final/pkg/db"
)

func (h *Handler) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		h.writeJSON(w, 400, ErrorResponse{Error: "id is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: "id must be digit"})
		return
	}

	task, err := h.Store.GetTask(r.Context(), idInt)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}

	if task.Repeat != "" { // ветка получения след даты
		taskDate, err := time.Parse(config.DatePattern, task.Date)
		if err != nil {
			h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
			return
		}

		if DateAfter(taskDate, time.Now().UTC()) || task.Date == time.Now().UTC().Format(config.DatePattern) { // Если дата в будущем
			task.Date, err = NextDate(time.Now().UTC(), task.Date, task.Repeat)
			if err != nil {
				h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
				return
			}
		}

		err = h.Store.UpdateTask(r.Context(), task)
		if err != nil {
			h.writeJSON(w, 400, ErrorResponse{err.Error()})
			return
		}
		h.writeJSON(w, 200, struct{}{})
		return
	}

	if task.Repeat == "" { // ветка удаления задачи
		err := h.Store.DeleteTask(r.Context(), idInt)
		if errors.Is(err, db.ErrTaskNotFound) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			resp, err := json.Marshal(&ErrorResponse{"task not found"})
			if err != nil {
				h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
				return
			}
			if _, err := w.Write(resp); err != nil {
				h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
				return
			}
			return
		}
		if err != nil {
			h.writeJSON(w, 400, ErrorResponse{err.Error()})
			return
		}
		h.writeJSON(w, 200, struct{}{})
	}

}
