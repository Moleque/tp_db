package controllers

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DataBase struct {
	instance   *sql.DB
	connection bool
}

var DB = &DataBase{}

func (db *DataBase) Connect(dsn string) {
	var err error
	db.instance, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Errorf("database connecting:%s", err)
	}
	db.instance.SetMaxOpenConns(10)
	err = db.instance.Ping()
	if err != nil {
		fmt.Errorf("database connecting:%s", err)
	}
	db.connection = true
	log.Println("database was connected")
}

func (db *DataBase) Disconnect() {
	db.instance.Close()
	db.connection = false
	log.Println("database was disconnected")
}

func (db *DataBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if !db.connection {
		return nil, fmt.Errorf("database query")
	}
	return db.instance.Query(query, args...)
}
