package models

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Vote struct {
	Nickname string  `json:"nickname"`
	Voice    float32 `json:"voice"`
}

const createVote = `
	INSERT INTO vote (thread_id, user_id, value)
	VALUES ($1, $2, $3)`

const addVoteToThread = `
`

func ThreadVote(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	thread := getThreadBySlugId(params.ByName("slug_or_id"))

	// 	vote := &Vote{}
	// 	if decode(r.Body, thread) != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 	}

	// 	//создание голоса пользователя
	// 	err := database.DB.QueryRow(createVote, thread.Id, nickname, forum).Scan(&thread.Id, &thread.Forum, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	// 	if err, ok := err.(*pq.Error); ok {
	// 		if err.Code.Name() == "unique_violation" {
	// 			// if database.DB.QueryRow(selectForum, forum.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) == nil {
	// 			// 	jsonForum, _ := json.Marshal(forum)
	// 			// w.WriteHeader(http.StatusConflict)
	// 			// 	w.Write(jsonForum)
	// 			// return
	// 			// }
	// 		}
	// 		w.WriteHeader(http.StatusBadGateway)
	// 		return
	// 	}

	// 	//добавление голоса к ветке
	// 	err := database.DB.QueryRow(createVote, thread.Slug, thread.Created, thread.Title, thread.Message, nickname, forum).Scan(&thread.Id, &thread.Forum, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	// 	if err, ok := err.(*pq.Error); ok {
	// 		if err.Code.Name() == "unique_violation" {
	// 			// if database.DB.QueryRow(selectForum, forum.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts) == nil {
	// 			// 	jsonForum, _ := json.Marshal(forum)
	// 			// w.WriteHeader(http.StatusConflict)
	// 			// 	w.Write(jsonForum)
	// 			// return
	// 			// }
	// 		}
	// 		w.WriteHeader(http.StatusBadGateway)
	// 		return
	// 	}

	// 	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	// 	w.Write(jsonThread)
}
