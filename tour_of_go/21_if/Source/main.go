// Go's if statements are like its for loops; the expression need not be surrounded by parentheses ( ) but the braces { } are required.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)


var pd  *pdgo.PlaydateAPI
	

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}
	
// initGame is called once when the game starts
func initGame() {
	
	
	pd.Graphics.DrawText(fmt.Sprint(sqrt(2), " and ", sqrt(-4)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
