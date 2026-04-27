// A map maps keys to values.

// The zero value of a map is nil. A nil map has no keys, nor can keys be added.

// The make function returns a map of the given type, initialized and ready for use.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


type Vertex struct {
	Lat, Long float64
}

var m map[string]Vertex

// initGame is called once when the game starts
func initGame() {
	m = make(map[string]Vertex)
	m["Bell Labs"] = Vertex{
		40.68433, -74.39967,
	}
	v := m["Bell Labs"]
	pd.Graphics.DrawText(fmt.Sprintf("Bell Labs: %.2f,%.2f", v.Lat, v.Long), 50, 50)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
