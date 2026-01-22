//go:build tinygo

// TinyGo implementation of Network API
// Note: Network is only available on Playdate Simulator, not on device

package pdgo

// Network provides access to network functions (Simulator only)
type Network struct{}

func newNetwork() *Network {
	return &Network{}
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	ptr uintptr
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// NetworkCallback is called when a request completes
type NetworkCallback func(response *HTTPResponse, err error)

// Get performs an HTTP GET request (Simulator only)
func (n *Network) Get(url string, callback NetworkCallback) *HTTPRequest {
	// Network is not available on device
	if callback != nil {
		callback(nil, &networkError{msg: "network not available on device"})
	}
	return nil
}

// Post performs an HTTP POST request (Simulator only)
func (n *Network) Post(url string, body []byte, contentType string, callback NetworkCallback) *HTTPRequest {
	// Network is not available on device
	if callback != nil {
		callback(nil, &networkError{msg: "network not available on device"})
	}
	return nil
}

// Cancel cancels a pending request
func (n *Network) Cancel(request *HTTPRequest) {
	// No-op on device
}

type networkError struct {
	msg string
}

func (e *networkError) Error() string {
	return "network: " + e.msg
}
