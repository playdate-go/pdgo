// pdgo Network API - unified CGO implementation
// Note: Network is only available on Playdate Simulator, not on device

package pdgo

/*
#include <stdint.h>

// HTTP API
int pd_http_requestAccess(const char* server, int port, int usessl, const char* purpose, void* userdata);
void* pd_http_newConnection(const char* server, int port, int usessl);
void pd_http_retain(void* conn);
void pd_http_release(void* conn);
void pd_http_setConnectTimeout(void* conn, int ms);
void pd_http_setKeepAlive(void* conn, int keepalive);
void pd_http_setByteRange(void* conn, int start, int end);
int pd_http_get(void* conn, const char* path, const char* headers, int headerlen);
int pd_http_post(void* conn, const char* path, const char* headers, int headerlen, const char* body, int bodylen);
int pd_http_query(void* conn, const char* method, const char* path, const char* headers, int headerlen, const char* body, int bodylen);
int pd_http_getError(void* conn);
void pd_http_getProgress(void* conn, int* read, int* total);
int pd_http_getResponseStatus(void* conn);
int pd_http_getBytesAvailable(void* conn);
void pd_http_setReadTimeout(void* conn, int ms);
void pd_http_setReadBufferSize(void* conn, int bytes);
int pd_http_read(void* conn, void* buf, int buflen);
void pd_http_close(void* conn);
void pd_http_setHeaderReceivedCallback(void* conn);
void pd_http_setHeadersReadCallback(void* conn);
void pd_http_setResponseCallback(void* conn);
void pd_http_setRequestCompleteCallback(void* conn);
void pd_http_setConnectionClosedCallback(void* conn);

// TCP API
int pd_tcp_requestAccess(const char* server, int port, int usessl, const char* purpose, void* userdata);
void* pd_tcp_newConnection(const char* server, int port, int usessl);
void pd_tcp_retain(void* conn);
void pd_tcp_release(void* conn);
int pd_tcp_getError(void* conn);
void pd_tcp_setConnectTimeout(void* conn, int ms);
int pd_tcp_open(void* conn, void* userdata);
int pd_tcp_close(void* conn);
void pd_tcp_setConnectionClosedCallback(void* conn);
void pd_tcp_setReadTimeout(void* conn, int ms);
void pd_tcp_setReadBufferSize(void* conn, int bytes);
int pd_tcp_getBytesAvailable(void* conn);
int pd_tcp_read(void* conn, void* buf, int length);
int pd_tcp_write(void* conn, const void* buf, int length);

// Network status
int pd_network_getStatus(void);
*/
import "C"
import (
	"sync"
	"unsafe"
)

// PDNetErr represents network error codes
type PDNetErr int

const (
	NetOK                 PDNetErr = 0
	NetNoDevice           PDNetErr = -1
	NetBusy               PDNetErr = -2
	NetWriteError         PDNetErr = -3
	NetWriteBusy          PDNetErr = -4
	NetWriteTimeout       PDNetErr = -5
	NetReadError          PDNetErr = -6
	NetReadBusy           PDNetErr = -7
	NetReadTimeout        PDNetErr = -8
	NetReadOverflow       PDNetErr = -9
	NetFrameError         PDNetErr = -10
	NetBadResponse        PDNetErr = -11
	NetErrorResponse      PDNetErr = -12
	NetResetTimeout       PDNetErr = -13
	NetBufferTooSmall     PDNetErr = -14
	NetUnexpectedResponse PDNetErr = -15
	NetNotConnectedToAP   PDNetErr = -16
	NetNotImplemented     PDNetErr = -17
	NetConnectionClosed   PDNetErr = -18
)

// AccessReply represents the result of an access request
type AccessReply int

const (
	AccessAsk   AccessReply = 0 // Need to ask user, callback will be called
	AccessDeny  AccessReply = 1 // Access denied (cached)
	AccessAllow AccessReply = 2 // Access allowed (cached)
)

// WifiStatus represents WiFi connection status
type WifiStatus int

const (
	WifiNotConnected WifiStatus = iota
	WifiConnected
	WifiNotAvailable
)

// Network provides access to network functions (Simulator only)
type Network struct {
	httpCallbacksMu sync.Mutex
	httpCallbacks   map[uintptr]*httpCallbackData

	tcpCallbacksMu sync.Mutex
	tcpCallbacks   map[uintptr]*tcpCallbackData

	accessCallbacksMu sync.Mutex
	accessCallbacks   map[int]AccessCallback
	accessCallbackID  int
}

