package main

import (
	"fmt"
	"log"
	"net/http"
	"tp_db/forum/db"

	"github.com/julienschmidt/httprouter"
)

const (
	DB_USER     = "docker"
	DB_PASSWORD = "docker"
	DB_NAME     = "forum"
)

var DB = &db.DataBase{}

type User struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// func ForumCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	fmt.Fprintf(w, "test, %s\n", params.ByName("value"))
// }

// func ThreadCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	fmt.Fprintf(w, "test, %s\n", params.ByName("value"))
// }

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

func ThreadsList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "Thread list for: %s\n", params.ByName("slug"))
}

// func UsersList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	fmt.Fprintf(w, "test, %s\n", params.ByName("value"))
// }

// func ThreadInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	fmt.Fprintf(w, "test, %s\n", params.ByName("value"))
// }

func main() {
	router := httprouter.New()

	router.GET("/", Index)
	// router.POST("/forum/create", ForumCreate)
	// router.POST("/forum/:slug/create", ThreadCreate)
	router.GET("/forum/:slug/details", ForumDetails)
	router.GET("/forum/:slug/threads", ThreadsList)
	// router.GET("/forum/:slug/users", UsersList)
	// router.GET("/post/:id/details", ThreadInfo)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	DB.Connect(dsn)
	defer DB.Disconnect()

	log.Printf("try to start http server http://127.0.0.1:5000")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatalf("listening error:%s", err)
	}
}
