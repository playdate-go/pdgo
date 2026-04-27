// A type assertion provides access to an interface value's underlying concrete value.

// t := i.(T)
// This statement asserts that the interface value i holds the concrete type T and assigns the underlying T value to the variable t.

// If i does not hold a T, the statement will trigger a panic.

// To test whether an interface value holds a specific type, a type assertion can return two values: the underlying value and a boolean value that reports whether the assertion succeeded.

// t, ok := i.(T)
// If i holds a T, then t will be the underlying value and ok will be true.

// If not, ok will be false and t will be the zero value of type T, and no panic occurs.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	var i interface{} = "hello"

	s := i.(string)
	pd.Graphics.DrawText(s, 10, 30)

	s, ok := i.(string)
	pd.Graphics.DrawText(fmt.Sprintf("%s %t", s, ok), 10, 50)

	f, ok := i.(float64)
	pd.Graphics.DrawText(fmt.Sprintf("%.0f %t", f, ok), 10, 70)

	// f = i.(float64) would panic — use ok check instead:
	if _, ok := i.(float64); !ok {
		pd.Graphics.DrawText("panic: interface{} is string, not float64", 10, 90)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
