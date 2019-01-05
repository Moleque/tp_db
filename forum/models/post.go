package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tp_db/forum/database"

	"github.com/julienschmidt/httprouter"
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

type Details struct {
	Post   Post        `json:"post"`
	User   interface{} `json:"author,omitempty"`
	Forum  interface{} `json:"forum,omitempty"`
	Thread interface{} `json:"thread,omitempty"`
}

const createPost = `
	INSERT INTO posts (created, message, username, forum, thread, parent) 
	VALUES `

const selectMainPost = `
	SELECT id, created, message, username, forum, thread, parent
	FROM posts
	WHERE thread = $1 AND forum = $2 AND parent = 0`

const selectPost = `
	SELECT id, created, message, username, forum, thread, parent, isedited
	FROM posts
	WHERE id = $1`

const selectPosts = `
	SELECT id, created, message, username, forum, thread, parent
	FROM posts`

const updatePost = `
	UPDATE posts
	SET message = $2, isedited = true
	WHERE id = $1
	RETURNING id, created, message, username, forum, thread, parent, isedited`

func PostsCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slugId := params.ByName("slug_or_id")

	thread := getThreadBySlugId(slugId)
	//проверка, что существует такая ветка
	if thread.Slug == "" {
		fmt.Println("NF1")
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find post thread by id:" + slugId))
		return
	}

	posts := []*Post{}
	if decode(r.Body, &posts) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(posts) == 0 {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("[]"))
		return
	}

	query := createPost
	vals := []interface{}{}

	createTime := time.Now() //время для всех постов
	iterator := 0

	database.DB.Exec("SET LOCAL synchronous_commit TO OFF")

	for _, post := range posts {
		shift := iterator * 6
		parentVal := fmt.Sprintf("$%d", shift+6)
		if post.Parent != 0 {
			parentVal = "(SELECT id FROM posts WHERE id = " + parentVal + " and thread = $" + strconv.Itoa(shift+5) + ")"
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, %s),", shift+1, shift+2, shift+3, shift+4, shift+5, parentVal)
		vals = append(vals, createTime, post.Message, post.Author, thread.Forum, thread.Id, post.Parent)
		iterator++
	}
	query = query[0:len(query)-1] + " RETURNING id, created, message, username, forum, thread, parent;"
	template, _ := database.DB.Prepare(query)

	rows, err := template.Query(vals...)
	if err != nil {
		if err.Error() == "pq: нулевое значение в столбце \"parent\" нарушает ограничение NOT NULL" || err.Error() == "pq: null value in column \"parent\" violates not-null constraint" {
			w.WriteHeader(http.StatusConflict)
			w.Write(conflict("Parent post was created in another thread"))
			return
		}
		if err.Error() == "pq: INSERT или UPDATE в таблице \"posts\" нарушает ограничение внешнего ключа \"posts_username_fkey\"" || err.Error() == "pq: insert or update on table \"posts\" violates foreign key constraint \"posts_username_fkey\"" {
			fmt.Println("NF2")
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find post author by nickname:"))
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer rows.Close()

	number := 0
	for rows.Next() {
		post := posts[number]
		rows.Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, &post.Parent)
		number++
	}

	jsonPosts, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonPosts)
}

func PostGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := params.ByName("id")

	postDetails := &Details{}

	post := &Post{}
	database.DB.QueryRow(selectPost, id).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, &post.Parent, &post.IsEdited)
	if post.Author == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find post with id:" + id))
		return
	}
	postDetails.Post = *post

	queryParams := strings.Split(r.URL.Query().Get("related"), ",")
	objects(postDetails, queryParams)

	jsonPost, _ := json.Marshal(postDetails)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPost)
}

func PostUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := params.ByName("id")

	postUpdate := &Post{}
	if decode(r.Body, postUpdate) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	post := &Post{}
	database.DB.QueryRow(selectPost, id).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, &post.Parent, &post.IsEdited)
	if isEmpty(postUpdate.Message) != nil && postUpdate.Message != post.Message {
		database.DB.QueryRow(updatePost, id, postUpdate.Message).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, &post.Parent, &post.IsEdited)
		if isEmpty(post.Author) == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find post with id:" + id))
			return
		}
	}

	jsonPost, _ := json.Marshal(post)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPost)
}

func objects(details *Details, params []string) {
	for _, param := range params {
		switch param {
		case "user":
			user := &User{}
			database.DB.QueryRow(selectUserByNickname, details.Post.Author).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
			details.User = *user
		case "forum":
			forum := &Forum{}
			database.DB.QueryRow(selectForum, details.Post.Forum).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts)
			details.Forum = *forum
		case "thread":
			thread := &Thread{}
			database.DB.QueryRow(selectThreadById, details.Post.Thread).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
			details.Thread = *thread
		}
	}
}
