// A struct literal denotes a newly allocated struct value by listing the values of its fields.

// You can list just a subset of fields by using the Name: syntax. (And the order of named fields is irrelevant.)

// The special prefix & returns a pointer to the struct value.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI

type Vertex struct {
	X int
	Y int
}	


var (
	v1 = Vertex{1, 2}  // has type Vertex
	v2 = Vertex{X: 1}  // Y:0 is implicit
	v3 = Vertex{}      // X:0 and Y:0
	p  = &Vertex{1, 2} // has type *Vertex
)
		
// initGame is called once when the game starts
func initGame() {
    pd.Graphics.DrawText(fmt.Sprintf("X:%d Y:%d", v1.X, v1.Y), 50, 50)                                            
	pd.Graphics.DrawText(fmt.Sprintf("X:%d Y:%d", p.X, p.Y), 50, 70)                                               
	pd.Graphics.DrawText(fmt.Sprintf("X:%d Y:%d", v2.X, v2.Y), 50, 90)
	pd.Graphics.DrawText(fmt.Sprintf("X:%d Y:%d", v3.X, v3.Y), 50, 110)
	                                                                                         
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
