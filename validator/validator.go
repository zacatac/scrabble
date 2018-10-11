package validator

import (
	"log"

	"github.com/zacatac/scrabble/board"
	"github.com/zacatac/scrabble/client"
	"sort"
)

type Validator interface {
	ValidateMove(board board.Board) bool
}

type validator struct {
	client client.ScrabbleClient
}

func NewValidator(client client.ScrabbleClient) Validator {
	return &validator{
		client: client,
	}
}

func (v *validator) ValidateMove(board board.Board) bool {
	return candidatesAdjacent(board) &&
		candidatesHaveValidPlacement(board) &&
		v.formedWordsAreValid(board)
}

func (v *validator) formedWordsAreValid(board board.Board) bool {
	wordsFormed, validWords := wordsFormedByCandidate(board)
	if !validWords {
		return false
	}
	for _, word := range wordsFormed {
		if !v.client.IsValidWord(word) {
			return false
		}
	}
	return true
}

const minimumLengthWord = 2

func wordsFormedByCandidate(board board.Board) ([]string, bool) {
	type index struct {
		row    int
		column int
	}

	type formedWord struct {
		indices []index
		word    string
	}

	formedWords := []formedWord{}

	applyCandidateToBoard(board)
	defer removeCandidateFromBoard(board)

	for i, row := range board.Board {
		currentWord := formedWord{}
		for j, letter := range row {
			if letter == "" && currentWord.word != "" {
				if len(currentWord.word) >= minimumLengthWord {
					formedWords = append(formedWords, currentWord)
				}
				currentWord = formedWord{}
			}
			currentWord.indices = append(
				currentWord.indices, index{row: i, column: j})
			currentWord.word += letter

		}
		if len(currentWord.word) >= minimumLengthWord {
			formedWords = append(formedWords, currentWord)
		}
	}

	for j := 0; j < len(board.Board[0]); j++ {
		currentWord := formedWord{}
		for i := 0; i < len(board.Board); i++ {
			letter := board.Board[i][j]
			if letter == "" && currentWord.word != "" {
				if len(currentWord.word) >= minimumLengthWord {
					formedWords = append(formedWords, currentWord)
				}
				currentWord = formedWord{}
			}
			currentWord.indices = append(
				currentWord.indices, index{row: i, column: j})
			currentWord.word += letter
		}
		if len(currentWord.word) >= minimumLengthWord {
			formedWords = append(formedWords, currentWord)
		}
	}

	words := []string{}
	for _, fw := range formedWords {
		// Verify that word is connected in other letters
		// caught by notattached case
		numCandidateLetters := 0
		for _, index := range fw.indices {
			for _, candidate := range board.Candidate {
				if index.column == candidate.Column &&
					index.row == candidate.Row {
					numCandidateLetters++
				}
			}
		}

		if numCandidateLetters == len(board.Candidate) &&
			len(fw.word) == numCandidateLetters { // caught by directions case
			return words, false
		}

		words = append(words, fw.word)
	}
	// TODO: Optimize by not checking words that don't contain a candidate word
	log.Println("formed words:", words)
	return words, true
}

const (
	rowAligned    = "row"
	columnAligned = "column"
	notAligned    = "notAligned"
)

func applyCandidateToBoard(board board.Board) {
	for _, candidate := range board.Candidate {
		board.Board[candidate.Row][candidate.Column] = candidate.Letter
	}
}
func removeCandidateFromBoard(board board.Board) {
	for _, candidate := range board.Candidate {
		board.Board[candidate.Row][candidate.Column] = ""
	}
}

func candidatesAdjacent(board board.Board) bool {
	adjaceny := candidatesAdjacency(board)
	log.Println("adjaceny:", adjaceny)
	return adjaceny != notAligned
}

func candidatesAdjacency(board board.Board) string {
	if len(board.Candidate) == 0 {
		return notAligned
	}
	firstCandidate := board.Candidate[0]
	isInSameRow := true
	isInSameColumn := true
	for _, candidate := range board.Candidate[1:] {
		if candidate.Column != firstCandidate.Column {
			isInSameColumn = false
		}
		if candidate.Row != firstCandidate.Row {
			isInSameRow = false
		}
	}

	adjaceny := notAligned
	if isInSameColumn {
		adjaceny = columnAligned
	}
	if isInSameRow {
		adjaceny = rowAligned
	}

	rowSorter := func(i, j int) bool {
		return board.Candidate[i].Row < board.Candidate[j].Row
	}
	columnSorter := func(i, j int) bool {
		return board.Candidate[i].Column < board.Candidate[j].Column
	}

	var sorter func(i, j int) bool
	// assumes adjaceny has been validated
	switch adjaceny {
	case rowAligned:
		sorter = columnSorter
	case columnAligned:
		sorter = rowSorter
	}

	// sorting the slice guarantees that an adjacent candidate in the
	// list is adjacent on the board
	sort.SliceStable(board.Candidate, sorter)

	if isInSameColumn {
		for i := 1; i < len(board.Candidate); i++ {
			if i+board.Candidate[0].Row != board.Candidate[i].Row {
				return notAligned
			}
		}
		return columnAligned
	}
	if isInSameRow {
		for i := 1; i < len(board.Candidate); i++ {
			if i+board.Candidate[0].Column != board.Candidate[i].Column {
				return notAligned
			}
		}

		return rowAligned
	}
	return notAligned
}

func candidatesHaveValidPlacement(board board.Board) bool {
	for _, candidate := range board.Candidate {
		if !candidateHasValidPlacement(board, candidate) {
			return false
		}
	}
	log.Println("candidate has valid placement")
	return true
}

func candidatesAreUnique(board board.Board) bool {
	type key struct {
		row    int
		column int
	}
	visitedCandidates := map[key]bool{}
	for _, candidate := range board.Candidate {
		k := key{
			row:    candidate.Row,
			column: candidate.Column,
		}
		if _, found := visitedCandidates[k]; found {
			return false
		}
		visitedCandidates[k] = true

	}
	return false
}

func candidateHasValidPlacement(board board.Board, candidate board.Candidate) bool {
	validRow := candidate.Row >= 0 && candidate.Row < len(board.Board)
	validColumn := candidate.Column >= 0 && candidate.Column < len(board.Board[candidate.Row])
	placeEmpty := board.Board[candidate.Row][candidate.Column] == ""
	return validRow && validColumn && placeEmpty
}
