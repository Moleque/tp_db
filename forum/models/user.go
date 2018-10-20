package models

import (
	"encoding/json"
	"net/http"
	"regexp"
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

const createUser = `
	INSERT INTO users (email, nickname, fullname, about)
	VALUES ($1, $2, $3, $4) RETURNING email, nickname, fullname, about`

const selectUser = `
	SELECT email, nickname, fullname, about
	FROM users
	WHERE nickname = $1 OR email = $2`

const selectUserByNickname = `
	SELECT email, nickname, fullname, about
	FROM users
	WHERE nickname = $1`

const updateUser = `
	UPDATE users 
	SET	email = $1, fullname = $2, about = $3
	WHERE nickname = $4
	RETURNING email, nickname, fullname, about`

func UserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	nickname := params.ByName("nickname")
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9_.]*$", nickname); !ok {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	user := &User{}
	if decode(r.Body, user) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := database.DB.QueryRow(createUser, user.Email, nickname, user.Fullname, user.About).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() == "unique_violation" {
			if rows, err := database.DB.Query(selectUser, nickname, user.Email); err == nil {
				defer rows.Close()
				users := []*User{}
				for rows.Next() {
					user := &User{}
					rows.Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
					users = append(users, user)
				}
				jsonUsers, _ := json.Marshal(users)
				w.WriteHeader(http.StatusConflict)
				w.Write(jsonUsers)
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
	nickname := params.ByName("nickname")

	user := &User{}
	if database.DB.QueryRow(selectUserByNickname, nickname).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About) != nil {
		message := Error{"Can't find user with id #42\n"}
		jsonMessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonMessage)
		return
	}

	jsonUser, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

func UserUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	nickname := params.ByName("nickname")

	user := &User{}
	if decode(r.Body, user) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := database.DB.QueryRow(updateUser, user.Email, user.Fullname, user.About, nickname).Scan(&user.Email, &user.Nickname, &user.Fullname, &user.About)
	if _, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	jsonUser, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}
