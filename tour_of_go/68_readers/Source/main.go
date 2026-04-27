// The io package specifies the io.Reader interface, which represents the read end of a stream of data.

// The Go standard library contains many implementations of this interface, including files, network connections, compressors, ciphers, and others.

// The io.Reader interface has a Read method:

// func (T) Read(b []byte) (n int, err error)
// Read populates the given byte slice with data and returns the number of bytes populated and an error value. It returns an io.EOF error when the stream ends.

package main

import (
	"github.com/playdate-go/pdgo"
	"fmt"
	"io"
	"strings"
)


var pd  *pdgo.PlaydateAPI


// initGame is called once when the game starts
func initGame() {
	r := strings.NewReader("Hello, Reader!")

	b := make([]byte, 8)
	y := 30
	for {
		n, err := r.Read(b)
		errStr := "nil"
		if err != nil {
			errStr = err.Error() // Call Error() directly — %v on custom Stringer crashes on device (Known Issue)
		}
		pd.Graphics.DrawText(fmt.Sprintf("n=%d err=%s", n, errStr), 10, y)
		// Print b[:n] manually — no %v on slices (Known Issue)
		s := ""
		for i := 0; i < n; i++ {
			s += string(b[i])
		}
		pd.Graphics.DrawText("b[:n]=" + s, 10, y+15)
		y += 35
		if err == io.EOF {
			break
		}
	}
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
