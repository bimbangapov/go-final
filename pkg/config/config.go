package config

import (
	"os"
	"strconv"
)

func GetPortFromEnv() int {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		return 7540
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return 7540
	}
	return intPort
}

func GetDBFromEnv() string {
	db := os.Getenv("TODO_DBFILE")
	if db == "" {
		return "scheduler.db"
	}
	return db
}
