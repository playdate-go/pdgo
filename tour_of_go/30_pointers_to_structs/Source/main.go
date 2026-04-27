// Struct fields can be accessed through a struct pointer.
// To access the field X of a struct when we have the struct pointer p we could write (*p).X. However, that notation is cumbersome, so the language permits us instead to write just p.X, without the explicit dereference.

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
		
// initGame is called once when the game starts
func initGame() {
	v := Vertex{1, 2}
	p := &v
	p.X = 1e9
    pd.Graphics.DrawText(fmt.Sprintf("X:%d Y:%d", v.X, v.Y), 50, 50)                                                        
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
