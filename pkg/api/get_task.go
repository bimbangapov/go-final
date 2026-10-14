package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, true, ErrorResponse{Error: "id is required"})
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.Store.GetTask(r.Context(), idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			resp, err := json.Marshal(ErrorResponse{Error: "task not found"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(resp)
			return
		}
	}
	writeJSON(w, false, task)
}
