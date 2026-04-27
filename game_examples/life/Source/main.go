// Game of Life example for Playdate in Go
// Port of the C example from Playdate SDK

package main

import "github.com/playdate-go/pdgo"

var (
	pd   *pdgo.PlaydateAPI
	game *Game
)

// initGame is called once when the game starts
func initGame() {
	game = NewGame(pd)
	game.Init()
}

// update is called every frame
func update() int {
	return game.Update()
}

func main() {}
