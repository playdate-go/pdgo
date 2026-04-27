// The expression T(v) converts the value v to the type T.

// Some numeric conversions:

// var i int = 42
// var f float64 = float64(i)
// var u uint = uint(f)
// Or, put more simply:

// i := 42
// f := float64(i)
// u := uint(f)
// Unlike in C, in Go assignment between items of different type requires an explicit conversion. Try removing the float64 or uint conversions in the example and see what happens.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI
	

	
// initGame is called once when the game starts
func initGame() {
	var x, y int = 3, 4
	var f float64 = math.Sqrt(float64(x*x + y*y))
	var z uint = uint(f)
	pd.Graphics.DrawText(fmt.Sprint(x, y, z), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
