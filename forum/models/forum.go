package models

import (
	"encoding/json"
	"log"
	"net/http"
	"tp_db/forum/database"

	"github.com/julienschmidt/httprouter"
)

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts"`
	Threads int32  `json:"threads"`
}

const createForum = `
INSERT INTO forums (slug, title, username)
VALUES ($1, $2, $3) RETURNING *`

func ForumCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := &Forum{}
	if decode(r.Body, forum) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	database.DB.Query(createForum, forum.Slug, forum.Title, forum.User)

	jsonForum, err := json.Marshal(forum)
	if err != nil {
		log.Printf("cannot marshal:%s", err)
	}
	w.Write(jsonForum)
	w.WriteHeader(http.StatusOK)
}

// func ForumDetails(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	log.Println("was queried0")
// 	rows, err := DB.Query("Select nickname, fullname, about, email From users")
// 	log.Println("was queried1")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	log.Println("was queried")

// 	for rows.Next() {
// 		user := &User{}
// 		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Println(user.Email)
// 	}
// }
