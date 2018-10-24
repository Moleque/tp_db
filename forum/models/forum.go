package models

import (
	"encoding/json"
	"net/http"
	"tp_db/forum/database"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
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
	VALUES ($1, $2, $3) RETURNING slug, title, username, threads, posts`

const selectForum = `
	SELECT slug, title, username, threads, posts
	FROM forums
	WHERE slug = $1`

func ForumCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := &Forum{}
	if decode(r.Body, forum) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var nickname, temp string
	if database.DB.QueryRow(selectUserByNickname, forum.User).Scan(&temp, &nickname, &temp, &temp) != nil {
		message := Error{"Can't find user by nickname:" + forum.User}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}

	err := database.DB.QueryRow(createForum, forum.Slug, forum.Title, nickname).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() == "unique_violation" {
			if database.DB.QueryRow(selectForum, forum.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) == nil {
				jsonForum, _ := json.Marshal(forum)
				w.WriteHeader(http.StatusConflict)
				w.Write(jsonForum)
				return
			}
		}
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	jsonForum, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonForum)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slug := params.ByName("slug")

	forum := &Forum{}
	if database.DB.QueryRow(selectForum, slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) != nil {
		message := Error{"Can't find user by nickname:" + slug}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}

	jsonForum, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonForum)
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

func ForumGetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
