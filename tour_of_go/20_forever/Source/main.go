// If you omit the loop condition it loops forever, so an infinite loop is compactly expressed.

package main

import (
	"github.com/playdate-go/pdgo"
)


var pd  *pdgo.PlaydateAPI
	
// initGame is called once when the game starts
func initGame() {
	for {
	}
	
	pd.Graphics.DrawText("hi", 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
