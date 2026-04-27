// The type [n]T is an array of n values of type T.

// The expression

// var a [10]int
// declares a variable a as an array of ten integers.

// An array's length is part of its type, so arrays cannot be resized. This seems limiting, but don't worry; Go provides a convenient way of working with arrays.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"

	pd.Graphics.DrawText(a[0]+" "+a[1], 50, 50)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	ps := "["
	for i, v := range primes {
		if i > 0 { ps += " " }
		ps += fmt.Sprint(v)
	}
	ps += "]"
	pd.Graphics.DrawText(ps, 50, 70)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
