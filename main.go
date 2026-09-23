package main

import (
	"database/sql"
	"errors"
	"fmt"
	"go-final/pkg/config"
	"go-final/pkg/db"
	"go-final/pkg/server"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	port := config.GetPortFromEnv()
	dbName := config.GetDBFromEnv()
	needDbInstall := db.NeedInstallDb(dbName)

	err := os.Mkdir("logs", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
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

	dbLogs := log.New(
		flog,
		`db`,
		log.LstdFlags|log.Lshortfile)

	myServ := server.NewServer(servLogs, port)

	dataBase, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Printf("error connecting to database: %v", err)
		dbLogs.Fatal(err)
	}
	defer dataBase.Close()

	store, err := db.NewSchedulerStore(dataBase, dbLogs)
	if err != nil {
		log.Printf("error connecting to database: %v", err)
		dbLogs.Fatal(err)
	}

	if needDbInstall {
		if err := store.Init(); err != nil {
			log.Printf("error connecting to database: %v", err)
			dbLogs.Fatal(err)
		}
	}

	log.Printf("start serv on port: %d\n", port)
	servLogs.Printf("start serv on port: %d\n", port)
	dbLogs.Printf("database coonnected")
	log.Printf("database coonnected")

	if err := myServ.Server.ListenAndServe(); err != nil {
		fmt.Printf("error serving on port: %d\n", port)
		servLogs.Fatal(err)
	}

}
