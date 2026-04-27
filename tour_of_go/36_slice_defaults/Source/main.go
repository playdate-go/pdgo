// When slicing, you may omit the high or low bounds to use their defaults instead.

// The default is zero for the low bound and the length of the slice for the high bound.

// For the array

// var a [10]int
// these slice expressions are equivalent:

// a[0:10]
// a[:10]
// a[0:]
// a[:]

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

// initGame is called once when the game starts
func initGame() {
	s := []int{2, 3, 5, 7, 11, 13}

	s = s[1:4]
	pd.Graphics.DrawText(joinInts(s), 50, 50)

	s = s[:2]
	pd.Graphics.DrawText(joinInts(s), 50, 70)

	s = s[1:]
	pd.Graphics.DrawText(joinInts(s), 50, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
