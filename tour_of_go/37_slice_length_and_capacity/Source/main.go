// A slice has both a length and a capacity.

// The length of a slice is the number of elements it contains.

// The capacity of a slice is the number of elements in the underlying array, counting from the first element in the slice.

// The length and capacity of a slice s can be obtained using the expressions len(s) and cap(s).

// You can extend a slice's length by re-slicing it, provided it has sufficient capacity.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


func printSlice(s []int, y int) string {
	r := fmt.Sprintf("len=%d cap=%d [", len(s), cap(s))
	for i, v := range s {
		if i > 0 { r += " " }
		r += fmt.Sprint(v)
	}
	r += "]"
	pd.Graphics.DrawText(r, 10, y)
	return r
}

// initGame is called once when the game starts
func initGame() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s, 30)

	s = s[:0]
	printSlice(s, 50)

	s = s[:4]
	printSlice(s, 70)

	s = s[2:]
	printSlice(s, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
