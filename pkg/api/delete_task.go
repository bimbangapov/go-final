package api

import (
	"net/http"
	"strconv"
)

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, true, ErrorResponse{Error: "id is required"})
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, true, ErrorResponse{Error: "id is invalid"})
		return
	}

	err = h.Store.DeleteTask(r.Context(), idInt)
	if err != nil {
		writeJSON(w, true, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, false, struct{}{})
}
