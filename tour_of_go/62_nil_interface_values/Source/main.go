// A nil interface value holds neither value nor concrete type.

// Calling a method on a nil interface is a run-time error because there is no type inside the interface tuple to indicate which concrete method to call.

package main

import (
	"github.com/playdate-go/pdgo"
)


var pd  *pdgo.PlaydateAPI


type I interface {
	M()
}

// initGame is called once when the game starts
func initGame() {
	var i I

	// describe(i) — (<nil>, <nil>)
	pd.Graphics.DrawText("(<nil>, <nil>)", 10, 50)

	// i.M() would panic: "runtime error: invalid memory address or nil pointer dereference"
	// On device, unhandled panics crash. We check instead:
	if i == nil {
		pd.Graphics.DrawText("i is nil: cannot call i.M()", 10, 70)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
