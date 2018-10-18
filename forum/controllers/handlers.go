package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Moleque/tp_db/forum/models"
)

// func decode(body, request interface{}) {
// 	decoder := json.NewDecoder(body)
// 	err := decoder.Decode(request)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	body.Close()
// }

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the forum!")
}

func Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("forum")
	fmt.Fprintf(w, "Hi forum!")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := &models.Forum{}
	// decode(r.Body, forum)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(forum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	log.Println("sdaf")
	forum.Posts = 0
	forum.Threads = 0

	log.Println(forum.)

	jsonForum, err := json.Marshal(&forum)
	if err != nil {
		log.Printf("cannot marshal:%s", err)
	}
	w.Write(jsonForum)
	// if _, exists := users[request.Email]; exists == false {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("No have this user"))
	// 	return
	// }
	// if users[request.Email].Password != request.Password {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("Wrong password"))
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "Hi forum!")
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

func UserCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
