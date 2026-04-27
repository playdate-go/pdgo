// Inside a function, the := short assignment statement can be used in place of a var declaration with implicit type.

// Outside a function, every statement begins with a keyword (var, func, and so on) and so the := construct is not available.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	

// initGame is called once when the game starts
func initGame() {
	var i, j int = 1, 2
	k := 3
	c, python, java := true, false, "no!"
	
	pd.Graphics.DrawText(fmt.Sprint(i, j,  k, c, python, " ", java), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
