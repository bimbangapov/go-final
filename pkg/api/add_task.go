package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go-final/pkg/db"
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
		h.Logger.Printf(http.StatusText(http.StatusBadRequest)+": %s", err.Error())
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		h.Logger.Printf(http.StatusText(http.StatusBadRequest)+": %s", err.Error())
		return
	}

	//Провекрка заполненности title
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

	id, err := h.Store.AddTask(ctx, &task)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		h.Logger.Printf(http.StatusText(http.StatusBadRequest)+": %s", err.Error())
		return
	}
	h.writeJSON(w, 200, SuccessResponseAddTask{id})

}
