package models

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}

func GetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	rows, err := DB.Query("Select nickname, fullname, about, email From users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	log.Println("was queried")

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(user.Email)
	}
}

func ForumDetails(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("was queried0")
	rows, err := DB.Query("Select nickname, fullname, about, email From users")
	log.Println("was queried1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	log.Println("was queried")

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(user.Email)
	}
}
