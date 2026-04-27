// Go does not have classes. However, you can define methods on types.

// A method is a function with a special receiver argument.

// The receiver appears in its own argument list between the func keyword and the method name.

// In this example, the Abs method has a receiver of type Vertex named v.

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

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// initGame is called once when the game starts
func initGame() {
	v := Vertex{3, 4}
	pd.Graphics.DrawText(fmt.Sprintf("Abs=%.2f", v.Abs()), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
