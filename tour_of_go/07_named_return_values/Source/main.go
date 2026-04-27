// Go's return values may be named. If so, they are treated as variables defined at the top of the function.

// These names should be used to document the meaning of the return values.

// A return statement without arguments returns the named return values. This is known as a "naked" return.

// Naked return statements should be used only in short functions, as with the example shown here. They can harm readability in longer functions.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprint(split(17)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
