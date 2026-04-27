// Go has only one looping construct, the for loop.

// The basic for loop has three components separated by semicolons:

// the init statement: executed before the first iteration
// the condition expression: evaluated before every iteration
// the post statement: executed at the end of every iteration
// The init statement will often be a short variable declaration, and the variables declared there are visible only in the scope of the for statement.

// The loop will stop iterating once the boolean condition evaluates to false.

// Note: Unlike other languages like C, Java, or JavaScript there are no parentheses surrounding the three components of the for statement and the braces { } are always required.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	
// initGame is called once when the game starts
func initGame() {
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	
	pd.Graphics.DrawText(fmt.Sprint(sum), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
