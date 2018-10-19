package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Logger(inner httprouter.Handle, name string) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		start := time.Now()
		inner(w, r, params)
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
