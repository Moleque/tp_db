package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the forum!")
}

func Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

//=======================
const createForum = `
INSERT INTO forums (slug, title, username)
VALUES ($1, $2, $3) RETURNING *`

func ForumCreate(w http.ResponseWriter, r *http.Request) {
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

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

//=======================
const createUser = `
	INSERT INTO users (email, nickname, fullname, about)
	VALUES ($1, $2, $3, $4) RETURNING *`

const selectUserByNickname = `
	SELECT email, nickname, fullname, about
	FROM users
	WHERE nickname = $1`

func UserCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	user := &models.User{}
	if decode(r.Body, user) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	rows, err := DB.Query(createUser, user.Email, mux.Vars(r)["nickname"], user.Fullname, user.About)
	// rows.Next()
	rows.Scan()
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() == "unique_violation" {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("cannot marshal:%s", err)
	}
	w.Write(jsonUser)
	w.WriteHeader(http.StatusCreated)
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
