// The var statement declares a list of variables; as in function argument lists, the type is last.

// A var statement can be at package or function level. We see both in this example.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)

var c, python, java bool

var pd  *pdgo.PlaydateAPI
	

// initGame is called once when the game starts
func initGame() {
	var i int
	pd.Graphics.DrawText(fmt.Sprint(i, c, python, java), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
