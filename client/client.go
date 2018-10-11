package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/zacatac/scrabble/board"
	"net/url"
)

const (
	boardPath       = "board"
	checkPath       = "check"
	isValidWordPath = "is-word"

	protocol    = "https"
	host        = "scrabble-server.now.sh"
	basePath    = "/api/scrabble"
	contentType = "application/json"
)

type ScrabbleClient interface {
	Board(level string) board.Board
	Check(params CheckParams) bool
	IsValidWord(word string) bool
}

type scrabbleClient struct{}

func NewScrabbleClient() ScrabbleClient {
	return &scrabbleClient{}
}

type boardRequest struct {
	Level string `json:"level"`
}

func (c *scrabbleClient) Board(level string) board.Board {
	var respBytes []byte
	b, _ := json.Marshal(boardRequest{
		Level: level,
	})

	resp, err := http.Post(c.buildURL(boardPath), contentType, bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return board.Board{}
	}

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(respBytes))
	return board.NewBoardFromJSON(respBytes)
}

type CheckParams struct {
	Level string
	Valid bool
}

type checkResponse struct {
	Correct   bool   `json:"correct"`
	NextLevel string `json:"nextLevel"`
	Message   string `json:"message"`
}

type checkRequest struct {
	Level string `json:"level"`
	Valid bool   `json:"valid"`
}

func (c *scrabbleClient) Check(params CheckParams) bool {
	var respBytes []byte
	b, _ := json.Marshal(checkRequest{
		Level: params.Level,
		Valid: params.Valid,
	})

	resp, err := http.Post(c.buildURL(checkPath), contentType, bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return false
	}
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	log.Println(string(respBytes))
	var response checkResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		log.Fatal(err)
	}
	log.Println(response)
	return response.Correct
}

type isValidWordResponse struct {
	IsWord bool `json:"isWord"`
}

type isValidWordRequest struct {
	Word string `json:"word"`
}

func (c *scrabbleClient) IsValidWord(word string) bool {
	var respBytes []byte
	b, _ := json.Marshal(isValidWordRequest{
		Word: word,
	})

	resp, err := http.Post(c.buildURL(isValidWordPath), contentType, bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return false
	}
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	var response isValidWordResponse
	if err = json.Unmarshal(respBytes, &response); err != nil {
		log.Println(err)
		return false
	}

	log.Println(string(respBytes))
	log.Println(response)
	return response.IsWord
}

func (c *scrabbleClient) buildURL(requestPath string) string {
	uri := &url.URL{
		Scheme: protocol,
		Host:   host,
		Path:   path.Join(basePath, requestPath),
	}
	return uri.String()
}

func NewMockScrabbleClient() ScrabbleClient {
	return &mockScrabbleClient{}
}

type mockScrabbleClient struct{}

func (c *mockScrabbleClient) Board(level string) board.Board {
	return board.Board{
		Board: [][]string{
			{"c", "a", "t", "s"},
			{"", "r", "", "a"},
			{"", "t", "", "p"},
			{"", "s", "", ""},
		},
		Candidate: []board.Candidate{
			{
				Row:    3,
				Column: 3,
				Letter: "s",
			},
		},
	}
}

func (c *mockScrabbleClient) Check(params CheckParams) bool {
	return true
}

func (c *mockScrabbleClient) IsValidWord(word string) bool {
	return true
}
