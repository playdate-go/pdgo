// Go programs express error state with error values.

// The error type is a built-in interface similar to fmt.Stringer:

// type error interface {
//     Error() string
// }
// A nil error denotes success; a non-nil error denotes failure.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


type MyError struct {
	When string
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %s, %s", e.When, e.What)
}

func run() error {
	return &MyError{
		"2026-04-27",
		"it didn't work",
	}
}

// initGame is called once when the game starts
func initGame() {
	if err := run(); err != nil {
		// Call Error() directly — fmt.Println(err) crashes on device with custom Stringer (Known Issue)
		pd.Graphics.DrawText(err.Error(), 10, 50)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
