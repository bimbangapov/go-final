package db

import (
	"context"
	"database/sql"
)

type TasksResp struct {
	Tasks []Task `json:"tasks.go"`
}

func (s *SchedulerStore) Tasks(ctx context.Context, limit int) ([]*Task, error) {
	tasks := []*Task{}

	rows, err := s.db.QueryContext(ctx, "SELECT id,date,title,comment,repeat FROM scheduler order by date asc  limit :limit",
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
