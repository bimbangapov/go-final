package api

import (
	"net/http"
	"strconv"
)

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		h.writeJSON(w, 400, ErrorResponse{Error: "id is required"})
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: "id is invalid"})
		return
	}

	err = h.Store.DeleteTask(r.Context(), idInt)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}
	h.writeJSON(w, 200, struct{}{})
}
