// A var declaration can include initializers, one per variable.

// If an initializer is present, the type can be omitted; the variable will take the type of the initializer.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)

var i, j int = 1, 2

var pd  *pdgo.PlaydateAPI
	

// initGame is called once when the game starts
func initGame() {
	var c, python, java = true, false, "no!"
	pd.Graphics.DrawText(fmt.Sprint(i, j, c, python, " ", java), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
