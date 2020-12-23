package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pulley.com/shakesearch/handler"
	"pulley.com/shakesearch/internal/searcher"
)

func main() {

	searcher := searcher.New()
	err := searcher.Load("completeworks.txt")

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handler.HandleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}