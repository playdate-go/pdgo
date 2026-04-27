// Test server for the Playdate networking example
// Run this before running the Playdate simulator
//
// Usage: go run main.go

package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	httpPort = 65433
	tcpPort  = 65431
)

func main() {
	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[HTTP] Request: %s %s\n", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG from HTTP server"))
	})

	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpPort))
	if err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
		return
	}

	go func() {
		fmt.Printf("[HTTP] Server listening on port %d\n", httpPort)
		if err := http.Serve(httpListener, nil); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[HTTP] Server error: %v\n", err)
		}
	}()

	// Start TCP server
	go func() {
		tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
		if err != nil {
			fmt.Printf("Failed to start TCP server: %v\n", err)
			return
		}
		defer tcpListener.Close()

		fmt.Printf("[TCP] Server listening on port %d\n", tcpPort)

		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				fmt.Printf("[TCP] Accept error: %v\n", err)
				continue
			}

			go handleTCPConnection(conn)
		}
	}()

	fmt.Println("\nTest servers running. Press Ctrl+C to stop.")
	fmt.Println("==========================================")

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nShutting down servers...")
	httpListener.Close()
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	fmt.Printf("[TCP] Connection from %s\n", addr)

	// Read data from client
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Printf("[TCP] Received from %s: %s", addr, line)

		// Send response
		response := "PONG from TCP server\n"
		conn.Write([]byte(response))
		fmt.Printf("[TCP] Sent to %s: %s", addr, response)
	}

	fmt.Printf("[TCP] Connection closed from %s\n", addr)
}
