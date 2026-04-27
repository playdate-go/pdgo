// Here we see the Abs and Scale methods rewritten as functions.

// Again, try removing the * from the Scale function. Can you see why the behavior changes? What else did you need to change for the example to compile?

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

func Scale(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

// initGame is called once when the game starts
func initGame() {
	v := Vertex{3, 4}
	Scale(&v, 10)
	pd.Graphics.DrawText(fmt.Sprintf("Abs=%.2f", Abs(v)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
