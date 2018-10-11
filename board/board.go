package board

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TODO: Rename to BoardAndCandidate
type Board struct {
	Board     [][]string
	Candidate []Candidate
}

type boardResponse struct {
	Board     [][]string  `json:"board"`
	Candidate []Candidate `json:"candidate"`
}

type Candidate struct {
	Row    int    `json:"row"`
	Column int    `json:"col"`
	Letter string `json:"letter"`
}

func NewBoardFromJSON(boardJSON []byte) Board {
	responseBoard := boardResponse{}
	json.Unmarshal(boardJSON, &responseBoard)
	return Board{
		Board:     responseBoard.Board,
		Candidate: responseBoard.Candidate,
	}
}

func (b Board) String() string {
	s := "Board:\n"
	for _, row := range b.Board {
		// Print gaps
		// if i == 0 || i == len(b.Board)-1 {
		s += fmt.Sprintln(strings.Repeat("-", len(b.Board[0])*2))
		// }
		for _, col := range row {
			if col == "" {
				col = " "
			}
			s += fmt.Sprintf("|%s", col)
		}
		s += "\n"
	}
	bytes, _ := json.MarshalIndent(b.Candidate, "", " ")
	return s + "\nCandidate:\n" + string(bytes)
}
