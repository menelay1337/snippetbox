package main

import (
	"log"
	"net/http"
)

func main() {
	// requests handling part
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on port 4000.\nPress Ctrl+C to stop the server.")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
