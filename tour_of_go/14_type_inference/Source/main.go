// When declaring a variable without specifying an explicit type (either by using the := syntax or var = expression syntax), the variable's type is inferred from the value on the right hand side.

// When the right hand side of the declaration is typed, the new variable is of that same type:

// var i int
// j := i // j is an int
// But when the right hand side contains an untyped numeric constant, the new variable may be an int, float64, or complex128 depending on the precision of the constant:

// i := 42           // int
// f := 3.142        // float64
// g := 0.867 + 0.5i // complex128
// Try changing the initial value of v in the example code and observe how its type is affected.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	

	
// initGame is called once when the game starts
func initGame() {
	v := 23.5
	pd.Graphics.DrawText(fmt.Sprintf("v is of type %T", v), 50, 50)
	
	v2 := 18
	pd.Graphics.DrawText(fmt.Sprintf("v is of type %T", v2), 50, 70)
	
	v3 := true
	pd.Graphics.DrawText(fmt.Sprintf("v is of type %T", v3), 50, 90)
	
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
