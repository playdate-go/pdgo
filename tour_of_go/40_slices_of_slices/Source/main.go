// Slices can contain any type, including other slices.

package main

import (
	"github.com/playdate-go/pdgo"
	"strings"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	// Create a tic-tac-toe board.
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	// The players take turns.
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		pd.Graphics.DrawText(strings.Join(board[i], " "), 50, 50+i*20)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
