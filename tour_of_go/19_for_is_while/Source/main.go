// At that point you can drop the semicolons: C's while is spelled for in Go.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	
// initGame is called once when the game starts
func initGame() {
	sum := 1
	for sum < 1000 {
		sum += sum
	}
	
	pd.Graphics.DrawText(fmt.Sprint(sum), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
