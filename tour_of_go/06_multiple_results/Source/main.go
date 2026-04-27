// A function can return any number of results.

// The swap function returns two strings.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	

func swap(x, y string) (string, string) {
	return y, x
}

// initGame is called once when the game starts
func initGame() {
	a, b := swap("Hello", "World")
	pd.Graphics.DrawText(fmt.Sprintf("%s %s", a, b), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
