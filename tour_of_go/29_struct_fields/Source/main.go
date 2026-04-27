// Struct fields are accessed using a dot.

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
	v := Vertex{1, 2}
	v.X = 4
	pd.Graphics.DrawText(fmt.Sprint(v.X), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
