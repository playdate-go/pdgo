// Variables declared without an explicit initial value are given their zero value.

// The zero value is:

// 0 for numeric types,
// false for the boolean type, and
// "" (the empty string) for strings.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
		

// initGame is called once when the game starts
func initGame() {
	var i int
	var f float64
	var b bool
	var s string
	pd.Graphics.DrawText(fmt.Sprintf("%v %v %v %q\n", i, f, b, s), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
