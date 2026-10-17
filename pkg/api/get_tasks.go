package api

import (
	"net/http"
	"time"

	"go-final/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func (h *Handler) TasksHadler(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	if search != "" {
		searchDate, err := time.Parse("02.01.2006", search)
		if err != nil { // Ветка поиска про Title
			tasks, err := h.Store.GetTasksByTitle(r.Context(), search, 50)
			if err != nil {
				writeJSON(w, true, ErrorResponse{Error: "StatusInternalServerError"})
				return
			}
			if len(tasks) == 0 {
				writeJSON(w, false, TasksResp{make([]*db.Task, 0)})
				return
			}
			writeJSON(w, false, TasksResp{tasks})
			return
		} else { // Поиска по Date
			search = searchDate.Format("20060102")
			tasks, err := h.Store.GetTasksByDate(r.Context(), search, 50)
			if err != nil {
				writeJSON(w, true, ErrorResponse{Error: "StatusInternalServerError"})
				return
			}
			if len(tasks) == 0 {
				writeJSON(w, false, TasksResp{make([]*db.Task, 0)})
				return
			}
			writeJSON(w, false, TasksResp{Tasks: tasks})
			return
		}

	}
	// ветка поиска с пустым search
	tasks, err := h.Store.GetTasks(r.Context(), 50)

	if err != nil {
		writeJSON(w, true, ErrorResponse{Error: "StatusInternalServerError"})
		return
	}
	if len(tasks) == 0 {
		writeJSON(w, false, TasksResp{make([]*db.Task, 0)})
		return
	}
	writeJSON(w, false, TasksResp{Tasks: tasks})
	return
}
