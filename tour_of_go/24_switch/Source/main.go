// A switch statement is a shorter way to write a sequence of if - else statements. It runs the first case whose value is equal to the condition expression.

// Go's switch is like the one in C, C++, Java, JavaScript, and PHP, except that Go only runs the selected case, not all the cases that follow. In effect, the break statement that is needed at the end of each case in those languages is provided automatically in Go. Another important difference is that Go's switch cases need not be constants, and the values involved need not be integers.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"runtime"
)


var pd  *pdgo.PlaydateAPI
	
	
// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText("Go runs on ", 50, 50)
	
	switch os := runtime.GOOS; os {
	case "darwin":
		pd.Graphics.DrawText("macOS ", 50, 70)
	case "linux":
				pd.Graphics.DrawText("Linux ", 50, 70)
	default:
		// freebsd, openbsd,
		// plan9, windows...
		pd.Graphics.DrawText(fmt.Sprintf("%s", os), 50, 70)
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
