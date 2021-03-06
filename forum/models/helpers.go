package models

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"tp_db/forum/database"
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

func conflict(textMessage string) []byte {
	message := Error{textMessage}
	jsonMessage, _ := json.Marshal(message)
	return jsonMessage
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
	_, err := strconv.Atoi(slugId)
	if err == nil {
		database.DB.QueryRow(selectThreadById, slugId).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	} else {
		database.DB.QueryRow(selectThreadBySlug, slugId).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)
	}
	return *thread
}

func getThreadIdQuery(slugId string) string {
	var query string
	_, err := strconv.Atoi(slugId)
	if err == nil {
		query = slugId
	} else {
		query = "(SELECT id FROM threads WHERE slug = '" + slugId + "')"
	}
	return query
}
