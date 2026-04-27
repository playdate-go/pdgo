// Under the hood, interface values can be thought of as a tuple of a value and a concrete type:

// (value, type)
// An interface value holds a value of a specific underlying concrete type.

// Calling a method on an interface value executes the method of the same name on its underlying type.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI


type I interface {
	M()
}

type T struct {
	S string
}

func (t *T) M() {
	pd.Graphics.DrawText(t.S, 10, 70)
}

type F float64

func (f F) M() {
	pd.Graphics.DrawText(fmt.Sprintf("%f", float64(f)), 10, 110)
}

// initGame is called once when the game starts
func initGame() {
	var i I

	i = &T{"Hello"}
	// describe(i) — (*T, *main.T) — print type name directly
	pd.Graphics.DrawText("(*T, *main.T)", 10, 50)
	i.M()

	i = F(math.Pi)
	// describe(i) — (3.141593, main.F) — print value and type name directly
	pd.Graphics.DrawText(fmt.Sprintf("(%.6f, main.F)", float64(i.(F))), 10, 90)
	i.M()
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
