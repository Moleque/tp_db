package models

import (
	"encoding/json"
	"net/http"
	"tp_db/forum/database"

	"github.com/julienschmidt/httprouter"
)

type Counts struct {
	Post   int32 `json:"post"`
	User   int32 `json:"user"`
	Forum  int32 `json:"forum"`
	Thread int32 `json:"thread"`
}

const selectForumsCount = `
	SELECT COUNT(*)
	FROM forums`

const selectPostsCount = `
	SELECT COUNT(*)
	FROM posts`

const selectThreadsCount = `
	SELECT COUNT(*)
	FROM threads`

const selectUsersCount = `
	SELECT COUNT(*)
	FROM users`

const deleteData = `
	TRUNCATE TABLE forums CASCADE;
	TRUNCATE TABLE posts CASCADE;
	TRUNCATE TABLE threads CASCADE;
	TRUNCATE TABLE users CASCADE;`

func Status(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	counts := &Counts{}
	database.DB.QueryRow(selectForumsCount).Scan(&counts.Forum)
	database.DB.QueryRow(selectPostsCount).Scan(&counts.Post)
	database.DB.QueryRow(selectThreadsCount).Scan(&counts.Thread)
	database.DB.QueryRow(selectUsersCount).Scan(&counts.User)

	jsonCounts, _ := json.Marshal(counts)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonCounts)
}

func Clear(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	database.DB.QueryRow(deleteData).Scan()

	counts := &Counts{}
	database.DB.QueryRow(selectForumsCount).Scan(&counts.Forum)
	database.DB.QueryRow(selectPostsCount).Scan(&counts.Post)
	database.DB.QueryRow(selectThreadsCount).Scan(&counts.Thread)
	database.DB.QueryRow(selectUsersCount).Scan(&counts.User)

	jsonCounts, _ := json.Marshal(counts)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonCounts)
}
