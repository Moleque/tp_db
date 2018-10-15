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

	log.Printf("Try to start http server http://127.0.0.1:5000")
	if err := http.ListenAndServe(":5000", mux); err != nil {
		log.Fatalf("listening error:%s", err)
	}
}
