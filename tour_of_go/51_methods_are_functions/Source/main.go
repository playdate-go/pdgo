// Remember: a method is just a function with a receiver argument.

// Here's Abs written as a regular function with no change in functionality.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI


type Vertex struct {
	X, Y float64
}

func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// initGame is called once when the game starts
func initGame() {
	v := Vertex{3, 4}
	pd.Graphics.DrawText(fmt.Sprintf("Abs=%.2f", Abs(v)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
