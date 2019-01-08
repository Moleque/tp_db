package main

import (
	"fmt"
	"net/http"

	"tp_db/forum/database"

	"tp_db/forum/controllers"
)

const (
	DB_USER     = "docker"
	DB_PASSWORD = "docker"
	DB_NAME     = "docker"
)

func main() {
	router := controllers.NewRouter()

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	database.DB.Connect(dsn)
	defer database.DB.Disconnect()

	http.ListenAndServe(":5000", router)
}
