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

const selectCounts = `
	SELECT * 
	FROM (SELECT COUNT(*) FROM forums) forums,
		(SELECT COUNT(*) FROM posts) posts,
		(SELECT COUNT(*) FROM threads) threads,
		(SELECT COUNT(*) FROM users) users`

const deleteData = `
	TRUNCATE TABLE users, forums, threads, posts, votes, members CASCADE`

func Status(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	counts := &Counts{}
	database.DB.QueryRow(selectCounts).Scan(&counts.Forum, &counts.Post, &counts.Thread, &counts.User)

	jsonCounts, _ := json.Marshal(counts)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonCounts)
}

func Clear(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	database.DB.QueryRow(deleteData).Scan()
	w.WriteHeader(http.StatusOK)
}
