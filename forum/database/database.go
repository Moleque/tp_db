package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DataBase struct {
	instance   *sql.DB
	connection bool
}

var DB = &DataBase{}

func (db *DataBase) Connect(dsn string) {
	db.instance, _ = sql.Open("postgres", dsn)
	db.instance.SetMaxOpenConns(10)
	db.instance.Ping()
	db.connection = true
}

func (db *DataBase) Disconnect() {
	db.instance.Close()
	db.connection = false
}

func (db *DataBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if !db.connection {
		return nil, fmt.Errorf("database query")
	}
	return db.instance.Query(query, args...)
}

func (db *DataBase) QueryRow(query string, args ...interface{}) *sql.Row {
	if !db.connection {
		return nil
	}
	return db.instance.QueryRow(query, args...)
}

func (db *DataBase) Prepare(query string) (*sql.Stmt, error) {
	if !db.connection {
		return nil, nil
	}
	return db.instance.Prepare(query)
}

func (db *DataBase) Exec(query string, args ...interface{}) (sql.Result, error) {
	if !db.connection {
		return nil, fmt.Errorf("database query")
	}
	return db.instance.Exec(query, args...)
}
