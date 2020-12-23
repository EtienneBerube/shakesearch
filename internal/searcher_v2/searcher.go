package searcher_v2

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"unicode"
)

type Searcher struct {
	CompleteWorks string
	Tree          *TrieTree
}

func New() Searcher {
	return Searcher{}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.Tree = InitTrieTree()
	s.processTree()
	return nil
}

func (s *Searcher) processTree() {
	var buffer strings.Builder
	for index, char := range s.CompleteWorks {
		 if !unicode.IsLetter(char) && buffer.Len() != 0 {
			s.Tree.insert(buffer.String(), index-buffer.Len())
			buffer.Reset()
		} else if unicode.IsLetter(char) {
			buffer.WriteRune(char)
		} else {
			continue
		}
	}
}

func (s *Searcher) Search(query string) []string {
	var idxs []int
	var results []string

	if strings.ContainsAny(query, " ") {
		// Find first word
		superset_idx := s.Tree.Find(query)

		for _, idx := range superset_idx {
			if strings.Contains(s.CompleteWorks[idx:idx+100], query) {
				idxs = append(idxs, idx)
			}
		}
	} else {
		idxs = s.Tree.Find(query)
	}

	if idxs == nil || idxs[0] == -1{
		return []string{}
	}

	idxs = optimizeIndices(idxs)

	for _, idx := range idxs {
		start, end := correctIndexBounds(s, idx)
		results = append(results, s.CompleteWorks[start:end])
	}
	return results
}

func optimizeIndices(indices []int) []int {
	sort.Ints(indices)

	current := indices[0]

	var optimized []int
	optimized = append(optimized, current)

	for _, val := range indices[1:] {
		if val < current+250 {
			continue
		} else {
			optimized = append(optimized, val)
			current = val
		}
	}

	fmt.Printf("INDICES -- Removed overlapping: Went from %d to %d\n", len(indices), len(optimized))

	return optimized
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
