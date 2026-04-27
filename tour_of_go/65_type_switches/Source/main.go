// A type switch is a construct that permits several type assertions in series.

// A type switch is like a regular switch statement, but the cases in a type switch specify types (not values), and those values are compared against the type of the value held by the given interface value.

// switch v := i.(type) {
// case T:
//     // here v has type T
// case S:
//     // here v has type S
// default:
//     // no match; here v has the same type as i
// }
// The declaration in a type switch has the same syntax as a type assertion i.(T), but the specific type T is replaced with the keyword type.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


func do(i interface{}, y int) {
	switch v := i.(type) {
	case int:
		pd.Graphics.DrawText(fmt.Sprintf("Twice %d is %d", v, v*2), 10, y)
	case string:
		pd.Graphics.DrawText(fmt.Sprintf("\"%s\" is %d bytes long", v, len(v)), 10, y)
	default:
		pd.Graphics.DrawText(fmt.Sprintf("I don't know about type %T!", v), 10, y)
	}
}

// initGame is called once when the game starts
func initGame() {
	do(21, 50)
	do("hello", 70)
	do(true, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
