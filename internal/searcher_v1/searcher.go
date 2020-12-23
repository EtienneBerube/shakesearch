package searcher

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
)

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *Searcher) Search(query string) []string {
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	results := []string{}
	for _, idx := range idxs {
		start, end := correctIndexBounds(s, idx)
		results = append(results, s.CompleteWorks[start:end])
	}
	return results
}

func correctIndexBounds(s *Searcher, idx int) (start int, end int) {

	if idx < 250 {
		start = 0
	} else {
		start = idx - 250
	}

	if idx > len(s.CompleteWorks)-250 {
		end = len(s.CompleteWorks) - 1
	} else {
		end = idx + 250
	}

	return start, end
}


