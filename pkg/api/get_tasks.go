package api

import (
	"go-final/pkg/config"
	"net/http"
	"time"

	"go-final/pkg/db"
)

const limit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func (h *Handler) TasksHadler(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	if search != "" {
		searchDate, err := time.Parse("02.01.2006", search)
		if err != nil { // Ветка поиска про Title
			tasks, err := h.Store.GetTasksByTitle(r.Context(), search, limit)
			if err != nil {
				h.writeJSON(w, 500, ErrorResponse{Error: "StatusInternalServerError"})
				return
			}

			h.writeJSON(w, 200, TasksResp{tasks})
			return
		} else { // Поиска по Date
			search = searchDate.Format(config.DatePattern)
			tasks, err := h.Store.GetTasksByDate(r.Context(), search, limit)
			if err != nil {
				h.writeJSON(w, 500, ErrorResponse{Error: "StatusInternalServerError"})
				return
			}
			if len(tasks) == 0 {
				h.writeJSON(w, 200, TasksResp{make([]*db.Task, 0)})
				return
			}
			h.writeJSON(w, 200, TasksResp{Tasks: tasks})
			return
		}

	}
	// ветка поиска с пустым search
	tasks, err := h.Store.GetTasks(r.Context(), limit)

	if err != nil {
		h.writeJSON(w, 500, ErrorResponse{Error: "StatusInternalServerError"})
		return
	}
	if len(tasks) == 0 {
		h.writeJSON(w, 500, TasksResp{make([]*db.Task, 0)})
		return
	}
	h.writeJSON(w, 200, TasksResp{Tasks: tasks})
	return
}
