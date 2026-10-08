package api

import (
	"bytes"
	"encoding/json"
	"go-final/pkg/db"
	"net/http"
)

type SuccessResponseAddTask struct {
	Id int64 `json:"id"`
}

func (h *Handler) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		task db.Task
		buf  bytes.Buffer
		ctx  = r.Context()
	)

	//Проверка корректности входного JSON
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Logger.Printf(http.StatusText(http.StatusBadRequest)+": %s", err.Error())
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Logger.Printf(http.StatusText(http.StatusBadRequest)+": %s", err.Error())
		return
	}

	//Провекрка заполненности title
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
	id, err := h.Store.AddTask(ctx, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Logger.Printf(http.StatusText(http.StatusBadRequest)+": %s", err.Error())
		return
	}
	writeJSON(w, false, SuccessResponseAddTask{id})

}
