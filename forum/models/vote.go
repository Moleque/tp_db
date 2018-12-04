package models

import (
	"encoding/json"
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
	VALUES ($1, $2, $3) RETURNING username, value`

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

	thread := getThreadBySlugId(slugId)
	if isEmpty(thread.Author) == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(conflict("Can't find thread by slug:" + slugId))
		return
	}

	vote := &Vote{}
	if decode(r.Body, vote) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	//проверка наличия голоса на данный момент
	var voice int32
	database.DB.QueryRow(selectVoice, thread.Id, vote.Nickname).Scan(&voice)

	if voice == 0 {
		//создание голоса пользователя
		if err := database.DB.QueryRow(createVote, thread.Id, vote.Nickname, vote.Voice).Scan(&vote.Nickname, &vote.Voice); err != nil {
			if err.Error() == "pq: insert or update on table \"votes\" violates foreign key constraint \"votes_username_fkey\"" {
				w.WriteHeader(http.StatusNotFound)
				w.Write(conflict("Can't find post author by nickname:"))
				return
			}
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	} else {
		value := voice
		if voice == 1 && vote.Voice == -1 {
			value = -2
		}
		if voice == -1 && vote.Voice == 1 {
			value = 1
		}
		database.DB.QueryRow(updateVoice, thread.Id, vote.Nickname, value).Scan(&vote.Voice)
	}

	database.DB.QueryRow(selectThreadById, thread.Id).Scan(&thread.Id, &thread.Slug, &thread.Created, &thread.Title, &thread.Message, &thread.Author, &thread.Forum, &thread.Votes)

	jsonThread, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonThread)
}
