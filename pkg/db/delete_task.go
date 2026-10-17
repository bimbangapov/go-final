package db

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

func (s *SchedulerStore) DeleteTask(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM scheduler WHERE id=:id",
		sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}
