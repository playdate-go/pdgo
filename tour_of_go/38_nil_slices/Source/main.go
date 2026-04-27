// The zero value of a slice is nil.

// A nil slice has a length and capacity of 0 and has no underlying array.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	var s []int
	pd.Graphics.DrawText(fmt.Sprintf("len=%d cap=%d", len(s), cap(s)), 50, 50)
	if s == nil {
		pd.Graphics.DrawText("nil!", 50, 70)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