func newNetwork() *Network {
	return &Network{
		httpCallbacks:   make(map[uintptr]*httpCallbackData),
		tcpCallbacks:    make(map[uintptr]*tcpCallbackData),
		accessCallbacks: make(map[int]AccessCallback),
	}
}

// AccessCallback is called when access request completes
type AccessCallback func(allowed bool)

// HTTPConnection represents an HTTP connection
type HTTPConnection struct {
	ptr     uintptr
	network *Network
}

// httpCallbackData stores Go callbacks for an HTTP connection
type httpCallbackData struct {
	headerReceived   func(key, value string)
	headersRead      func()
	response         func()
	requestComplete  func()
	connectionClosed func()
}

// TCPConnection represents a TCP connection
type TCPConnection struct {
	ptr     uintptr
	network *Network
}

// tcpCallbackData stores Go callbacks for a TCP connection
type tcpCallbackData struct {
	openCallback   func(err PDNetErr)
	closedCallback func(err PDNetErr)
}

// GetWifiStatus returns the current WiFi status
func (n *Network) GetWifiStatus() WifiStatus {
	return WifiStatus(C.pd_network_getStatus())
}

// HTTP returns the HTTP API
func (n *Network) HTTP() *HTTPAPI {
	return &HTTPAPI{network: n}
}

// TCP returns the TCP API
func (n *Network) TCP() *TCPAPI {
	return &TCPAPI{network: n}
}

// HTTPAPI provides HTTP networking functions
type HTTPAPI struct {
	network *Network
}

// RequestAccess requests permission to connect to a server.
// Returns AccessReply indicating whether access was granted, denied, or pending.
// If not AccessAsk, the callback will NOT be called and the reply should be used directly.
func (h *HTTPAPI) RequestAccess(server string, port int, useSSL bool, purpose string, callback AccessCallback) AccessReply {
	cServer := make([]byte, len(server)+1)
	copy(cServer, server)
	cPurpose := make([]byte, len(purpose)+1)
	copy(cPurpose, purpose)

	var useSSLInt C.int
	if useSSL {
		useSSLInt = 1
	}

	var userdata C.int
	if callback != nil {
		h.network.accessCallbacksMu.Lock()
		h.network.accessCallbackID++
		id := h.network.accessCallbackID
		h.network.accessCallbacks[id] = callback
		h.network.accessCallbacksMu.Unlock()
		userdata = C.int(id)
	}

	reply := C.pd_http_requestAccess(
		(*C.char)(unsafe.Pointer(&cServer[0])),
		C.int(port),
		useSSLInt,
		(*C.char)(unsafe.Pointer(&cPurpose[0])),
		unsafe.Pointer(uintptr(userdata)),
	)

	// If access was immediately granted/denied, callback won't be called
	// Clean up the callback and call it now
	if reply != 0 && callback != nil { // not kAccessAsk
		h.network.accessCallbacksMu.Lock()
		delete(h.network.accessCallbacks, int(userdata))
		h.network.accessCallbacksMu.Unlock()
		// Call the callback with the result
		go callback(reply == 2) // kAccessAllow = 2
	}

	return AccessReply(reply)
}

// NewConnection creates a new HTTP connection
func (h *HTTPAPI) NewConnection(server string, port int, useSSL bool) *HTTPConnection {
	cServer := make([]byte, len(server)+1)
	copy(cServer, server)

	var useSSLInt C.int
	if useSSL {
		useSSLInt = 1
	}

	ptr := C.pd_http_newConnection(
		(*C.char)(unsafe.Pointer(&cServer[0])),
		C.int(port),
		useSSLInt,
	)

	if ptr == nil {
		return nil
	}

	conn := &HTTPConnection{
		ptr:     uintptr(unsafe.Pointer(ptr)),
		network: h.network,
	}

	// Initialize callback data
	h.network.httpCallbacksMu.Lock()
	h.network.httpCallbacks[conn.ptr] = &httpCallbackData{}
	h.network.httpCallbacksMu.Unlock()

	return conn
}

// Retain increments the reference count
func (c *HTTPConnection) Retain() {
	C.pd_http_retain(unsafe.Pointer(c.ptr))
}

// Release decrements the reference count and frees when zero
func (c *HTTPConnection) Release() {
	C.pd_http_release(unsafe.Pointer(c.ptr))
	c.network.httpCallbacksMu.Lock()
	delete(c.network.httpCallbacks, c.ptr)
	c.network.httpCallbacksMu.Unlock()
}

