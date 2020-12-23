package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pulley.com/shakesearch/internal/searcher_v1"
	"pulley.com/shakesearch/internal/searcher_v2"
	"time"
)

func HandleSearch(searcher searcher_v1.Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		start := time.Now().UnixNano()
		results := searcher.Search(query[0])
		end := time.Now().UnixNano()
		fmt.Printf("Start: %d, End: %d\n", start, end)
		fmt.Printf("DEFAULT TOOK: %d ns\n", end-start)

		resp := make(map[string]interface{})
		resp["results"] = results
		resp["time"] = fmt.Sprintf("%d ns", end-start)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func HandleSearchNew(searcher searcher_v2.Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		//if len(query[0]) < 3 {
		//	w.WriteHeader(http.StatusBadRequest)
		//	w.Write([]byte("search query insufficient length"))
		//}

		if len(query[0]) > 100 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("search query too long"))
		}

		start := time.Now().UnixNano()
		results := searcher.Search(query[0])
		end := time.Now().UnixNano()

		fmt.Printf("NEW TOOK: %d ns\n", end-start)

		resp := make(map[string]interface{})
		resp["results"] = results
		resp["time"] = fmt.Sprintf("%d ns", end-start)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(resp)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}
