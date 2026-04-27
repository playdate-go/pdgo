// Go has pointers. A pointer holds the memory address of a value.

// The type *T is a pointer to a T value. Its zero value is nil.

// var p *int
// The & operator generates a pointer to its operand.

// i := 42
// p = &i
// The * operator denotes the pointer's underlying value.

// fmt.Println(*p) // read i through the pointer p
// *p = 21         // set i through the pointer p
// This is known as "dereferencing" or "indirecting".

// Unlike C, Go has no pointer arithmetic.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	
	
// initGame is called once when the game starts
func initGame() {
	i, j := 42, 2701

	p := &i         // point to i
	pd.Graphics.DrawText(fmt.Sprint(*p), 50, 50) // read i through the pointer
	*p = 21         // set i through the pointer
	pd.Graphics.DrawText(fmt.Sprint(i), 50, 70) // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer

	pd.Graphics.DrawText(fmt.Sprint(j), 50, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
