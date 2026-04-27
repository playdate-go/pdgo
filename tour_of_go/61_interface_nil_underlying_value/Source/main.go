// If the concrete value inside the interface itself is nil, the method will be called with a nil receiver.

// In some languages this would trigger a null pointer exception, but in Go it is common to write methods that gracefully handle being called with a nil receiver (as with the method M in this example.)

// Note that an interface value that holds a nil concrete value is itself non-nil.

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

func (t *T) M() {
	if t == nil {
		pd.Graphics.DrawText("<nil>", 10, 70)
		return
	}
	pd.Graphics.DrawText(t.S, 10, 110)
}

// initGame is called once when the game starts
func initGame() {
	var i I

	var t *T
	i = t
	// describe(i) — (<nil>, *main.T)
	pd.Graphics.DrawText("(<nil>, *main.T)", 10, 50)
	i.M()

	i = &T{"hello"}
	// describe(i) — (&{hello}, *main.T)
	pd.Graphics.DrawText("(&{hello}, *main.T)", 10, 90)
	i.M()
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
