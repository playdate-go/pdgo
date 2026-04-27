// The equivalent thing happens in the reverse direction.

// Functions that take a value argument must take a value of that specific type:

// var v Vertex
// fmt.Println(AbsFunc(v))  // OK
// fmt.Println(AbsFunc(&v)) // Compile error!
// while methods with value receivers take either a value or a pointer as the receiver when they are called:

// var v Vertex
// fmt.Println(v.Abs()) // OK
// p := &v
// fmt.Println(p.Abs()) // OK
// In this case, the method call p.Abs() is interpreted as (*p).Abs().

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

func AbsFunc(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// initGame is called once when the game starts
func initGame() {
	v := Vertex{3, 4}
	pd.Graphics.DrawText(fmt.Sprintf("v.Abs()=%.2f", v.Abs()), 10, 30)
	pd.Graphics.DrawText(fmt.Sprintf("AbsFunc(v)=%.2f", AbsFunc(v)), 10, 50)

	p := &Vertex{4, 3}
	pd.Graphics.DrawText(fmt.Sprintf("p.Abs()=%.2f", p.Abs()), 10, 70)
	pd.Graphics.DrawText(fmt.Sprintf("AbsFunc(*p)=%.2f", AbsFunc(*p)), 10, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
