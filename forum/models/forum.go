package models

import (
	"encoding/json"
	"log"
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

const selectForumUsers = `
	SELECT DISTINCT *
	FROM (SELECT email, nickname, fullname, about
		FROM forums JOIN threads ON (forums.slug = threads.forum) 
		JOIN users ON (users.nickname = threads.username)
		WHERE forums.slug = $1 
		UNION
		SELECT email, nickname, fullname, about
		FROM forums JOIN posts ON (forums.slug = posts.forum) 
		JOIN users ON (users.nickname = posts.username)
		WHERE forums.slug = $1) AS users`

// Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.

// Пользователи выводятся отсортированные по nickname в порядке возрастания.
// Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.

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
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find user by nickname:" + slug))
		return
	}

	jsonForum, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonForum)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := params.ByName("slug")

	var count int
	database.DB.QueryRow(countThread, forum).Scan(&count)
	if count == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find forum by slug:" + forum))
		return
	}

	query := paramsGetThreads(selectThread, r)

	rows, err := database.DB.Query(query, forum)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	defer rows.Close()
	threads := []*Thread{}
	for rows.Next() {
		thread := &Thread{}
		rows.Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
		threads = append(threads, thread)
	}
	jsonThreads, _ := json.Marshal(threads)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThreads)
	return
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slug := params.ByName("slug")

	forum := &Forum{}
	database.DB.QueryRow(selectForum, slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts)
	//проверка,что форум существует
	if isEmpty(forum.Title) == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find forum by slug:" + slug))
		return
	}

	query := paramsGetUsers(selectForumUsers, r)
	rows, err := database.DB.Query(query, slug)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := []*User{}
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
		users = append(users, user)
	}

	jsonUsers, _ := json.Marshal(users)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUsers)
}

func paramsGetUsers(query string, r *http.Request) string {
	since := r.URL.Query().Get("since")
	order := r.URL.Query().Get("desc")
	limit := r.URL.Query().Get("limit")

	if since != "" {
		if order == "true" {
			query += "\nWHERE users.nickname < '" + since + "'"
		} else {
			query += "\nWHERE users.nickname > '" + since + "'"
		}
	}
	// query += "\nGROUP BY users.nickname, users.email"
	if order == "true" {
		query += "\nORDER BY users.nickname DESC"
	} else {
		query += "\nORDER BY users.nickname ASC"
	}
	if limit != "" {
		query += "\nLIMIT " + limit
	}
	return query
}
