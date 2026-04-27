// The range form of the for loop iterates over a slice or map.

// When ranging over a slice, two values are returned for each iteration. The first is the index, and the second is a copy of the element at that index.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI

var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

// initGame is called once when the game starts
func initGame() {
	for i, v := range pow {
		pd.Graphics.DrawText(fmt.Sprintf("2**%d = %d", i, v), 50, 30+i*20)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
