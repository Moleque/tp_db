package models

import (
	"encoding/json"
	"net/http"
	"tp_db/forum/database"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
)

type User struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}

// func GetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	rows, err := DB.Query("Select nickname, fullname, about, email From users")
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

const createUser = `
	INSERT INTO users (email, nickname, fullname, about)
	VALUES ($1, $2, $3, $4) RETURNING email, nickname, fullname, about`

const selectUser = `
	SELECT email, nickname, fullname, about
	FROM users
	WHERE nickname = $1 AND email = $2`

func UserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	nickname := params.ByName("nickname")
	user := &User{}
	if decode(r.Body, user) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := database.DB.QueryRow(createUser, user.Email, nickname, user.Fullname, user.About).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() == "unique_violation" {
			if database.DB.QueryRow(selectUser, nickname, user.Email).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About) == nil {
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
