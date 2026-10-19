package main

import (
	"database/sql"
	"errors"
	"fmt"
	"go-final/pkg/db"
	"log"
	"os"

	"go-final/pkg/config"
	"go-final/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {
	port := config.GetPortFromEnv()
	dbName := config.GetDBFromEnv()

	password := config.GetPasswordFromEnv()

	err := os.Mkdir("logs", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Println(err)
		return
	}

	flog, err := os.OpenFile("logs/server.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println("Ошибка создания файла для логов: ", err)
		return
	}
	defer flog.Close()

	servLogs := log.New(
		flog,
		`serv`,
		log.LstdFlags|log.Lshortfile)

	dataBase, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Printf("error connecting to database: %v", err)
		servLogs.Println(err)
		return
	}
	defer dataBase.Close()

	store, err := db.NewSchedulerStore(dataBase)
	if err != nil {
		log.Printf("error connecting to database: %v", err)
		servLogs.Println("error connecting to database: ", err)
		return
	}

	myServ, err := server.NewServer(servLogs, port, store, password)
	if err != nil {
		servLogs.Println("error creating server: ", err)
		return
	}

	log.Printf("start serv on port: %d\n", port)
	servLogs.Printf("start serv on port: %d\n", port)
	servLogs.Printf("database coonnected")
	log.Printf("database coonnected")

	if err := myServ.Server.ListenAndServe(); err != nil {
		fmt.Printf("error serving on port: %d\n", port)
		servLogs.Println(err)
		return
	}

}
