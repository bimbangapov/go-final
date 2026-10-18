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
		h.writeJSON(w, 400, ErrorResponse{Error: "id is required"})
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
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
			if _, err := w.Write(resp); err != nil {
				h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
				return
			}
			return
		}
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}
	h.writeJSON(w, 200, task)
}