// SetConnectTimeout sets the connection timeout in milliseconds
func (c *HTTPConnection) SetConnectTimeout(ms int) {
	C.pd_http_setConnectTimeout(unsafe.Pointer(c.ptr), C.int(ms))
}

// SetKeepAlive enables or disables keep-alive
func (c *HTTPConnection) SetKeepAlive(keepalive bool) {
	var flag C.int
	if keepalive {
		flag = 1
	}
	C.pd_http_setKeepAlive(unsafe.Pointer(c.ptr), flag)
}

// SetByteRange sets the byte range for partial content requests
func (c *HTTPConnection) SetByteRange(start, end int) {
	C.pd_http_setByteRange(unsafe.Pointer(c.ptr), C.int(start), C.int(end))
}

// Get performs an HTTP GET request
func (c *HTTPConnection) Get(path string, headers []byte) PDNetErr {
	cPath := make([]byte, len(path)+1)
	copy(cPath, path)

	var headersPtr *C.char
	var headersLen C.int
	if len(headers) > 0 {
		headersPtr = (*C.char)(unsafe.Pointer(&headers[0]))
		headersLen = C.int(len(headers))
	}

	return PDNetErr(C.pd_http_get(unsafe.Pointer(c.ptr),
		(*C.char)(unsafe.Pointer(&cPath[0])),
		headersPtr, headersLen))
}

// Post performs an HTTP POST request
func (c *HTTPConnection) Post(path string, headers, body []byte) PDNetErr {
	cPath := make([]byte, len(path)+1)
	copy(cPath, path)

	var headersPtr, bodyPtr *C.char
	var headersLen, bodyLen C.int
	if len(headers) > 0 {
		headersPtr = (*C.char)(unsafe.Pointer(&headers[0]))
		headersLen = C.int(len(headers))
	}
	if len(body) > 0 {
		bodyPtr = (*C.char)(unsafe.Pointer(&body[0]))
		bodyLen = C.int(len(body))
	}

	return PDNetErr(C.pd_http_post(unsafe.Pointer(c.ptr),
		(*C.char)(unsafe.Pointer(&cPath[0])),
		headersPtr, headersLen, bodyPtr, bodyLen))
}

// Query performs a custom HTTP request
func (c *HTTPConnection) Query(method, path string, headers, body []byte) PDNetErr {
	cMethod := make([]byte, len(method)+1)
	copy(cMethod, method)
	cPath := make([]byte, len(path)+1)
	copy(cPath, path)

	var headersPtr, bodyPtr *C.char
	var headersLen, bodyLen C.int
	if len(headers) > 0 {
		headersPtr = (*C.char)(unsafe.Pointer(&headers[0]))
		headersLen = C.int(len(headers))
	}
	if len(body) > 0 {
		bodyPtr = (*C.char)(unsafe.Pointer(&body[0]))
		bodyLen = C.int(len(body))
	}

	return PDNetErr(C.pd_http_query(unsafe.Pointer(c.ptr),
		(*C.char)(unsafe.Pointer(&cMethod[0])),
		(*C.char)(unsafe.Pointer(&cPath[0])),
		headersPtr, headersLen, bodyPtr, bodyLen))
}

// GetError returns the last error for the connection
func (c *HTTPConnection) GetError() PDNetErr {
	return PDNetErr(C.pd_http_getError(unsafe.Pointer(c.ptr)))
}

// GetProgress returns the download progress
func (c *HTTPConnection) GetProgress() (read, total int) {
	var cRead, cTotal C.int
	C.pd_http_getProgress(unsafe.Pointer(c.ptr), &cRead, &cTotal)
	return int(cRead), int(cTotal)
}

// GetResponseStatus returns the HTTP response status code
func (c *HTTPConnection) GetResponseStatus() int {
	return int(C.pd_http_getResponseStatus(unsafe.Pointer(c.ptr)))
}

// GetBytesAvailable returns the number of bytes available to read
func (c *HTTPConnection) GetBytesAvailable() int {
	return int(C.pd_http_getBytesAvailable(unsafe.Pointer(c.ptr)))
}

// SetReadTimeout sets the read timeout in milliseconds
func (c *HTTPConnection) SetReadTimeout(ms int) {
	C.pd_http_setReadTimeout(unsafe.Pointer(c.ptr), C.int(ms))
}

