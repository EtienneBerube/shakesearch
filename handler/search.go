package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pulley.com/shakesearch/internal/searcher"
	"strconv"
	"strings"
	"time"
)

func HandleSearch(searcher searcher.Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var start int64
		var end int64
		var results []string

		queries, ok := r.URL.Query()["q"]
		if !ok || len(queries[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		query := strings.ToLower(queries[0])

		work, ok := r.URL.Query()["w"]
		if ok {
			id, err := strconv.Atoi(work[0])
			if err != nil{
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad work id"))
				return
			}
			start = time.Now().UnixNano()
			results = searcher.SearchByWork(query, id)
			end = time.Now().UnixNano()
		}else{
			start = time.Now().UnixNano()
			results = searcher.Search(query)
			end = time.Now().UnixNano()
		}

		fmt.Printf("TOOK: %d ns\n", end-start)

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