package models

import (
	"encoding/json"
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
VALUES ($1, $2, $3) RETURNING slug, title, username, threads, posts`

// func UserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	nickname := params.ByName("nickname")
// 	if ok, _ := regexp.MatchString("^[A-z0-9_.]*$", nickname); !ok {
// 		w.WriteHeader(http.StatusBadGateway)
// 		return
// 	}

// 	user := &User{}
// 	if decode(r.Body, user) != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}

// 	err := database.DB.QueryRow(createUser, user.Email, nickname, user.Fullname, user.About).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About, &)
// 	if err, ok := err.(*pq.Error); ok {
// 		if err.Code.Name() == "unique_violation" {
// 			if rows, err := database.DB.Query(selectUser, nickname, user.Email); err == nil {
// 				defer rows.Close()
// 				users := []*User{}
// 				for rows.Next() {
// 					user := &User{}
// 					rows.Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
// 					users = append(users, user)
// 				}
// 				jsonUsers, _ := json.Marshal(users)
// 				w.WriteHeader(http.StatusConflict)
// 				w.Write(jsonUsers)
// 				return
// 			}
// 		}
// 		w.WriteHeader(http.StatusBadGateway)
// 		return
// 	}
// 	jsonUser, _ := json.Marshal(user)
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(jsonUser)
// }

func ForumCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	forum := &Forum{}
	if decode(r.Body, forum) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// err :=
	database.DB.QueryRow(createForum, forum.Slug, forum.Title, forum.User).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts)

	jsonForum, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonForum)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// func ForumDetails(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	log.Println("was queried0")
// 	rows, err := DB.Query("Select nickname, fullname, about, email From users")
// 	log.Println("was queried1")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	log.Println("was queried")

// 	for rows.Next() {
// 		user := &User{}
// 		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Println(user.Email)
// 	}
// }
