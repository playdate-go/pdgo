// Go's basic types are

// bool

// string

// int  int8  int16  int32  int64
// uint uint8 uint16 uint32 uint64 uintptr

// byte // alias for uint8

// rune // alias for int32
//      // represents a Unicode code point

// float32 float64

// complex64 complex128.

package main

import (
	"github.com/playdate-go/pdgo"
	"math/cmplx"
	"fmt"
)


var pd  *pdgo.PlaydateAPI
	
var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)
	

// initGame is called once when the game starts
func initGame() {
	pd.Graphics.DrawText(fmt.Sprintf("Type: %T Value: %v\n", ToBe, ToBe), 50, 50)
	pd.Graphics.DrawText(fmt.Sprintf("Type: %T Value: %v\n", MaxInt, MaxInt), 50, 70)
    pd.Graphics.DrawText(fmt.Sprintf("Type: %T Value: %v\n", z, z), 50, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
