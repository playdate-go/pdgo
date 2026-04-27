// An interface type is defined as a set of method signatures.

// A value of interface type can hold any value that implements those methods.

// Note: There is an error in the example code on line 22. Vertex (the value type) doesn't implement Abser because the Abs method is defined only on *Vertex (the pointer type).

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI


type Abser interface {
	Abs() float64
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// initGame is called once when the game starts
func initGame() {
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	a = f  // a MyFloat implements Abser
	pd.Graphics.DrawText(fmt.Sprintf("MyFloat Abs=%.4f", a.Abs()), 10, 30)

	a = &v // a *Vertex implements Abser
	pd.Graphics.DrawText(fmt.Sprintf("*Vertex Abs=%.2f", a.Abs()), 10, 50)

	// a = v would NOT compile: Vertex does not implement Abser, only *Vertex does
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
