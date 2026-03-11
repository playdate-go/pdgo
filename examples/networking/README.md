# Networking Demo

This example demonstrates HTTP and TCP networking in Playdate Go.

**Note: Network is only available on Playdate Simulator, not on device.**

## How it works

The demo cycles through the following states:

1. **HTTP Access Request** - Requests permission to connect to an HTTP server
2. **HTTP GET Request** - Makes a GET request to `/PING` and logs the response
3. **TCP Access Request** - Requests permission to connect to a TCP server
4. **TCP Connection** - Opens a TCP connection, sends "PING\n", and logs the response

## Running the example

### Step 1: Start the test server

A Go test server is included in `testserver/`:

```bash
cd testserver
go run main.go
```

This starts both:
- HTTP server on port 65433
- TCP server on port 65431

### Step 2: Build and run the Playdate app

```bash
cd ..
./build.sh
```

Then open the generated `Networking_sim.pdx` in the Playdate Simulator.

## Configuration

To use different servers, modify the constants in `Source/main.go`:

```go
const (
    httpServer = "localhost"
    httpPort   = 65433
    tcpServer  = "localhost"
    tcpPort    = 65431
)
```
