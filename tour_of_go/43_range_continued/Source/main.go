// You can skip the index or value by assigning to _.

// for i, _ := range pow
// for _, value := range pow
// If you only want the index, you can omit the second variable.

// for i := range pow

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	pow := make([]int, 10)
	for i := range pow {
		pow[i] = 1 << uint(i) // == 2**i
	}
	for _, value := range pow {
		fmt.Sprintf("%d", value) // value unused, just demonstrating syntax
	}
	// Print results manually
	for i, value := range pow {
		pd.Graphics.DrawText(fmt.Sprintf("2**%d = %d", i, value), 10, 20+i*20)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
