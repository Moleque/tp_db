package models

import (
	"encoding/json"
	"net/http"
	"time"
	"tp_db/forum/database"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
)

type Post struct {
	Id       int32     `json:"id,omitempty"`
	Parent   float32   `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int32     `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
}

const createPost = `
	INSERT INTO posts (message, username, forum, thread)
	VALUES ($1, $2, $3, $4) RETURNING id, created, message, username, forum, thread`

func PostGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostsCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slug := params.ByName("slug_or_id")

	posts := []*Post{}
	if decode(r.Body, &posts) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var thread, forum string
	database.DB.QueryRow(selectThreadById, slug).Scan(&thread, &forum)
	if isEmpty(forum) == nil {
		database.DB.QueryRow(selectThreadBySlug, slug).Scan(&thread, &forum)
	}

	for _, post := range posts {
		err := database.DB.QueryRow(createPost, post.Message, post.Author, forum, thread).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				// if rows, err := database.DB.Query(selectUser, nickname, user.Email); err == nil {
				// 	defer rows.Close()
				// 	users := []*User{}
				// 	for rows.Next() {
				// 		user := &User{}
				// 		rows.Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
				// 		users = append(users, user)
				// 	}
				// jsonUsers, _ := json.Marshal(users)
				w.WriteHeader(http.StatusConflict)
				// w.Write(jsonUsers)
				return
				// }
			}
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	jsonPosts, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonPosts)
}
