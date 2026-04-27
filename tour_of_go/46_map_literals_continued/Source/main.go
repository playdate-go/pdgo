// If the top-level type is just a type name, you can omit it from the elements of the literal.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


type Vertex struct {
	Lat, Long float64
}

var m = map[string]Vertex{
	"Bell Labs": {40.68433, -74.39967},
	"Google":    {37.42202, -122.08408},
}

// initGame is called once when the game starts
func initGame() {
	for name, v := range m {
		pd.Graphics.DrawText(fmt.Sprintf("%s: %.2f,%.2f", name, v.Lat, v.Long), 10, 30)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
