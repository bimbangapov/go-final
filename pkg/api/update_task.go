package api

import (
	"bytes"
	"encoding/json"
	"go-final/pkg/db"
	"net/http"
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
		writeJSON(w, true, respErr)
		return
	}

	if err := checkDate(&task); err != nil {
		respErr := &ErrorResponse{
			err.Error(),
		}
		writeJSON(w, true, respErr)
		return
	}

	err = h.Store.UpdateTask(r.Context(), &task)
	if err != nil {
		writeJSON(w, true, ErrorResponse{Error: "Error updating task"})
		return
	}
	writeJSON(w, false, struct{}{})
}
