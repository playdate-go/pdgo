// Slices can be created with the built-in make function; this is how you create dynamically-sized arrays.

// The make function allocates a zeroed array and returns a slice that refers to that array:

// a := make([]int, 5)  // len(a)=5
// To specify a capacity, pass a third argument to make:

// b := make([]int, 0, 5) // len(b)=0, cap(b)=5

// b = b[:cap(b)] // len(b)=5, cap(b)=5
// b = b[1:]      // len(b)=4, cap(b)=4

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


func printSlice(label string, x []int, y int) {
	r := fmt.Sprintf("%s len=%d cap=%d [", label, len(x), cap(x))
	for i, v := range x {
		if i > 0 { r += " " }
		r += fmt.Sprint(v)
	}
	r += "]"
	pd.Graphics.DrawText(r, 10, y)
}

// initGame is called once when the game starts
func initGame() {
	a := make([]int, 5)
	printSlice("a", a, 50)

	b := make([]int, 0, 5)
	printSlice("b", b, 70)

	c := b[:2]
	printSlice("c", c, 90)

	d := c[2:5]
	printSlice("d", d, 110)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
