// Insert or update an element in map m:

// m[key] = elem
// Retrieve an element:

// elem = m[key]
// Delete an element:

// delete(m, key)
// Test that a key is present with a two-value assignment:

// elem, ok = m[key]
// If key is in m, ok is true. If not, ok is false.

// If key is not in the map, then elem is the zero value for the map's element type.

// Note: If elem or ok have not yet been declared you could use a short declaration form:

// elem, ok := m[key]

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	m := make(map[string]int)

	m["Answer"] = 42
	pd.Graphics.DrawText(fmt.Sprintf("The value: %d", m["Answer"]), 50, 30)

	m["Answer"] = 48
	pd.Graphics.DrawText(fmt.Sprintf("The value: %d", m["Answer"]), 50, 50)

	delete(m, "Answer")
	pd.Graphics.DrawText(fmt.Sprintf("The value: %d", m["Answer"]), 50, 70)

	v, ok := m["Answer"]
	pd.Graphics.DrawText(fmt.Sprintf("The value: %d Present? %t", v, ok), 50, 90)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
