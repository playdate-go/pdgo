// A function can take zero or more arguments.

// In this example, add takes two parameters of type int.

// Notice that the type comes after the variable name.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	

func add(x int, y int) int {
	return x + y
}


// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprint(add(42, 13)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
