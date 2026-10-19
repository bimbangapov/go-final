package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (s *SchedulerStore) AddTask(ctx context.Context, task *Task) (int64, error) {
	var id int64

	query := fmt.Sprintf("INSERT INTO scheduler (date,title,comment,repeat) values (:date,:title,:comment,:repeat)")

	res, err := s.db.ExecContext(ctx, query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
