// Variables declared inside an if short statement are also available inside any of the else blocks.

// (Both calls to pow return their results before the call to fmt.Println in main begins.)

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI
	

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	// can't use v here, though
	return lim
}
	
// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprint(pow(3,2,10), " and ", pow(3,3,20)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
