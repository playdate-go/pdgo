// Switch without a condition is the same as switch true.

// This construct can be a clean way to write long if-then-else chains.

package main

import (
	"github.com/playdate-go/pdgo"
	"time"
)


var pd  *pdgo.PlaydateAPI
	
	
// initGame is called once when the game starts
func initGame() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		pd.Graphics.DrawText("Good morning!", 50, 50)
	case t.Hour() < 17:
		pd.Graphics.DrawText("Good afternoon!", 50, 50)
	default:
		pd.Graphics.DrawText("Good evening!", 50, 50)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
