// Constants are declared like variables, but with the const keyword.

// Constants can be character, string, boolean, or numeric values.

// Constants cannot be declared using the := syntax.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"math"
)

const Pi = 3.14

var pd  *pdgo.PlaydateAPI
	
	
// initGame is called once when the game starts
func initGame() {
	const World = "世界"
	pd.Graphics.DrawText(World, 50, 50)	
	
	pd.Graphics.DrawText(fmt.Sprint("Hello ", World), 50, 70)
		
	pd.Graphics.DrawText(fmt.Sprint("Happy ", math.Pi, " ", "Day"), 50, 90)	


	const Truth = true
	pd.Graphics.DrawText(fmt.Sprint("Go rules? ", Truth), 50, 110)
	
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
