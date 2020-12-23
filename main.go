package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pulley.com/shakesearch/handler"
	"pulley.com/shakesearch/internal/searcher_v1"
	"pulley.com/shakesearch/internal/searcher_v2"
)

func main() {
	searcherV2 := searcher_v2.New()
	searcherV1 := searcher_v1.Searcher{}
	err := searcherV2.Load("completeworks.txt")

	if err != nil {
		log.Fatal(err)
	}

	err = searcherV1.Load("completeworks.txt")

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search/default", handler.HandleSearch(searcherV1))
	http.HandleFunc("/search/new", handler.HandleSearchNew(searcherV2))

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