// Slice literals
// A slice literal is like an array literal without the length.

// This is an array literal:

// [3]bool{true, true, false}
// And this creates the same array as above, then builds a slice that references it:

// []bool{true, true, false}

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


func joinInts(s []int) string {
	r := "["
	for i, v := range s {
		if i > 0 { r += " " }
		r += fmt.Sprint(v)
	}
	return r + "]"
}

func joinBools(s []bool) string {
	r := "["
	for i, v := range s {
		if i > 0 { r += " " }
		r += fmt.Sprint(v)
	}
	return r + "]"
}

// initGame is called once when the game starts
func initGame() {
	q := []int{2, 3, 5, 7, 11, 13}
	pd.Graphics.DrawText(joinInts(q), 50, 30)

	r := []bool{true, false, true, true, false, true}
	pd.Graphics.DrawText(joinBools(r), 50, 50)

	type entry struct {
		i int
		b bool
	}
	s := []entry{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	es := "["
	for i, e := range s {
		if i > 0 { es += " " }
		es += fmt.Sprintf("{%d %t}", e.i, e.b)
	}
	es += "]"
	pd.Graphics.DrawText(es, 50, 70)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
