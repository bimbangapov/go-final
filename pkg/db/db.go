package db

import (
	"database/sql"
	"os"
)

var (
	DbSchema = "CREATE TABLE IF NOT EXISTS scheduler\n(\n    id INTEGER PRIMARY KEY AUTOINCREMENT,\n    date CHAR(8) NOT NULL DEFAULT \"\",\n    title VARCHAR(256) NOT NULL,\n    comment TEXT,\n    repeat VARCHAR(128)\n);"
	DbIndex  = "CREATE INDEX IF NOT EXISTS scheduler_date on scheduler(date);"
)

type SchedulerStore struct {
	db *sql.DB
}

func NewSchedulerStore(db *sql.DB) (*SchedulerStore, error) {

	schedulerStore := &SchedulerStore{db: db}
	err := schedulerStore.Init()

	if err != nil {
		return nil, err
	}

	return schedulerStore, nil
}

func (s *SchedulerStore) Init() error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(DbSchema)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(DbIndex)
	if err != nil {
		return err
	}
	return nil
}

func NeedInstallDb(dataBase string) bool {
	dbFile := dataBase
	_, err := os.Stat(dbFile)

	install := false
	if err != nil {
		install = true
	}
	return install
}
