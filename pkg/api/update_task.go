package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go-final/pkg/db"
)

func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		task db.Task
		buf  bytes.Buffer
	)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !isValidInput(task.Title) {
		respErr := &ErrorResponse{
			"title обязательный для заполнения",
		}
		h.writeJSON(w, 400, respErr)
		return
	}

	if err := checkDate(&task); err != nil {
		respErr := &ErrorResponse{
			err.Error(),
		}
		h.writeJSON(w, 400, respErr)
		return
	}

	err = h.Store.UpdateTask(r.Context(), &task)
	if err != nil {
		h.writeJSON(w, 500, ErrorResponse{Error: "Error updating task"})
		return
	}
	h.writeJSON(w, 200, struct{}{})
}
