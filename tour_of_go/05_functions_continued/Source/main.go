// When two or more consecutive named function parameters share a type, you can omit the type from all but the last.

// In this example, we shortened

// x int, y int

// to

// x, y int

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	

func add(x, y int) int {
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
