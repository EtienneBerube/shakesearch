package searcher

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"sort"
	"strings"
)

var WORKS = [...]string{"THE SONNETS", "ALL’S WELL THAT ENDS WELL", "ANTONY AND CLEOPATRA", "AS YOU LIKE IT", "THE COMEDY OF ERRORS", "THE TRAGEDY OF CORIOLANUS", "CYMBELINE", "THE TRAGEDY OF HAMLET, PRINCE OF DENMARK", "THE FIRST PART OF KING HENRY THE FOURTH", "THE SECOND PART OF KING HENRY THE FOURTH", "THE LIFE OF KING HENRY V", "THE FIRST PART OF HENRY THE SIXTH", "THE SECOND PART OF KING HENRY THE SIXTH", "THE THIRD PART OF KING HENRY THE SIXTH", "KING HENRY THE EIGHTH", "KING JOHN", "THE TRAGEDY OF JULIUS CAESAR", "THE TRAGEDY OF KING LEAR", "LOVE’S LABOUR’S LOST", "MACBETH", "MEASURE FOR MEASURE", "THE MERCHANT OF VENICE", "THE MERRY WIVES OF WINDSOR", "A MIDSUMMER NIGHT’S DREAM", "MUCH ADO ABOUT NOTHING", "OTHELLO", "PERICLES, PRINCE OF TYRE", "KING RICHARD THE SECOND", "KING RICHARD THE THIRD", "THE TRAGEDY OF ROMEO AND JULIET", "THE TAMING OF THE SHREW", "THE TEMPEST", "THE LIFE OF TIMON OF ATHENS", "THE TRAGEDY OF TITUS ANDRONICUS", "THE HISTORY OF TROILUS AND CRESSIDA", "TWELFTH NIGHT", "THE TWO GENTLEMEN OF VERONA", "THE TWO NOBLE KINSMEN", "THE WINTER’S TALE", "A LOVER’S COMPLAINT", "THE PASSIONATE PILGRIM", "THE PHOENIX AND THE TURTLE", "THE RAPE OF LUCRECE", "VENUS AND ADONIS"}
/* ID to WORK
[
	0:  "THE SONNETS",
	1:  "ALL’S WELL THAT ENDS WELL",
	2:  "ANTONY AND CLEOPATRA",
	3:  "AS YOU LIKE IT",
	4:  "THE COMEDY OF ERRORS",
	5:  "THE TRAGEDY OF CORIOLANUS",
	6:  "CYMBELINE",
	7:  "THE TRAGEDY OF HAMLET, PRINCE OF DENMARK",
	8:  "THE FIRST PART OF KING HENRY THE FOURTH",
	9: "THE SECOND PART OF KING HENRY THE FOURTH",
	10: "THE LIFE OF KING HENRY V",
	11: "THE FIRST PART OF HENRY THE SIXTH",
	12: "THE SECOND PART OF KING HENRY THE SIXTH",
	13: "THE THIRD PART OF KING HENRY THE SIXTH",
	14: "KING HENRY THE EIGHTH",
	15: "KING JOHN",
	16: "THE TRAGEDY OF JULIUS CAESAR",
	17: "THE TRAGEDY OF KING LEAR",
	18: "LOVE’S LABOUR’S LOST",
	19: "MACBETH",
	20: "MEASURE FOR MEASURE",
	21: "THE MERCHANT OF VENICE",
	22: "THE MERRY WIVES OF WINDSOR",
	23: "A MIDSUMMER NIGHT’S DREAM",
	24: "MUCH ADO ABOUT NOTHING",
	25: "OTHELLO",
	26: "PERICLES, PRINCE OF TYRE",
	27: "KING RICHARD THE SECOND",
	28: "KING RICHARD THE THIRD",
	29: "THE TRAGEDY OF ROMEO AND JULIET",
	30: "THE TAMING OF THE SHREW",
	31: "THE TEMPEST",
	32: "THE LIFE OF TIMON OF ATHENS",
	33: "THE TRAGEDY OF TITUS ANDRONICUS",
	34: "THE HISTORY OF TROILUS AND CRESSIDA",
	35: "TWELFTH NIGHT",
	36: "THE TWO GENTLEMEN OF VERONA",
	37: "THE TWO NOBLE KINSMEN",
	38: "THE WINTER’S TALE",
	39: "A LOVER’S COMPLAINT",
	40: "THE PASSIONATE PILGRIM",
	41: "THE PHOENIX AND THE TURTLE",
	42: "THE RAPE OF LUCRECE",
	43: "VENUS AND ADONIS"
]
 */

type Searcher struct {
	CompleteWorks       string
	WorksSuffixArrayMap map[int]*suffixarray.Index
	AllWorksSuffixArray *suffixarray.Index
}

func New() Searcher {
	return Searcher{
		WorksSuffixArrayMap: make(map[int]*suffixarray.Index),
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.setup(dat)
	return nil
}
func (s *Searcher) setup(dat []byte){
	s.CompleteWorks = string(dat)

	var completeWorkLowercase = []byte(strings.ToLower(s.CompleteWorks))
	s.AllWorksSuffixArray = suffixarray.New(completeWorkLowercase)

	s.processWorksSuffixIndex(dat, completeWorkLowercase)
}

// Split all works with their respective Suffix Arrays
func (s *Searcher) processWorksSuffixIndex(data []byte, workLower []byte) {
	var nameIdx = [len(WORKS)]int{}
	var localSuffixArray = suffixarray.New(data)

	for i, work := range WORKS {
		idxs := localSuffixArray.Lookup([]byte(work), -1)
		if len(idxs) == 0 {
			fmt.Printf("ERROR %s", work)
		}
		sort.Ints(idxs)
		nameIdx[i] = idxs[1]
	}
	for i, _ := range WORKS {
		start := nameIdx[i]
		var end int
		if i == len(WORKS)-1 {
			end = len(workLower) - 1
		} else {
			end = nameIdx[i+1]
		}
		s.WorksSuffixArrayMap[i] = suffixarray.New(workLower[start:end])
	}
}

func (s *Searcher) Search(query string) []string {
	idxs := s.AllWorksSuffixArray.Lookup([]byte(query), -1)
	results := []string{}

	if idxs == nil{
		return results
	}

	idxs = optimizeIndices(idxs)

	for _, idx := range idxs {
		start, end := correctIndexBounds(s, idx)
		results = append(results, s.CompleteWorks[start:end])
	}
	return results
}

func (s *Searcher) SearchByWork(query string, workId int) []string {
	idxs := s.WorksSuffixArrayMap[workId].Lookup([]byte(query), -1)
	results := []string{}

	if idxs == nil{
		return results
	}

	idxs = optimizeIndices(idxs)

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
