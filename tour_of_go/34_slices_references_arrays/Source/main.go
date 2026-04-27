// A slice does not store any data, it just describes a section of an underlying array.

// Changing the elements of a slice modifies the corresponding elements of its underlying array.

// Other slices that share the same underlying array will see those changes.

package main

import (
	"github.com/playdate-go/pdgo"
)


var pd  *pdgo.PlaydateAPI


func joinStrs(s []string) string {
	r := "["
	for i, v := range s {
		if i > 0 { r += " " }
		r += v
	}
	return r + "]"
}

// initGame is called once when the game starts
func initGame() {
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	pd.Graphics.DrawText(joinStrs(names[:]), 50, 30)

	a := names[0:2]
	b := names[1:3]
	pd.Graphics.DrawText(joinStrs(a)+" "+joinStrs(b), 50, 50)

	b[0] = "XXX"
	pd.Graphics.DrawText(joinStrs(a)+" "+joinStrs(b), 50, 70)
	pd.Graphics.DrawText(joinStrs(names[:]), 50, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
