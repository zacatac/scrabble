package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zacatac/scrabble/client"
	"github.com/zacatac/scrabble/validator"
)

const (
	mainDebug  = false
	mockClient = false
	flags      = log.Lshortfile | log.Ltime

	currentLevel = "start"
)

var scrabbleClient client.ScrabbleClient

func init() {
	log.SetFlags(flags)
	if mainDebug {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	if mockClient {
		scrabbleClient = client.NewMockScrabbleClient()
	} else {
		scrabbleClient = client.NewScrabbleClient()
	}
	completedLevels := []string{
		"start",
		"overrun",
		"nonlinear",
		"notaword",
		"notattached", // missed this candidate validation
		"directions",  // fix for notattached case too stringent
	}
	log.Println(completedLevels)
	playLevel(currentLevel)
}

func playLevel(level string) {
	validator := validator.NewValidator(scrabbleClient)
	b := scrabbleClient.Board(level)
	isValid := validator.ValidateMove(b)
	log.Println("move is valid:", isValid)
	providedCorrectAnswer := scrabbleClient.Check(client.CheckParams{
		Level: level,
		Valid: isValid,
	})
	fmt.Println("level:", level)
	if !providedCorrectAnswer {
		fmt.Printf("incorrect answer\n%s\nisValid:%t\n", b, isValid)
		return
	}
	fmt.Printf("correct answer\n%s\nisValid:%t\n", b, isValid)
}
