package db

import (
	"context"
	"database/sql"
)

func (s *SchedulerStore) GetTask(ctx context.Context, id int) (*Task, error) {
	var task Task

	row := s.db.QueryRowContext(ctx, "SELECT id,date,title,comment,repeat FROM scheduler WHERE id = :id order by date asc limit 1",
		sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
