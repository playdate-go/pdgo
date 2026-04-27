// Like for, the if statement can start with a short statement to execute before the condition.

// Variables declared by the statement are only in scope until the end of the if.

// (Try using v in the last return statement.)

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
	}
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
