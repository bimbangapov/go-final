package db

import (
	"context"
	"database/sql"
)

type TasksResp struct {
	Tasks []Task `json:"tasks.go"`
}

func (s *SchedulerStore) GetTasks(ctx context.Context, limit int) ([]*Task, error) {
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

func (s *SchedulerStore) GetTasksByTitle(ctx context.Context, title string, limit int) ([]*Task, error) {
	tasks := []*Task{}
	rows, err := s.db.QueryContext(ctx, "SELECT id,date,title,comment,repeat FROM scheduler WHERE title like :title OR comment like :title order by date asc limit :limit",
		sql.Named("title", "%"+title+"%"),
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
func (s *SchedulerStore) GetTasksByDate(ctx context.Context, date string, limit int) ([]*Task, error) {
	tasks := []*Task{}

	rows, err := s.db.QueryContext(ctx, "SELECT id,date,title,comment,repeat FROM scheduler WHERE date like :date order by date asc limit :limit",
		sql.Named("date", date),
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
