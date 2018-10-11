# Scrabble move checker

This repository contains a move checker for the game scrabble.

Run with `go run main.go`

You can configure the level to evaluate by setting the currentLevel constant in main.go.

If you want to see additional logs, the mainDebug flag to true.

## Scrabble Details

### Board

	- Peach pieces are the current game board state
	- Gray pieces are candidate moves
	- Use the `API` below to fetch the board state and candidate move

### Rules

	- Played moves have to be adjacent to the current board (in at least one place)
	- Valid words are all adjacent pieces from top to bottom, left to right.
	- Played moves have to be adjacent, and all in the same row or column
	- Every individual piece in a move must make valid dictionary words with all adjacent pieces

## Plan

	- Create a client for the scrabble API
	  - http POST "https://scrabble-server.now.sh/api/scrabble/board" {"level": YOUR-LEVEL}
	  - {board: matrix[i][j] = letter|null, candidate: array[{row: i, col: j, letter:a-z}]}
	  - http POST "https://scrabble-server.now.sh/api/scrabble/check" {"level": YOUR-LEVEL, valid: boolean}
	  - {nextLevel: keyword, correct: true}
	  - {correct: false, message: string}
	  - http POST "https://scrabble-server.now.sh/api/scrabble/is-word" {"word": string}
	- Define data model for words, moves, and boards such that we can efficiently check the validity of a move
	- Start implementing data models
	- More clearly define how we will check the validity of moves
	- Implement move checkers
