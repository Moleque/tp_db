package controllers

import (
	"fmt"
	"net/http"

	// "github.com/gorilla/mux"

	"github.com/julienschmidt/httprouter"
)

//=======================

func Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "Welcome to the forum!")
}

func Clear(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

//=======================
