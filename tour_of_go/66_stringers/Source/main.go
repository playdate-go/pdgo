// One of the most ubiquitous interfaces is Stringer defined by the fmt package.

// type Stringer interface {
//     String() string
// }
// A Stringer is a type that can describe itself as a string. The fmt package (and many others) look for this interface to print values.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
)


var pd  *pdgo.PlaydateAPI


type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

// initGame is called once when the game starts
func initGame() {
	a := Person{"Arthur Dent", 42}
	z := Person{"Zaphod Beeblebrox", 9001}
	// Call String() directly — fmt.Println(struct) crashes on device (Known Issue)
	pd.Graphics.DrawText(a.String(), 10, 50)
	pd.Graphics.DrawText(z.String(), 10, 70)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
