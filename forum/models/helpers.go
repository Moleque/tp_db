package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/moleque/tp_db/forum/database"
)

func decode(body io.ReadCloser, request interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(request)
	if err != nil {
		return fmt.Errorf("decode error:", err)
	}
	body.Close()
	return nil
}

func isEmpty(str string) interface{} {
	if str == "" {
		return nil
	} else {
		return str
	}
}

func getThreadBySlugId(slugId string) Thread {
	thread := &Thread{}
	database.DB.QueryRow(selectThreadById, slugId).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	if isEmpty(thread.Slug) == nil {
		database.DB.QueryRow(selectThreadBySlug, slugId).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	}
	return *thread
}

func conflict(textMessage string) []byte {
	message := Error{textMessage}
	jsonMessage, _ := json.Marshal(message)
	return jsonMessage
}

func Creator(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	path1 := params.ByName("path1")
	if path1 == "create" {
		ForumCreate(w, r, params)
	} else {
		ThreadCreate(w, r, params)
	}
}
