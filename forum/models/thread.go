package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"

	"tp_db/forum/database"
)

type Thread struct {
	Id      int32     `json:"id,omitempty"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum,omitempty"`
	Message string    `json:"message"`
	Votes   int32     `json:"votes,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

const createThread = `
	INSERT INTO threads (slug, created, title, message, username, forum)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, slug, created, title, message, username, forum, votes`

const selectThread = `
	SELECT id, slug, created, title, message, username, forum, votes
	FROM threads
	WHERE slug = $1`

func ThreadCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := params.ByName("path1")

	thread := &Thread{}
	if decode(r.Body, thread) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var nickname, temp string
	if database.DB.QueryRow(selectUserByNickname, thread.Author).Scan(&temp, &nickname, &temp, &temp) != nil {
		message := Error{"Can't find user by nickname:" + thread.Author}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}

	err := database.DB.QueryRow(createThread, thread.Slug, thread.Created, thread.Title, thread.Message, nickname, forum).Scan(&thread.Id, &thread.Forum, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	log.Println(nickname, *thread)
	if err, ok := err.(*pq.Error); ok {
		log.Println(err.Code.Name())
		if err.Code.Name() == "unique_violation" {
			// if database.DB.QueryRow(selectForum, forum.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) == nil {
			// 	jsonForum, _ := json.Marshal(forum)
			w.WriteHeader(http.StatusConflict)
			// 	w.Write(jsonForum)
			return
			// }
		}
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonThread)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slug := params.ByName("slug")

	thread := &Thread{}
	if database.DB.QueryRow(selectThread, slug).Scan(&thread.Id, &thread.Forum, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes) != nil {
		message := Error{"Can't find forum by slug:" + slug}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}

	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThread)
}
