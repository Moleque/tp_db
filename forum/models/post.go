package models

import (
	"encoding/json"
	"log"
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

type Details struct {
	Post   Post   `json:"post"`
	User   User   `json:"user,omitempty"`
	Forum  Forum  `json:"forum,omitempty"`
	Thread Thread `json:"thread,omitempty"`
}

const createPost = `
	INSERT INTO posts (created, message, username, forum, thread, parent)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created, message, username, forum, thread, parent`

const selectMainPost = `
	SELECT id, created, message, username, forum, thread, parent
	FROM posts
	WHERE thread = $1 AND forum = $2 AND parent = 0`

const selectPost = `
	SELECT id, created, message, username, forum, thread, parent
	FROM posts
	WHERE id = $1`

const selectPosts = `
	SELECT id, created, message, username, forum, thread, parent
	FROM posts
	WHERE thread = $1 AND forum = $2`

//ОПТИМИЗАЦИЯ: триггер
const updatePostsCount = `
	UPDATE forums
	SET posts = posts + 1
	WHERE slug = $1`

func PostsCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slugId := params.ByName("slug_or_id")

	posts := []*Post{}
	if decode(r.Body, &posts) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// var thread string
	// var forum string
	// database.DB.QueryRow(selectThreadById, slug).Scan(&thread, &forum)
	// if isEmpty(forum) == nil {
	// 	database.DB.QueryRow(selectThreadBySlug, slug).Scan(&thread, &forum)
	// }
	thread := getThreadBySlugId(slugId)

	log.Println(thread.Id, thread.Forum)
	createTime := time.Now()
	for _, post := range posts {
		log.Println("parent -", post.Parent)
		if post.Parent == 0 {
			database.DB.QueryRow(selectMainPost, thread.Id, thread.Forum).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread)
			log.Println(post.Thread, post.Id)
			if post.Id != 0 {
				jsonPost, _ := json.Marshal(post)
				w.WriteHeader(http.StatusConflict)
				w.Write(jsonPost)
				return
			}
		}

		err := database.DB.QueryRow(createPost, createTime, post.Message, post.Author, thread.Forum, thread.Id, post.Parent).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, post.Parent)
		if err, ok := err.(*pq.Error); ok {
			log.Println(err.Code.Name())
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
		database.DB.QueryRow(updatePostsCount, thread.Forum).Scan()
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
	database.DB.QueryRow(selectPost, id).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, post.Parent)
	if post.Author == "" {
		message := Error{"Can't find post with id:" + id}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMessage)
		return
	}
	postDetails.Post = *post
	objects(r.URL.Query().Get("user"), postDetails)

	jsonPost, _ := json.Marshal(postDetails)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPost)
}

func PostUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func objects(param string, details *Details) {
	// for _, param := range params {
	switch param {
	case "user":
		user := &User{}
		database.DB.QueryRow(selectUserByNickname, details.Post.Author).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
		details.User = *user
	case "forum":
	case "thread":
	}
	// }
}
