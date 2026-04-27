// You can declare a method on non-struct types, too.

// In this example we see a numeric type MyFloat with an Abs method.

// You can only declare a method with a receiver whose type is defined in the same package as the method. You cannot declare a method with a receiver whose type is defined in another package (which includes the built-in types such as int).

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI


type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

// initGame is called once when the game starts
func initGame() {
	f := MyFloat(-math.Sqrt2)
	pd.Graphics.DrawText(fmt.Sprintf("Abs=%.4f", f.Abs()), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
