// A slice does not store any data, it just describes a section of an underlying array.

// Changing the elements of a slice modifies the corresponding elements of its underlying array.

// Other slices that share the same underlying array will see those changes.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	primes := [6]int{2, 3, 5, 7, 11, 13}

	var s []int = primes[1:4]

	// Manual formatting — fmt.Sprint on slices crashes on device (Known Issue)
	out := "["
	for i, v := range s {
		if i > 0 { out += " " }
		out += fmt.Sprint(v)
	}
	out += "]"
	pd.Graphics.DrawText(out, 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
