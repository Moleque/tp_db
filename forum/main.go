package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Docker test!")
	})

	log.Println("Server started!")
	http.ListenAndServe(":5000", mux)
}
