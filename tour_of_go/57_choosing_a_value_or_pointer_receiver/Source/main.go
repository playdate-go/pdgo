// There are two reasons to use a pointer receiver.

// The first is so that the method can modify the value that its receiver points to.

// The second is to avoid copying the value on each method call. This can be more efficient if the receiver is a large struct, for example.

// In this example, both Scale and Abs are methods with receiver type *Vertex, even though the Abs method needn't modify its receiver.

// In general, all methods on a given type should have either value or pointer receivers, but not a mixture of both.

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

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// initGame is called once when the game starts
func initGame() {
	v := &Vertex{3, 4}
	pd.Graphics.DrawText(fmt.Sprintf("Before: %.0f,%.0f Abs=%.2f", v.X, v.Y, v.Abs()), 10, 50)
	v.Scale(5)
	pd.Graphics.DrawText(fmt.Sprintf("After:  %.0f,%.0f Abs=%.2f", v.X, v.Y, v.Abs()), 10, 70)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
