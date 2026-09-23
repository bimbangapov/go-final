package db

import (
	"database/sql"
	"log"
	"os"
)

var (
	DbSchema = "CREATE TABLE scheduler\n(\n    id INTEGER PRIMARY KEY AUTOINCREMENT,\n    date CHAR(8) NOT NULL DEFAULT \"\",\n    title VARCHAR(256) NOT NULL,\n    comment TEXT,\n    repeat VARCHAR(128)\n);"
	DbIndex  = "CREATE INDEX scheduler_date on scheduler(date);"
)

type SchedulerStore struct {
	db  *sql.DB
	log *log.Logger
}

func NewSchedulerStore(db *sql.DB, log *log.Logger) (*SchedulerStore, error) {
	return &SchedulerStore{
		db:  db,
		log: log,
	}, nil
}

func (s *SchedulerStore) Init() error {
	err := s.db.Ping()
	if err != nil {
		return err
	}
	s.log.Printf("scheduler store initialized, succesfull ping")

	_, err = s.db.Exec(DbSchema)
	if err != nil {
		return err
	}
	s.log.Printf("scheduler store initialized, succesfull create table")

	_, err = s.db.Exec(DbIndex)
	if err != nil {
		return err
	}
	s.log.Printf("scheduler store initialized, succesfull create index")
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
