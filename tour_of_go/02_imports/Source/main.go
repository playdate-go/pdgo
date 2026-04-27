// This code groups the imports into a parenthesized, "factored" import statement.

// You can also write multiple import statements, like:

// import "fmt"
// import "math/cmplx"

// But it is good style to use the factored import statement.


package main

import (
	"github.com/playdate-go/pdgo"
    "fmt"
    "math/cmplx"	
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprintf("Now you have %.0f problems", cmplx.Sqrt(-7)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
