// A type implements an interface by implementing its methods. There is no explicit declaration of intent, no "implements" keyword.

// Implicit interfaces decouple the definition of an interface from its implementation, which could then appear in any package without prearrangement.

package main

import (
	"github.com/playdate-go/pdgo"
)


var pd  *pdgo.PlaydateAPI


type I interface {
	M()
}

type T struct {
	S string
}

// This method means type T implements the interface I,
// but we don't need to explicitly declare that it does so.
func (t T) M() {
	pd.Graphics.DrawText(t.S, 50, 50)
}

// initGame is called once when the game starts
func initGame() {
	var i I = T{"hello"}
	i.M()
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
