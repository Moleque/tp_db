package models

import (
	"encoding/json"
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

const countThread = `
	SELECT COUNT(*)
	FROM threads
	WHERE forum = $1`

const selectThread = `
	SELECT id, slug, created, title, message, username, forum, votes
	FROM threads
	WHERE forum = $1`

const selectThreadById = `
	SELECT id, slug, created, title, message, username, forum, votes
	FROM threads
	WHERE id = $1`

const selectThreadBySlug = `
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

	//проверка, что существует такой пользователь
	var nickname, temp string
	if database.DB.QueryRow(selectUserByNickname, thread.Author).Scan(&temp, &nickname, &temp, &temp) != nil {
		message := Error{"Can't find user by nickname:" + thread.Author}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}

	//проверка, что существует такой форум
	var check string
	database.DB.QueryRow(selectForum, forum).Scan(&check, &temp, &temp, &temp, &temp)
	if isEmpty(check) == nil {
		message := Error{"Can't find thread forum by slug:" + forum}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}
	forum = check

	//проверка, существует-ли уже данный thread
	if isEmpty(thread.Slug) != nil {
		database.DB.QueryRow(selectThreadBySlug, thread.Slug).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
		if thread.Id != 0 {
			jsonThread, _ := json.Marshal(thread)
			w.WriteHeader(http.StatusConflict)
			w.Write(jsonThread)
			return
		}
	}

	err := database.DB.QueryRow(createThread, thread.Slug, thread.Created, thread.Title, thread.Message, nickname, forum).Scan(&thread.Id, &thread.Forum, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() == "unique_violation" {
			// if database.DB.QueryRow(selectForum, forum.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) == nil {
			// 	jsonForum, _ := json.Marshal(forum)
			// w.WriteHeader(http.StatusConflict)
			// 	w.Write(jsonForum)
			// return
			// }
		}
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonThread)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slugId := params.ByName("slug_or_id")

	thread := getThreadBySlugId(slugId)
	if isEmpty(thread.Forum) == nil {
		message := Error{"Can't find thread by slug:" + slugId}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}

	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThread)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	thread := getThreadBySlugId(params.ByName("slug_or_id"))

	query := paramsGetPosts(selectPosts, r)

	rows, err := database.DB.Query(query, thread.Id, thread.Forum)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		post := &Post{}
		rows.Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, &post.Parent)
		posts = append(posts, post)
	}

	jsonPosts, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPosts)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func paramsGetThreads(query string, r *http.Request) string {
	since := r.URL.Query().Get("since")
	order := r.URL.Query().Get("desc")
	limit := r.URL.Query().Get("limit")

	if since != "" {
		if order == "true" {
			query += " AND created <= '" + since + "'"
		} else {
			query += " AND created >= '" + since + "'"
		}
	}
	if order == "true" {
		query += "\nORDER BY created DESC"
	} else {
		query += "\nORDER BY created"
	}
	if limit != "" {
		query += "\nLIMIT " + limit
	}
	return query
}

func paramsGetPosts(query string, r *http.Request) string {
	since := r.URL.Query().Get("since")
	limit := r.URL.Query().Get("limit")
	if since != "" {
		query += " AND id > '" + since + "'"
	}
	if limit != "" {
		query += "\nLIMIT " + limit
	}
	return query
}