// SetReadBufferSize sets the read buffer size
func (c *HTTPConnection) SetReadBufferSize(bytes int) {
	C.pd_http_setReadBufferSize(unsafe.Pointer(c.ptr), C.int(bytes))
}

// Read reads data from the connection
func (c *HTTPConnection) Read(buf []byte) int {
	if len(buf) == 0 {
		return 0
	}
	return int(C.pd_http_read(unsafe.Pointer(c.ptr),
		unsafe.Pointer(&buf[0]), C.int(len(buf))))
}

// Close closes the connection
func (c *HTTPConnection) Close() {
	C.pd_http_close(unsafe.Pointer(c.ptr))
}

// SetHeaderReceivedCallback sets the callback for when a header is received
func (c *HTTPConnection) SetHeaderReceivedCallback(callback func(key, value string)) {
	c.network.httpCallbacksMu.Lock()
	if data, ok := c.network.httpCallbacks[c.ptr]; ok {
		data.headerReceived = callback
	}
	c.network.httpCallbacksMu.Unlock()
	C.pd_http_setHeaderReceivedCallback(unsafe.Pointer(c.ptr))
}

// SetHeadersReadCallback sets the callback for when all headers are read
func (c *HTTPConnection) SetHeadersReadCallback(callback func()) {
	c.network.httpCallbacksMu.Lock()
	if data, ok := c.network.httpCallbacks[c.ptr]; ok {
		data.headersRead = callback
	}
	c.network.httpCallbacksMu.Unlock()
	C.pd_http_setHeadersReadCallback(unsafe.Pointer(c.ptr))
}

// SetResponseCallback sets the callback for when response starts
func (c *HTTPConnection) SetResponseCallback(callback func()) {
	c.network.httpCallbacksMu.Lock()
	if data, ok := c.network.httpCallbacks[c.ptr]; ok {
		data.response = callback
	}
	c.network.httpCallbacksMu.Unlock()
	C.pd_http_setResponseCallback(unsafe.Pointer(c.ptr))
}

// SetRequestCompleteCallback sets the callback for when request is complete
func (c *HTTPConnection) SetRequestCompleteCallback(callback func()) {
	c.network.httpCallbacksMu.Lock()
	if data, ok := c.network.httpCallbacks[c.ptr]; ok {
		data.requestComplete = callback
	}
	c.network.httpCallbacksMu.Unlock()
	C.pd_http_setRequestCompleteCallback(unsafe.Pointer(c.ptr))
}

// SetConnectionClosedCallback sets the callback for when connection closes
func (c *HTTPConnection) SetConnectionClosedCallback(callback func()) {
	c.network.httpCallbacksMu.Lock()
	if data, ok := c.network.httpCallbacks[c.ptr]; ok {
		data.connectionClosed = callback
	}
	c.network.httpCallbacksMu.Unlock()
	C.pd_http_setConnectionClosedCallback(unsafe.Pointer(c.ptr))
}

// TCPAPI provides TCP networking functions
type TCPAPI struct {
	network *Network
}

// RequestAccess requests permission to connect to a server.
// Returns AccessReply indicating whether access was granted, denied, or pending.
// If not AccessAsk, the callback will NOT be called and the reply should be used directly.
func (t *TCPAPI) RequestAccess(server string, port int, useSSL bool, purpose string, callback AccessCallback) AccessReply {
	cServer := make([]byte, len(server)+1)
	copy(cServer, server)
	cPurpose := make([]byte, len(purpose)+1)
	copy(cPurpose, purpose)

	var useSSLInt C.int
	if useSSL {
		useSSLInt = 1
	}

	var userdata C.int
	if callback != nil {
		t.network.accessCallbacksMu.Lock()
		t.network.accessCallbackID++
		id := t.network.accessCallbackID
		t.network.accessCallbacks[id] = callback
		t.network.accessCallbacksMu.Unlock()
		userdata = C.int(id)
	}

	reply := C.pd_tcp_requestAccess(
		(*C.char)(unsafe.Pointer(&cServer[0])),
		C.int(port),
		useSSLInt,
		(*C.char)(unsafe.Pointer(&cPurpose[0])),
		unsafe.Pointer(uintptr(userdata)),
	)

	// If access was immediately granted/denied, callback won't be called
	// Clean up the callback and call it now
	if reply != 0 && callback != nil { // not kAccessAsk
		t.network.accessCallbacksMu.Lock()
		delete(t.network.accessCallbacks, int(userdata))
		t.network.accessCallbacksMu.Unlock()
		// Call the callback with the result
		go callback(reply == 2) // kAccessAllow = 2
	}

	return AccessReply(reply)
}

