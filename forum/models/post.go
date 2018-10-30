package models

import (
	"encoding/json"
	"net/http"
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
	Post Post `json:"post"`
	// User   User   `json:"user,omitempty"`
	// Forum  Forum  `json:"forum,omitempty"`
	// Thread Thread `json:"thread,omitempty"`
}

const createPost = `
	INSERT INTO posts (created, message, username, forum, thread, parent)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created, message, username, forum, thread, parent`

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
	FROM posts
	WHERE thread = $1 AND forum = $2`

const updatePost = `
	UPDATE posts
	SET message = $2, isedited = true
	WHERE id = $1
	RETURNING id, created, message, username, forum, thread, parent, isedited`

//ОПТИМИЗАЦИЯ: триггер
const updatePostsCount = `
	UPDATE forums
	SET posts = posts + 1
	WHERE slug = $1`

func PostsCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slugId := params.ByName("slug_or_id")

	thread := getThreadBySlugId(slugId)
	//проверка, что существует такая ветка
	if thread.Slug == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find post thread by id:" + slugId))
		return
	}

	posts := []*Post{}
	if decode(r.Body, &posts) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createTime := time.Now() //время для всех постов
	for _, post := range posts {
		if post.Parent != 0 { //если родитель есть
			parent := &Post{} //получаем родительский пост
			database.DB.QueryRow(selectPost, post.Parent).Scan(&parent.Id, &parent.Created, &parent.Message, &parent.Author, &parent.Forum, &parent.Thread, &parent.Parent, &parent.IsEdited)
			if thread.Id != parent.Thread || isEmpty(parent.Author) == nil {
				w.WriteHeader(http.StatusConflict)
				w.Write(conflict("Parent post was created in another thread"))
				return
			}
		}

		//проверка, что существует такой пользователь
		var nickname, temp string
		if database.DB.QueryRow(selectUserByNickname, post.Author).Scan(&temp, &nickname, &temp, &temp) != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find post author by nickname:" + nickname))
			return
		}

		database.DB.QueryRow(createPost, createTime, post.Message, post.Author, thread.Forum, thread.Id, post.Parent).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, &post.Parent)
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
	database.DB.QueryRow(selectPost, id).Scan(&post.Id, &post.Created, &post.Message, &post.Author, &post.Forum, &post.Thread, &post.Parent, &post.IsEdited)
	if post.Author == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find post with id:" + id))
		return
	}
	postDetails.Post = *post

	var queryParams []string
	queryParams = append(queryParams, r.URL.Query().Get("user"))
	// queryParams = append(queryParams, r.URL.Query().Get("forum"))
	// queryParams = append(queryParams, r.URL.Query().Get("thread"))
	objects(queryParams, postDetails)

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

func objects(params []string, details *Details) {
	// for _, param := range params {
	// 	// switch param {
	// 	// case "user":
	// 	// 	user := &User{}
	// 	// 	database.DB.QueryRow(selectUserByNickname, details.Post.Author).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
	// 	// 	details.User = *user
	// 	// case "forum":
	// 	// 	forum := &Forum{}
	// 	// 	// database.DB.QueryRow(selectUserByNickname, details.Post.Author).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
	// 	// 	details.Forum = *forum
	// 	// case "thread":
	// 	// 	thread := &Thread{}
	// 	// 	// database.DB.QueryRow(selectUserByNickname, details.Post.Author).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
	// 	// 	details.Thread = *thread
	// 	// }
	// }
}
