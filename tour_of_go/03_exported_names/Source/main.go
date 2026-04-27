// In Go, a name is exported if it begins with a capital letter. For example, Pizza is an exported name, as is Pi, which is exported from the math package.

// pizza and pi do not start with a capital letter, so they are not exported.

// When importing a package, you can refer only to its exported names. Any "unexported" names are not accessible from outside the package.

// Run the code. Notice the error message.

// Please pay attention that pd.Graphics.DrawText(math.pi, 50, 50) where Pi is in lowercase will not pass.


package main

import (
	"github.com/playdate-go/pdgo"
    "math"	
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(math.Pi, 50,50)
}

// update is called every frame
func update() int {
	return 1
}
