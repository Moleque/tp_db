package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tp_db/forum/database"

	"github.com/julienschmidt/httprouter"
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
}

const createVote = `
	INSERT INTO votes (thread_id, username, value)
	VALUES (%s, $1, $2) 
	ON CONFLICT (thread_id, username) DO UPDATE SET value = $2
	RETURNING username, value`

const selectVoice = `
	SELECT value
	FROM votes
	WHERE thread_id = $1 AND username = $2`

const updateVoice = `
	UPDATE votes
	SET value = value + $3
	WHERE thread_id = $1 AND username = $2
	RETURNING value`

const updateVotes = `
	UPDATE threads
	SET votes = votes + $2
	WHERE id = $1`

const addVoteToThread = `
	UPDATE threads 
	SET votes = votes + $2
	WHERE id = $1
	RETURNING id, slug, created, title, message, username, forum, votes`

func ThreadVote(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	slugId := params.ByName("slug_or_id")

	threadIdQuery := getThreadIdQuery(slugId)

	vote := &Vote{}
	if decode(r.Body, vote) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	query := fmt.Sprintf(createVote, threadIdQuery)
	if err := database.DB.QueryRow(query, vote.Nickname, vote.Voice).Scan(&vote.Nickname, &vote.Voice); err != nil {
		fmt.Println(err)
		if err.Error() == "pq: INSERT или UPDATE в таблице \"votes\" нарушает ограничение внешнего ключа \"votes_username_fkey\"" || err.Error() == "pq: insert or update on table \"votes\" violates foreign key constraint \"votes_username_fkey\"" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find post author by nickname:"))
			return
		}
		if err.Error() == "pq: INSERT или UPDATE в таблице \"votes\" нарушает ограничение внешнего ключа \"votes_thread_id_fkey\"" || err.Error() == "pq: insert or update on table \"votes\" violates foreign key constraint \"votes_thread_id_fkey\"" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find thread by id:" + slugId))
			return
		}
		if err.Error() == "pq: нулевое значение в столбце \"thread_id\" нарушает ограничение NOT NULL" || err.Error() == "pq: null value in column \"thread_id\" violates not-null constraint" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(conflict("Can't find thread by slug:" + slugId))
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	thread := getThreadBySlugId(slugId)

	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThread)
}
