// The interface type that specifies zero methods is known as the empty interface:

// interface{}
// An empty interface may hold values of any type. (Every type implements at least zero methods.)

// Empty interfaces are used by code that handles values of unknown type. For example, fmt.Print takes any number of arguments of type interface{}.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	var i interface{}

	// describe(i) — (<nil>, <nil>)
	pd.Graphics.DrawText("(<nil>, <nil>)", 10, 30)

	i = 42
	// describe(i) — (42, int)
	pd.Graphics.DrawText(fmt.Sprintf("(%d, int)", i.(int)), 10, 50)

	i = "hello"
	// describe(i) — (hello, string)
	pd.Graphics.DrawText(fmt.Sprintf("(%s, string)", i.(string)), 10, 70)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