// NewConnection creates a new TCP connection
func (t *TCPAPI) NewConnection(server string, port int, useSSL bool) *TCPConnection {
	cServer := make([]byte, len(server)+1)
	copy(cServer, server)

	var useSSLInt C.int
	if useSSL {
		useSSLInt = 1
	}

	ptr := C.pd_tcp_newConnection(
		(*C.char)(unsafe.Pointer(&cServer[0])),
		C.int(port),
		useSSLInt,
	)

	if ptr == nil {
		return nil
	}

	conn := &TCPConnection{
		ptr:     uintptr(unsafe.Pointer(ptr)),
		network: t.network,
	}

	// Initialize callback data
	t.network.tcpCallbacksMu.Lock()
	t.network.tcpCallbacks[conn.ptr] = &tcpCallbackData{}
	t.network.tcpCallbacksMu.Unlock()

	return conn
}

// Retain increments the reference count
func (c *TCPConnection) Retain() {
	C.pd_tcp_retain(unsafe.Pointer(c.ptr))
}

// Release decrements the reference count and frees when zero
func (c *TCPConnection) Release() {
	C.pd_tcp_release(unsafe.Pointer(c.ptr))
	c.network.tcpCallbacksMu.Lock()
	delete(c.network.tcpCallbacks, c.ptr)
	c.network.tcpCallbacksMu.Unlock()
}

// GetError returns the last error for the connection
func (c *TCPConnection) GetError() PDNetErr {
	return PDNetErr(C.pd_tcp_getError(unsafe.Pointer(c.ptr)))
}

// SetConnectTimeout sets the connection timeout in milliseconds
func (c *TCPConnection) SetConnectTimeout(ms int) {
	C.pd_tcp_setConnectTimeout(unsafe.Pointer(c.ptr), C.int(ms))
}

// Open opens the TCP connection
func (c *TCPConnection) Open(callback func(err PDNetErr)) PDNetErr {
	c.network.tcpCallbacksMu.Lock()
	if data, ok := c.network.tcpCallbacks[c.ptr]; ok {
		data.openCallback = callback
	}
	c.network.tcpCallbacksMu.Unlock()
	return PDNetErr(C.pd_tcp_open(unsafe.Pointer(c.ptr), nil))
}

// Close closes the connection
func (c *TCPConnection) Close() PDNetErr {
	return PDNetErr(C.pd_tcp_close(unsafe.Pointer(c.ptr)))
}

// SetConnectionClosedCallback sets the callback for when connection closes
func (c *TCPConnection) SetConnectionClosedCallback(callback func(err PDNetErr)) {
	c.network.tcpCallbacksMu.Lock()
	if data, ok := c.network.tcpCallbacks[c.ptr]; ok {
		data.closedCallback = callback
	}
	c.network.tcpCallbacksMu.Unlock()
	C.pd_tcp_setConnectionClosedCallback(unsafe.Pointer(c.ptr))
}

// SetReadTimeout sets the read timeout in milliseconds
func (c *TCPConnection) SetReadTimeout(ms int) {
	C.pd_tcp_setReadTimeout(unsafe.Pointer(c.ptr), C.int(ms))
}

// SetReadBufferSize sets the read buffer size
func (c *TCPConnection) SetReadBufferSize(bytes int) {
	C.pd_tcp_setReadBufferSize(unsafe.Pointer(c.ptr), C.int(bytes))
}

// GetBytesAvailable returns the number of bytes available to read
func (c *TCPConnection) GetBytesAvailable() int {
	return int(C.pd_tcp_getBytesAvailable(unsafe.Pointer(c.ptr)))
}

// Read reads data from the connection
func (c *TCPConnection) Read(buf []byte) int {
	if len(buf) == 0 {
		return 0
	}
	return int(C.pd_tcp_read(unsafe.Pointer(c.ptr),
		unsafe.Pointer(&buf[0]), C.int(len(buf))))
}

// Write writes data to the connection
func (c *TCPConnection) Write(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	return int(C.pd_tcp_write(unsafe.Pointer(c.ptr),
		unsafe.Pointer(&data[0]), C.int(len(data))))
}

// C callback trampolines - called from C, dispatch to Go callbacks

