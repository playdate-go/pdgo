// Every Go program is made up of packages.

// Programs start running in package main.

// This program is using the packages with import paths "fmt" and "math/cmplx".

// By convention, the package name is the same as the last element of the import path. For instance, the "math/cmplx" package comprises files that begin with the statement package cmplx.


package main

import 	"github.com/playdate-go/pdgo"
import 	"fmt"
import	"math/cmplx"


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprintf("Square complex root of -4 = %.0f", cmplx.Sqrt(-4)), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
