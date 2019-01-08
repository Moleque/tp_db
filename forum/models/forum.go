package models

import (
	"encoding/json"
	"fmt"
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
	VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3))
	RETURNING slug, title, username, threads, posts`

const selectForum = `
	SELECT slug, title, username, threads, posts
	FROM forums
	WHERE slug = $1`

const selectForumUsers = `
	SELECT email, nickname, fullname, about
	FROM members JOIN users ON members.username = users.nickname AND members.forum = $1`

func ForumCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	forum := &Forum{}
	if decode(r.Body, forum) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := database.DB.QueryRow(createForum, forum.Slug, forum.Title, forum.User).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts)
	if err != nil {
		fmt.Println(err)
		if err.Error() == "pq: нулевое значение в столбце \"username\" нарушает ограничение NOT NULL" || err.Error() == "pq: null value in column \"username\" violates not-null constraint" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find user by nickname:" + forum.User))
			return
		}
		// ОПТИМИЗАЦИЯ: при конфликте возвращать существующую запись
		if err.Error() == "pq: повторяющееся значение ключа нарушает ограничение уникальности \"forums_slug_key\"" || err.Error() == "pq: duplicate key value violates unique constraint \"forums_slug_key\"" {
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
	slug := params.ByName("slug")

	query := paramsGetThreads(selectThread, r)
	rows, err := database.DB.Query(query, slug)
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

	if len(threads) == 0 {
		//проверка,что форум существует
		forum := &Forum{}
		if database.DB.QueryRow(selectForum, slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find forum by slug:" + slug))
			return
		}
	}

	jsonThreads, _ := json.Marshal(threads)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThreads)
	return
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slug := params.ByName("slug")

	query := paramsGetUsers(selectForumUsers, r)
	fmt.Println(query)
	rows, err := database.DB.Query(query, slug)
	if err != nil {
		fmt.Println(err)
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

	if len(users) == 0 {
		//проверка,что форум существует
		forum := &Forum{}
		if database.DB.QueryRow(selectForum, slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find forum by slug:" + slug))
			return
		}
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
			query += " AND users.nickname < '" + since + "'"
		} else {
			query += " AND users.nickname > '" + since + "'"
		}
	}
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