//export pdgo_http_access_callback
func pdgo_http_access_callback(allowed C.int, userdata unsafe.Pointer) {
	if api == nil || api.Network == nil {
		return
	}
	id := int(uintptr(userdata))
	api.Network.accessCallbacksMu.Lock()
	callback, ok := api.Network.accessCallbacks[id]
	if ok {
		delete(api.Network.accessCallbacks, id)
	}
	api.Network.accessCallbacksMu.Unlock()

	if callback != nil {
		callback(allowed != 0)
	}
}

//export pdgo_http_header_callback
func pdgo_http_header_callback(conn unsafe.Pointer, key, value *C.char) {
	if api == nil || api.Network == nil {
		return
	}
	ptr := uintptr(conn)
	api.Network.httpCallbacksMu.Lock()
	data, ok := api.Network.httpCallbacks[ptr]
	api.Network.httpCallbacksMu.Unlock()

	if ok && data != nil && data.headerReceived != nil {
		data.headerReceived(C.GoString(key), C.GoString(value))
	}
}

//export pdgo_http_headers_read_callback
func pdgo_http_headers_read_callback(conn unsafe.Pointer) {
	if api == nil || api.Network == nil {
		return
	}
	ptr := uintptr(conn)
	api.Network.httpCallbacksMu.Lock()
	data, ok := api.Network.httpCallbacks[ptr]
	api.Network.httpCallbacksMu.Unlock()

	if ok && data != nil && data.headersRead != nil {
		data.headersRead()
	}
}

//export pdgo_http_response_callback
func pdgo_http_response_callback(conn unsafe.Pointer) {
	if api == nil || api.Network == nil {
		return
	}
	ptr := uintptr(conn)
	api.Network.httpCallbacksMu.Lock()
	data, ok := api.Network.httpCallbacks[ptr]
	api.Network.httpCallbacksMu.Unlock()

	if ok && data != nil && data.response != nil {
		data.response()
	}
}

//export pdgo_http_request_complete_callback
func pdgo_http_request_complete_callback(conn unsafe.Pointer) {
	if api == nil || api.Network == nil {
		return
	}
	ptr := uintptr(conn)
	api.Network.httpCallbacksMu.Lock()
	data, ok := api.Network.httpCallbacks[ptr]
	api.Network.httpCallbacksMu.Unlock()

	if ok && data != nil && data.requestComplete != nil {
		data.requestComplete()
	}
}

//export pdgo_http_connection_closed_callback
func pdgo_http_connection_closed_callback(conn unsafe.Pointer) {
	if api == nil || api.Network == nil {
		return
	}
	ptr := uintptr(conn)
	api.Network.httpCallbacksMu.Lock()
	data, ok := api.Network.httpCallbacks[ptr]
	api.Network.httpCallbacksMu.Unlock()

	if ok && data != nil && data.connectionClosed != nil {
		data.connectionClosed()
	}
}

//export pdgo_tcp_access_callback
func pdgo_tcp_access_callback(allowed C.int, userdata unsafe.Pointer) {
	if api == nil || api.Network == nil {
		return
	}
	id := int(uintptr(userdata))
	api.Network.accessCallbacksMu.Lock()
	callback, ok := api.Network.accessCallbacks[id]
	if ok {
		delete(api.Network.accessCallbacks, id)
	}
	api.Network.accessCallbacksMu.Unlock()

	if callback != nil {
		callback(allowed != 0)
	}
}

//export pdgo_tcp_open_callback
func pdgo_tcp_open_callback(conn unsafe.Pointer, err C.int, userdata unsafe.Pointer) {
	if api == nil || api.Network == nil {
		return
	}
	ptr := uintptr(conn)
	api.Network.tcpCallbacksMu.Lock()
	data, ok := api.Network.tcpCallbacks[ptr]
	api.Network.tcpCallbacksMu.Unlock()

	if ok && data != nil && data.openCallback != nil {
		data.openCallback(PDNetErr(err))
	}
}

//export pdgo_tcp_closed_callback
func pdgo_tcp_closed_callback(conn unsafe.Pointer, err C.int) {
	if api == nil || api.Network == nil {
		return
	}
	ptr := uintptr(conn)
	api.Network.tcpCallbacksMu.Lock()
	data, ok := api.Network.tcpCallbacks[ptr]
	api.Network.tcpCallbacksMu.Unlock()

	if ok && data != nil && data.closedCallback != nil {
		data.closedCallback(PDNetErr(err))
	}
}
