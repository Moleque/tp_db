package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// "github.com/gorilla/mux"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"

	"tp_db/forum/models"
)

//=======================

func decode(body io.ReadCloser, request interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(request)
	if err != nil {
		return fmt.Errorf("decode error:", err)
	}
	body.Close()
	return nil
}

func Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "Welcome to the forum!")
}

func Clear(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

//=======================
const createForum = `
INSERT INTO forums (slug, title, username)
VALUES ($1, $2, $3) RETURNING *`

func ForumCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := &models.Forum{}
	if decode(r.Body, forum) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	DB.Query(createForum, forum.Slug, forum.Title, forum.User)

	jsonForum, err := json.Marshal(forum)
	if err != nil {
		log.Printf("cannot marshal:%s", err)
	}
	w.Write(jsonForum)
	w.WriteHeader(http.StatusOK)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

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
	w.WriteHeader(http.StatusOK)
}

func Status(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadVote(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

//=======================
const createUser = `
	INSERT INTO users (email, nickname, fullname, about)
	VALUES ($1, $2, $3, $4) RETURNING email, nickname, fullname, about`

const selectUser = `
	SELECT email, nickname, fullname, about
	FROM users
	WHERE nickname = $1 AND email = $2`

func UserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	nickname := mux.Vars(r)["nickname"]
	user := &models.User{}
	if decode(r.Body, user) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := DB.QueryRow(createUser, user.Email, nickname, user.Fullname, user.About).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
	if err, ok := err.(*pq.Error); ok {
		log.Println(err.Code.Name())
		if err.Code.Name() == "unique_violation" {
			if DB.QueryRow(selectUser, nickname, user.Email).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About) == nil {
				jsonUser, _ := json.Marshal(user)
				w.WriteHeader(http.StatusConflict)
				w.Write(jsonUser)
				return
			}
		}
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	jsonUser, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonUser)
}

func UserGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UserUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
