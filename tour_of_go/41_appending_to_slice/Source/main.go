// It is common to append new elements to a slice, and so Go provides a built-in append function.

// func append(s []T, vs ...T) []T
// The first parameter s of append is a slice of type T, and the rest are T values to append to the slice.

// The resulting value of append is a slice containing all the elements of the original slice plus the provided values.

// If the backing array of s is too small to fit all the given values a bigger array will be allocated. The returned slice will point to the newly allocated array.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


func printSlice(s []int, y int) {
	r := fmt.Sprintf("len=%d cap=%d [", len(s), cap(s))
	for i, v := range s {
		if i > 0 { r += " " }
		r += fmt.Sprint(v)
	}
	r += "]"
	pd.Graphics.DrawText(r, 10, y)
}

// initGame is called once when the game starts
func initGame() {
	var s []int
	printSlice(s, 50)

	s = append(s, 0)
	printSlice(s, 70)

	s = append(s, 1)
	printSlice(s, 90)

	s = append(s, 2, 3, 4)
	printSlice(s, 110)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
