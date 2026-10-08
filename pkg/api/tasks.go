package api

import (
	"encoding/json"
	"go-final/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func (h *Handler) TasksHadler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Store.Tasks(r.Context(), 50)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		resp, err := json.Marshal(ErrorResponse{Error: "StatusInternalServerError"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
		return
	}
	if len(tasks) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(TasksResp{make([]*db.Task, 0)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(TasksResp{Tasks: tasks})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
	return
}
