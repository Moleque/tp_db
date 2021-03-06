package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

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
	VALUES ($1, $2, $3, $4, $5, (SELECT slug FROM forums WHERE slug = $6))
	RETURNING id, slug, created, title, message, username, forum, votes`

const countThread = `
	SELECT COUNT(*)
	FROM threads
	WHERE forum = $1`

const selectThread = `
	SELECT *
	FROM threads
	WHERE forum = $1`

const selectThreadById = `
	SELECT *
	FROM threads
	WHERE id = $1`

const selectThreadBySlug = `
	SELECT *
	FROM threads
	WHERE slug = $1`

const updateThread = `
	UPDATE threads
	SET title = COALESCE($1, title), message = COALESCE($2, message)
	WHERE id = %s
	RETURNING id, slug, created, title, message, username, forum, votes`

func ThreadCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := params.ByName("path1")

	thread := &Thread{}
	if decode(r.Body, thread) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var slug interface{}
	if isEmpty(thread.Slug) != nil {
		slug = thread.Slug
	} else {
		slug = nil
	}

	err := database.DB.QueryRow(createThread, slug, thread.Created, thread.Title, thread.Message, thread.Author, forum).Scan(&thread.Id, &thread.Forum, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	if err != nil {
		if err.Error() == "pq: insert or update on table \"threads\" violates foreign key constraint \"threads_username_fkey\"" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find user by nickname:"))
			return
		}
		if err.Error() == "pq: INSERT или UPDATE в таблице \"threads\" нарушает ограничение внешнего ключа \"threads_username_fkey\"" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find thread author by nickname:" + thread.Author))
			return
		}
		//проверка, что существует такой форум
		if err.Error() == "pq: нулевое значение в столбце \"forum\" нарушает ограничение NOT NULL" || err.Error() == "pq: null value in column \"forum\" violates not-null constraint" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find thread forum by slug:" + thread.Slug))
			return
		}
		//проверка, существует-ли уже данный thread
		if err.Error() == "pq: повторяющееся значение ключа нарушает ограничение уникальности \"threads_slug_key\"" || err.Error() == "pq: duplicate key value violates unique constraint \"threads_slug_key\"" {
			database.DB.QueryRow(selectThreadBySlug, thread.Slug).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
			jsonThread, _ := json.Marshal(thread)
			w.WriteHeader(http.StatusConflict)
			w.Write(jsonThread)
			return
		}
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
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find thread by slug:" + slugId))
		return
	}

	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThread)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	thread := getThreadBySlugId(params.ByName("slug_or_id"))
	if isEmpty(thread.Author) == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find thread by slug:" + thread.Slug))
		return
	}

	query := paramsGetPosts(selectPosts, r)

	rows, err := database.DB.Query(query, thread.Id)
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
	slugId := params.ByName("slug_or_id")

	querySlugId := getThreadIdQuery(slugId)

	thread := &Thread{}
	if decode(r.Body, thread) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	query := fmt.Sprintf(updateThread, querySlugId)
	if err := database.DB.QueryRow(query, isEmpty(thread.Title), isEmpty(thread.Message)).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find thread by slug:" + slugId))
		return
	}
	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThread)
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
	order := r.URL.Query().Get("desc")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	switch sort {
	case "tree":
		query += "\nWHERE thread = $1"
		if since != "" {
			if order == "true" {
				query += " AND path < (SELECT path FROM posts WHERE id = " + since + ")"
			} else {
				query += " AND path > (SELECT path FROM posts WHERE id = " + since + ")"
			}
		}
		if order == "true" {
			query += "\nORDER BY path DESC"
		} else {
			query += "\nORDER BY path ASC"
		}
		if limit != "" {
			query += "\nLIMIT " + limit
		}
	case "parent_tree":
		query += "\nWHERE root IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0"
		if since != "" {
			if order == "true" {
				query += " AND id < (SELECT root FROM posts WHERE id = " + since + ")"
			} else {
				query += " AND id > (SELECT root FROM posts WHERE id = " + since + ")"
			}
		}
		if limit != "" {
			limit = "\nLIMIT " + limit
		}
		if order == "true" {
			query += "\nORDER BY id DESC" + limit + ") ORDER BY root DESC, path"
		} else {
			query += "\nORDER BY id ASC" + limit + ") ORDER BY path"
		}
	default:
		query += "\nWHERE thread = $1"
		if since != "" {
			if order == "true" {
				query += " AND id < '" + since + "'"
			} else {
				query += " AND id > '" + since + "'"
			}
		}
		if order == "true" {
			query += "\nORDER BY id DESC"
		} else {
			query += "\nORDER BY id ASC"
		}
		if limit != "" {
			query += "\nLIMIT " + limit
		}
	}
	return query
}
