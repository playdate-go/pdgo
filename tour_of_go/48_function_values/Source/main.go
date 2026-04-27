// Functions are values too. They can be passed around just like other values.

// Function values may be used as function arguments and return values.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI


func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

// initGame is called once when the game starts
func initGame() {
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	pd.Graphics.DrawText(fmt.Sprintf("hypot(5,12)=%.2f", hypot(5, 12)), 50, 30)
	pd.Graphics.DrawText(fmt.Sprintf("compute(hypot)=%.2f", compute(hypot)), 50, 50)
	pd.Graphics.DrawText(fmt.Sprintf("compute(Pow)=%.2f", compute(math.Pow)), 50, 70)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
