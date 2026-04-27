// Sprite Game - Shoot-em-up example for Playdate
package main

import (
	"spritegame/core"

	"github.com/playdate-go/pdgo"
)

var (
	pd   *pdgo.PlaydateAPI
	game *core.Game
)

// initGame is called once when the game starts
func initGame() {
	pd.Display.SetRefreshRate(30)

	game = core.New(pd)
	game.Setup()
}

// update is called every frame
func update() int {
	return game.Update()
}

func main() {}
