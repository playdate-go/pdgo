// Networking example for Playdate in Go
// Ported from Playdate SDK C_API/Examples/networking
// Note: Network is only available on Playdate Simulator, not on device

package main

import (
	"fmt"

	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

// State machine for the networking demo
type State int

const (
	StateRequestingHTTPAccess State = iota
	StateWaitingForHTTPAccess
	StateHTTPGet
	StateWaitingForHTTPResponse
	StateRequestingTCPAccess
	StateWaitingForTCPAccess
	StateTCPGet
	StateWaitingForTCPResponse
)

var state State = StateRequestingHTTPAccess

// HTTP connection
var httpConn *pdgo.HTTPConnection

// TCP connection
var tcpConn *pdgo.TCPConnection

// Server configuration - change these to match your test server
const (
	httpServer = "localhost"
	httpPort   = 65433
	tcpServer  = "localhost"
	tcpPort    = 65431
)

// initGame is called once when the game starts
func initGame() {
	pd.System.LogToConsole("Networking Demo: Starting")
}

// requestHTTPAccess requests permission to make HTTP connections
func requestHTTPAccess() {
	pd.System.LogToConsole("Requesting HTTP access...")
	http := pd.Network.HTTP()
	reply := http.RequestAccess(httpServer, httpPort, false, "for testing", func(allowed bool) {
		pd.System.LogToConsole(fmt.Sprintf("HTTPAccessCallback: %v", allowed))
		if allowed {
			state = StateHTTPGet
		} else {
			pd.System.LogToConsole("HTTP access denied, skipping to TCP")
			state = StateRequestingTCPAccess
		}
	})

	// Handle immediate reply (cached permission)
	if reply != pdgo.AccessAsk {
		pd.System.LogToConsole(fmt.Sprintf("HTTP access reply immediate: %d", reply))
		if reply == pdgo.AccessAllow {
			state = StateHTTPGet
		} else {
			pd.System.LogToConsole("HTTP access denied, skipping to TCP")
			state = StateRequestingTCPAccess
		}
	} else {
		state = StateWaitingForHTTPAccess
	}
}

// requestHTTP performs the HTTP GET request
func requestHTTP() {
	pd.System.LogToConsole("Creating HTTP connection...")
	http := pd.Network.HTTP()
	httpConn = http.NewConnection(httpServer, httpPort, false)
	if httpConn == nil {
		pd.System.LogToConsole("Failed to create HTTP connection")
		state = StateRequestingTCPAccess
		return
	}

	// Set up callbacks
	httpConn.SetHeaderReceivedCallback(func(key, value string) {
		pd.System.LogToConsole(fmt.Sprintf("HTTP header: %s = %s", key, value))
	})

	httpConn.SetHeadersReadCallback(func() {
		pd.System.LogToConsole("HTTP headers read")
	})

	httpConn.SetRequestCompleteCallback(func() {
		avail := httpConn.GetBytesAvailable()
		err := httpConn.GetError()

		if err != pdgo.NetOK {
			pd.System.LogToConsole(fmt.Sprintf("HTTP request complete, err=%d", err))
		} else {
			pd.System.LogToConsole(fmt.Sprintf("HTTP request complete, %d bytes available", avail))
		}

		// Read all available data
		for avail > 0 {
			buf := make([]byte, 256)
			n := httpConn.Read(buf)
			if n > 0 {
				pd.System.LogToConsole(fmt.Sprintf("HTTP data: %s", string(buf[:n])))
			}
			avail = httpConn.GetBytesAvailable()
		}

		state = StateRequestingTCPAccess
	})

	httpConn.SetConnectionClosedCallback(func() {
		pd.System.LogToConsole("HTTP connection closed")
	})

	// Perform GET request
	err := httpConn.Get("/PING", nil)
	pd.System.LogToConsole(fmt.Sprintf("HTTP GET err=%d", err))
}

// requestTCPAccess requests permission to make TCP connections
func requestTCPAccess() {
	pd.System.LogToConsole("Requesting TCP access...")
	tcp := pd.Network.TCP()
	reply := tcp.RequestAccess(tcpServer, tcpPort, false, "for testing", func(allowed bool) {
		pd.System.LogToConsole(fmt.Sprintf("TCPAccessCallback: %v", allowed))
		if allowed {
			state = StateTCPGet
		} else {
			pd.System.LogToConsole("TCP access denied, looping back to HTTP")
			state = StateRequestingHTTPAccess
		}
	})

	// Handle immediate reply (cached permission)
	if reply != pdgo.AccessAsk {
		pd.System.LogToConsole(fmt.Sprintf("TCP access reply immediate: %d", reply))
		if reply == pdgo.AccessAllow {
			state = StateTCPGet
		} else {
			pd.System.LogToConsole("TCP access denied, looping back to HTTP")
			state = StateRequestingHTTPAccess
		}
	} else {
		state = StateWaitingForTCPAccess
	}
}

// requestTCP performs the TCP connection
func requestTCP() {
	pd.System.LogToConsole("Creating TCP connection...")
	tcp := pd.Network.TCP()
	tcpConn = tcp.NewConnection(tcpServer, tcpPort, false)
	if tcpConn == nil {
		pd.System.LogToConsole("Failed to create TCP connection")
		state = StateRequestingTCPAccess
		return
	}

	// Set up callbacks
	tcpConn.SetConnectionClosedCallback(func(err pdgo.PDNetErr) {
		pd.System.LogToConsole("TCP connection closed")

		avail := tcpConn.GetBytesAvailable()
		if err != pdgo.NetOK {
			pd.System.LogToConsole(fmt.Sprintf("TCP request complete, err=%d", err))
		} else {
			pd.System.LogToConsole(fmt.Sprintf("TCP request complete, %d bytes available", avail))
		}

		// Read all available data
		for avail > 0 {
			buf := make([]byte, 256)
			n := tcpConn.Read(buf)
			if n > 0 {
				pd.System.LogToConsole(fmt.Sprintf("TCP data: %s", string(buf[:n])))
			}
			avail = tcpConn.GetBytesAvailable()
		}

		// Loop back to HTTP
		state = StateRequestingHTTPAccess
	})

	// Open the connection
	err := tcpConn.Open(func(err pdgo.PDNetErr) {
		pd.System.LogToConsole("TCP connection open, sending PING")
		tcpConn.Write([]byte("PING\n"))
	})
	pd.System.LogToConsole(fmt.Sprintf("TCP open err=%d", err))
}

// update is called every frame
func update() int {
	// Clear screen
	pd.Graphics.Clear(pdgo.SolidWhite)

	// Draw status text
	pd.Graphics.DrawText("Networking Demo", 10, 10)
	pd.Graphics.DrawText("See console for output", 10, 30)

	// Display current state
	var stateText string
	switch state {
	case StateRequestingHTTPAccess:
		stateText = "Requesting HTTP access..."
		requestHTTPAccess()
	case StateWaitingForHTTPAccess:
		stateText = "Waiting for HTTP access..."
	case StateHTTPGet:
		stateText = "Making HTTP request..."
		requestHTTP()
		state = StateWaitingForHTTPResponse
	case StateWaitingForHTTPResponse:
		stateText = "Waiting for HTTP response..."
	case StateRequestingTCPAccess:
		stateText = "Requesting TCP access..."
		requestTCPAccess()
	case StateWaitingForTCPAccess:
		stateText = "Waiting for TCP access..."
	case StateTCPGet:
		stateText = "Making TCP connection..."
		requestTCP()
		state = StateWaitingForTCPResponse
	case StateWaitingForTCPResponse:
		stateText = "Waiting for TCP response..."
	}

	pd.Graphics.DrawText(stateText, 10, 50)

	// Draw FPS
	pd.System.DrawFPS(0, 0)

	return 1
}

func main() {}
