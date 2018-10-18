package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Moleque/tp_db/forum/controllers"
)

const (
	DB_USER     = "docker"
	DB_PASSWORD = "docker"
	DB_NAME     = "forum"
)

var DB = &controllers.DataBase{}

func main() {
	router := controllers.NewRouter()

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	DB.Connect(dsn)
	defer DB.Disconnect()

	log.Printf("try to start http server http://127.0.0.1:5000")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatalf("listening error:%s", err)
	}
}
