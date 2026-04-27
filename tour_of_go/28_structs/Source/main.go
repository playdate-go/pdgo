// A struct is a collection of fields.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI

type Vertex struct {
	X int
	Y int
}	
		
// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprint(Vertex{1, 2}), 50, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
