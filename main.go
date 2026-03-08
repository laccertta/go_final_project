package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {

	// инициализация БД
	err := db.Init("scheduler.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.DB.Close()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir := "./web"

	// регистрация API обработчиков
	api.Init()

	//сервер
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server started on http://localhost:" + port)

	// запуск
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
	}
}
