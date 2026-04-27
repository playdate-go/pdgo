// Go functions may be closures. A closure is a function value that references variables from outside its body. The function may access and assign to the referenced variables; in this sense the function is "bound" to the variables.

// For example, the adder function returns a closure. Each closure is bound to its own sum variable.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

// initGame is called once when the game starts
func initGame() {
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		pd.Graphics.DrawText(fmt.Sprintf("pos(%d)=%d neg(%d)=%d", i, pos(i), -2*i, neg(-2*i)), 10, 20+i*20)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
