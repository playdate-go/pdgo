
// Switch cases evaluate cases from top to bottom, stopping when a case succeeds.

// (For example,

// switch i {
// case 0:
// case f():
// }
// does not call f if i==0.)

// Note: Time in the Go playground always appears to start at 2009-11-10 23:00:00 UTC, a value whose significance is left as an exercise for the reader.

package main

import (
	"github.com/playdate-go/pdgo"
	"time"
)


var pd  *pdgo.PlaydateAPI
	
	
// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText("When's Saturday?", 50, 50)
	
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		pd.Graphics.DrawText("Today.", 50, 70)
	case today + 1:
		pd.Graphics.DrawText("Tomorrow.", 50, 70)
	case today + 2:
		pd.Graphics.DrawText("In two days.", 50, 70)
	default:
		pd.Graphics.DrawText("Too far away.", 50, 70)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
