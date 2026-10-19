package db

import (
	"context"
	"database/sql"
	"errors"
	"go-final/pkg/config"
	"strconv"
	"time"
)

func (s *SchedulerStore) UpdateTask(ctx context.Context, task *Task) error {
	idInt, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}
	_, err = time.Parse(config.DatePattern, task.Date)
	if err != nil {
		return err
	}

	res, err := s.db.ExecContext(ctx, "UPDATE scheduler SET date = :date, title=:title, comment=:comment, repeat= :repeat  WHERE id=:id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", idInt),
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("incorrect task id for update")
	}
	return nil
}
