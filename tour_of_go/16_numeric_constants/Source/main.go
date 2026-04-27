// Numeric Constants
// Numeric constants are high-precision values.

// An untyped constant takes the type needed by its context.

// Try printing needInt(Big) too.

// (An int can store at maximum a 64-bit integer, and sometimes less.)

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)

const (
	// Create a huge number by shifting a 1 bit left 100 places.
	// In other words, the binary number that is 1 followed by 100 zeroes.
	Big = 1 << 100
	// Shift it right again 99 places, so we end up with 1<<1, or 2.
	Small = Big >> 99
)

var pd  *pdgo.PlaydateAPI
	
func needInt(x int) int { return x*10 + 1 }
func needFloat(x float64) float64 {
	return x * 0.1
}

// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprint(needInt(Small)), 50, 50)
	pd.Graphics.DrawText(fmt.Sprint(needFloat(Small)), 50, 70)	
	pd.Graphics.DrawText(fmt.Sprint(needFloat(Big)), 50, 90)			
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
