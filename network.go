//go:build !tinygo

package pdgo

/*
#cgo CFLAGS: -DTARGET_EXTENSION=1
#include "pd_api.h"
#include <stdlib.h>

// Network API helper functions

// HTTP functions
static enum accessReply net_http_requestAccess(const struct playdate_network* net, const char* server, int port, bool usessl, const char* purpose) {
    // Note: Callback is not supported in this simplified binding
    return net->http->requestAccess(server, port, usessl, purpose, NULL, NULL);
}

static HTTPConnection* net_http_newConnection(const struct playdate_network* net, const char* server, int port, bool usessl) {
    return net->http->newConnection(server, port, usessl);
}

static HTTPConnection* net_http_retain(const struct playdate_network* net, HTTPConnection* http) {
    return net->http->retain(http);
}

static void net_http_release(const struct playdate_network* net, HTTPConnection* http) {
    net->http->release(http);
}

static void net_http_setConnectTimeout(const struct playdate_network* net, HTTPConnection* conn, int ms) {
    net->http->setConnectTimeout(conn, ms);
}

static void net_http_setKeepAlive(const struct playdate_network* net, HTTPConnection* conn, bool keepalive) {
    net->http->setKeepAlive(conn, keepalive);
}

static void net_http_setByteRange(const struct playdate_network* net, HTTPConnection* conn, int start, int end) {
    net->http->setByteRange(conn, start, end);
}

static void net_http_setUserdata(const struct playdate_network* net, HTTPConnection* conn, void* userdata) {
    net->http->setUserdata(conn, userdata);
}

static void* net_http_getUserdata(const struct playdate_network* net, HTTPConnection* conn) {
    return net->http->getUserdata(conn);
}

static PDNetErr net_http_get(const struct playdate_network* net, HTTPConnection* conn, const char* path, const char* headers, size_t headerlen) {
    return net->http->get(conn, path, headers, headerlen);
}

static PDNetErr net_http_post(const struct playdate_network* net, HTTPConnection* conn, const char* path, const char* headers, size_t headerlen, const char* body, size_t bodylen) {
    return net->http->post(conn, path, headers, headerlen, body, bodylen);
}

static PDNetErr net_http_query(const struct playdate_network* net, HTTPConnection* conn, const char* method, const char* path, const char* headers, size_t headerlen, const char* body, size_t bodylen) {
    return net->http->query(conn, method, path, headers, headerlen, body, bodylen);
}

static PDNetErr net_http_getError(const struct playdate_network* net, HTTPConnection* conn) {
    return net->http->getError(conn);
}

static void net_http_getProgress(const struct playdate_network* net, HTTPConnection* conn, int* read, int* total) {
    net->http->getProgress(conn, read, total);
}

static int net_http_getResponseStatus(const struct playdate_network* net, HTTPConnection* conn) {
    return net->http->getResponseStatus(conn);
}

static size_t net_http_getBytesAvailable(const struct playdate_network* net, HTTPConnection* conn) {
    return net->http->getBytesAvailable(conn);
}

static void net_http_setReadTimeout(const struct playdate_network* net, HTTPConnection* conn, int ms) {
    net->http->setReadTimeout(conn, ms);
}

static void net_http_setReadBufferSize(const struct playdate_network* net, HTTPConnection* conn, int bytes) {
    net->http->setReadBufferSize(conn, bytes);
}

static int net_http_read(const struct playdate_network* net, HTTPConnection* conn, void* buf, unsigned int buflen) {
    return net->http->read(conn, buf, buflen);
}

static void net_http_close(const struct playdate_network* net, HTTPConnection* conn) {
    net->http->close(conn);
}

// TCP functions
static enum accessReply net_tcp_requestAccess(const struct playdate_network* net, const char* server, int port, bool usessl, const char* purpose) {
    return net->tcp->requestAccess(server, port, usessl, purpose, NULL, NULL);
}

static TCPConnection* net_tcp_newConnection(const struct playdate_network* net, const char* server, int port, bool usessl) {
    return net->tcp->newConnection(server, port, usessl);
}

static TCPConnection* net_tcp_retain(const struct playdate_network* net, TCPConnection* tcp) {
    return net->tcp->retain(tcp);
}

static void net_tcp_release(const struct playdate_network* net, TCPConnection* tcp) {
    net->tcp->release(tcp);
}

static PDNetErr net_tcp_getError(const struct playdate_network* net, TCPConnection* conn) {
    return net->tcp->getError(conn);
}

static void net_tcp_setConnectTimeout(const struct playdate_network* net, TCPConnection* conn, int ms) {
    net->tcp->setConnectTimeout(conn, ms);
}

static void net_tcp_setUserdata(const struct playdate_network* net, TCPConnection* conn, void* userdata) {
    net->tcp->setUserdata(conn, userdata);
}

static void* net_tcp_getUserdata(const struct playdate_network* net, TCPConnection* conn) {
    return net->tcp->getUserdata(conn);
}

static PDNetErr net_tcp_open(const struct playdate_network* net, TCPConnection* conn) {
    return net->tcp->open(conn, NULL, NULL);
}

static PDNetErr net_tcp_close(const struct playdate_network* net, TCPConnection* conn) {
    return net->tcp->close(conn);
}

static void net_tcp_setReadTimeout(const struct playdate_network* net, TCPConnection* conn, int ms) {
    net->tcp->setReadTimeout(conn, ms);
}

static void net_tcp_setReadBufferSize(const struct playdate_network* net, TCPConnection* conn, int bytes) {
    net->tcp->setReadBufferSize(conn, bytes);
}

static size_t net_tcp_getBytesAvailable(const struct playdate_network* net, TCPConnection* conn) {
    return net->tcp->getBytesAvailable(conn);
}

static int net_tcp_read(const struct playdate_network* net, TCPConnection* conn, void* buffer, size_t length) {
    return net->tcp->read(conn, buffer, length);
}

static int net_tcp_write(const struct playdate_network* net, TCPConnection* conn, const void* buffer, size_t length) {
    return net->tcp->write(conn, buffer, length);
}

// Global network functions
static WifiStatus net_getStatus(const struct playdate_network* net) {
    return net->getStatus();
}

static void net_setEnabled(const struct playdate_network* net, bool flag) {
    net->setEnabled(flag, NULL);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// PDNetErr represents network errors
type PDNetErr int

const (
	NetOK                 PDNetErr = C.NET_OK
	NetNoDevice           PDNetErr = C.NET_NO_DEVICE
	NetBusy               PDNetErr = C.NET_BUSY
	NetWriteError         PDNetErr = C.NET_WRITE_ERROR
	NetWriteBusy          PDNetErr = C.NET_WRITE_BUSY
	NetWriteTimeout       PDNetErr = C.NET_WRITE_TIMEOUT
	NetReadError          PDNetErr = C.NET_READ_ERROR
	NetReadBusy           PDNetErr = C.NET_READ_BUSY
	NetReadTimeout        PDNetErr = C.NET_READ_TIMEOUT
	NetReadOverflow       PDNetErr = C.NET_READ_OVERFLOW
	NetFrameError         PDNetErr = C.NET_FRAME_ERROR
	NetBadResponse        PDNetErr = C.NET_BAD_RESPONSE
	NetErrorResponse      PDNetErr = C.NET_ERROR_RESPONSE
	NetResetTimeout       PDNetErr = C.NET_RESET_TIMEOUT
	NetBufferTooSmall     PDNetErr = C.NET_BUFFER_TOO_SMALL
	NetUnexpectedResponse PDNetErr = C.NET_UNEXPECTED_RESPONSE
	NetNotConnectedToAP   PDNetErr = C.NET_NOT_CONNECTED_TO_AP
	NetNotImplemented     PDNetErr = C.NET_NOT_IMPLEMENTED
	NetConnectionClosed   PDNetErr = C.NET_CONNECTION_CLOSED
)

func (e PDNetErr) Error() string {
	switch e {
	case NetOK:
		return "OK"
	case NetNoDevice:
		return "no device"
	case NetBusy:
		return "busy"
	case NetWriteError:
		return "write error"
	case NetWriteBusy:
		return "write busy"
	case NetWriteTimeout:
		return "write timeout"
	case NetReadError:
		return "read error"
	case NetReadBusy:
		return "read busy"
	case NetReadTimeout:
		return "read timeout"
	case NetReadOverflow:
		return "read overflow"
	case NetFrameError:
		return "frame error"
	case NetBadResponse:
		return "bad response"
	case NetErrorResponse:
		return "error response"
	case NetResetTimeout:
		return "reset timeout"
	case NetBufferTooSmall:
		return "buffer too small"
	case NetUnexpectedResponse:
		return "unexpected response"
	case NetNotConnectedToAP:
		return "not connected to AP"
	case NetNotImplemented:
		return "not implemented"
	case NetConnectionClosed:
		return "connection closed"
	default:
		return "unknown error"
	}
}

// WifiStatus represents WiFi connection status
type WifiStatus int

const (
	WifiNotConnected WifiStatus = C.kWifiNotConnected
	WifiConnected    WifiStatus = C.kWifiConnected
	WifiNotAvailable WifiStatus = C.kWifiNotAvailable
)

// AccessReply represents access request response
type AccessReply int

const (
	AccessAsk   AccessReply = C.kAccessAsk
	AccessDeny  AccessReply = C.kAccessDeny
	AccessAllow AccessReply = C.kAccessAllow
)

// HTTPConnection wraps an HTTP connection
type HTTPConnection struct {
	ptr *C.HTTPConnection
	net *C.struct_playdate_network
}

// TCPConnection wraps a TCP connection
type TCPConnection struct {
	ptr *C.TCPConnection
	net *C.struct_playdate_network
}

// Network wraps the playdate_network API
type Network struct {
	ptr  *C.struct_playdate_network
	HTTP *HTTPAPI
	TCP  *TCPAPI
}

func newNetwork(ptr *C.struct_playdate_network) *Network {
	n := &Network{ptr: ptr}
	n.HTTP = &HTTPAPI{net: ptr}
	n.TCP = &TCPAPI{net: ptr}
	return n
}

// GetStatus returns the WiFi status
func (n *Network) GetStatus() WifiStatus {
	return WifiStatus(C.net_getStatus(n.ptr))
}

// SetEnabled enables or disables networking
func (n *Network) SetEnabled(enabled bool) {
	C.net_setEnabled(n.ptr, C.bool(enabled))
}

// HTTPAPI wraps HTTP functions
type HTTPAPI struct {
	net *C.struct_playdate_network
}

// RequestAccess requests network access
func (h *HTTPAPI) RequestAccess(server string, port int, useSSL bool, purpose string) AccessReply {
	cserver := cString(server)
	defer freeCString(cserver)
	cpurpose := cString(purpose)
	defer freeCString(cpurpose)

	return AccessReply(C.net_http_requestAccess(h.net, cserver, C.int(port), C.bool(useSSL), cpurpose))
}

// NewConnection creates a new HTTP connection
func (h *HTTPAPI) NewConnection(server string, port int, useSSL bool) *HTTPConnection {
	cserver := cString(server)
	defer freeCString(cserver)

	ptr := C.net_http_newConnection(h.net, cserver, C.int(port), C.bool(useSSL))
	if ptr == nil {
		return nil
	}
	return &HTTPConnection{ptr: ptr, net: h.net}
}

// Retain retains an HTTP connection
func (h *HTTPAPI) Retain(conn *HTTPConnection) *HTTPConnection {
	if conn == nil {
		return nil
	}
	ptr := C.net_http_retain(h.net, conn.ptr)
	if ptr == nil {
		return nil
	}
	return &HTTPConnection{ptr: ptr, net: h.net}
}

// Release releases an HTTP connection
func (h *HTTPAPI) Release(conn *HTTPConnection) {
	if conn != nil && conn.ptr != nil {
		C.net_http_release(h.net, conn.ptr)
		conn.ptr = nil
	}
}

// SetConnectTimeout sets the connection timeout
func (conn *HTTPConnection) SetConnectTimeout(ms int) {
	if conn != nil && conn.ptr != nil {
		C.net_http_setConnectTimeout(conn.net, conn.ptr, C.int(ms))
	}
}

// SetKeepAlive sets keep-alive
func (conn *HTTPConnection) SetKeepAlive(keepAlive bool) {
	if conn != nil && conn.ptr != nil {
		C.net_http_setKeepAlive(conn.net, conn.ptr, C.bool(keepAlive))
	}
}

// SetByteRange sets the byte range
func (conn *HTTPConnection) SetByteRange(start, end int) {
	if conn != nil && conn.ptr != nil {
		C.net_http_setByteRange(conn.net, conn.ptr, C.int(start), C.int(end))
	}
}

// SetUserdata sets user data
func (conn *HTTPConnection) SetUserdata(userdata unsafe.Pointer) {
	if conn != nil && conn.ptr != nil {
		C.net_http_setUserdata(conn.net, conn.ptr, userdata)
	}
}

// GetUserdata returns user data
func (conn *HTTPConnection) GetUserdata() unsafe.Pointer {
	if conn == nil || conn.ptr == nil {
		return nil
	}
	return C.net_http_getUserdata(conn.net, conn.ptr)
}

// Get performs a GET request
func (conn *HTTPConnection) Get(path string, headers string) error {
	if conn == nil || conn.ptr == nil {
		return errors.New("connection is nil")
	}

	cpath := cString(path)
	defer freeCString(cpath)

	var cheaders *C.char
	var headerLen C.size_t
	if headers != "" {
		cheaders = cString(headers)
		defer freeCString(cheaders)
		headerLen = C.size_t(len(headers))
	}

	err := PDNetErr(C.net_http_get(conn.net, conn.ptr, cpath, cheaders, headerLen))
	if err != NetOK {
		return err
	}
	return nil
}

// Post performs a POST request
func (conn *HTTPConnection) Post(path string, headers string, body []byte) error {
	if conn == nil || conn.ptr == nil {
		return errors.New("connection is nil")
	}

	cpath := cString(path)
	defer freeCString(cpath)

	var cheaders *C.char
	var headerLen C.size_t
	if headers != "" {
		cheaders = cString(headers)
		defer freeCString(cheaders)
		headerLen = C.size_t(len(headers))
	}

	var cbody *C.char
	var bodyLen C.size_t
	if len(body) > 0 {
		cbody = (*C.char)(unsafe.Pointer(&body[0]))
		bodyLen = C.size_t(len(body))
	}

	err := PDNetErr(C.net_http_post(conn.net, conn.ptr, cpath, cheaders, headerLen, cbody, bodyLen))
	if err != NetOK {
		return err
	}
	return nil
}

// Query performs a custom HTTP request
func (conn *HTTPConnection) Query(method, path string, headers string, body []byte) error {
	if conn == nil || conn.ptr == nil {
		return errors.New("connection is nil")
	}

	cmethod := cString(method)
	defer freeCString(cmethod)
	cpath := cString(path)
	defer freeCString(cpath)

	var cheaders *C.char
	var headerLen C.size_t
	if headers != "" {
		cheaders = cString(headers)
		defer freeCString(cheaders)
		headerLen = C.size_t(len(headers))
	}

	var cbody *C.char
	var bodyLen C.size_t
	if len(body) > 0 {
		cbody = (*C.char)(unsafe.Pointer(&body[0]))
		bodyLen = C.size_t(len(body))
	}

	err := PDNetErr(C.net_http_query(conn.net, conn.ptr, cmethod, cpath, cheaders, headerLen, cbody, bodyLen))
	if err != NetOK {
		return err
	}
	return nil
}

// GetError returns the connection error
func (conn *HTTPConnection) GetError() PDNetErr {
	if conn == nil || conn.ptr == nil {
		return NetNoDevice
	}
	return PDNetErr(C.net_http_getError(conn.net, conn.ptr))
}

// GetProgress returns read progress
func (conn *HTTPConnection) GetProgress() (read, total int) {
	if conn == nil || conn.ptr == nil {
		return 0, 0
	}
	var r, t C.int
	C.net_http_getProgress(conn.net, conn.ptr, &r, &t)
	return int(r), int(t)
}

// GetResponseStatus returns the HTTP response status
func (conn *HTTPConnection) GetResponseStatus() int {
	if conn == nil || conn.ptr == nil {
		return 0
	}
	return int(C.net_http_getResponseStatus(conn.net, conn.ptr))
}

// GetBytesAvailable returns available bytes
func (conn *HTTPConnection) GetBytesAvailable() int {
	if conn == nil || conn.ptr == nil {
		return 0
	}
	return int(C.net_http_getBytesAvailable(conn.net, conn.ptr))
}

// SetReadTimeout sets the read timeout
func (conn *HTTPConnection) SetReadTimeout(ms int) {
	if conn != nil && conn.ptr != nil {
		C.net_http_setReadTimeout(conn.net, conn.ptr, C.int(ms))
	}
}

// SetReadBufferSize sets the read buffer size
func (conn *HTTPConnection) SetReadBufferSize(bytes int) {
	if conn != nil && conn.ptr != nil {
		C.net_http_setReadBufferSize(conn.net, conn.ptr, C.int(bytes))
	}
}

// Read reads data from the connection
func (conn *HTTPConnection) Read(buf []byte) (int, error) {
	if conn == nil || conn.ptr == nil {
		return 0, errors.New("connection is nil")
	}
	if len(buf) == 0 {
		return 0, nil
	}

	n := C.net_http_read(conn.net, conn.ptr, unsafe.Pointer(&buf[0]), C.uint(len(buf)))
	if n < 0 {
		return 0, PDNetErr(n)
	}
	return int(n), nil
}

// Close closes the connection
func (conn *HTTPConnection) Close() {
	if conn != nil && conn.ptr != nil {
		C.net_http_close(conn.net, conn.ptr)
	}
}

// TCPAPI wraps TCP functions
type TCPAPI struct {
	net *C.struct_playdate_network
}

// RequestAccess requests network access
func (t *TCPAPI) RequestAccess(server string, port int, useSSL bool, purpose string) AccessReply {
	cserver := cString(server)
	defer freeCString(cserver)
	cpurpose := cString(purpose)
	defer freeCString(cpurpose)

	return AccessReply(C.net_tcp_requestAccess(t.net, cserver, C.int(port), C.bool(useSSL), cpurpose))
}

// NewConnection creates a new TCP connection
func (t *TCPAPI) NewConnection(server string, port int, useSSL bool) *TCPConnection {
	cserver := cString(server)
	defer freeCString(cserver)

	ptr := C.net_tcp_newConnection(t.net, cserver, C.int(port), C.bool(useSSL))
	if ptr == nil {
		return nil
	}
	return &TCPConnection{ptr: ptr, net: t.net}
}

// Retain retains a TCP connection
func (t *TCPAPI) Retain(conn *TCPConnection) *TCPConnection {
	if conn == nil {
		return nil
	}
	ptr := C.net_tcp_retain(t.net, conn.ptr)
	if ptr == nil {
		return nil
	}
	return &TCPConnection{ptr: ptr, net: t.net}
}

// Release releases a TCP connection
func (t *TCPAPI) Release(conn *TCPConnection) {
	if conn != nil && conn.ptr != nil {
		C.net_tcp_release(t.net, conn.ptr)
		conn.ptr = nil
	}
}

// SetConnectTimeout sets the connection timeout
func (conn *TCPConnection) SetConnectTimeout(ms int) {
	if conn != nil && conn.ptr != nil {
		C.net_tcp_setConnectTimeout(conn.net, conn.ptr, C.int(ms))
	}
}

// SetUserdata sets user data
func (conn *TCPConnection) SetUserdata(userdata unsafe.Pointer) {
	if conn != nil && conn.ptr != nil {
		C.net_tcp_setUserdata(conn.net, conn.ptr, userdata)
	}
}

// GetUserdata returns user data
func (conn *TCPConnection) GetUserdata() unsafe.Pointer {
	if conn == nil || conn.ptr == nil {
		return nil
	}
	return C.net_tcp_getUserdata(conn.net, conn.ptr)
}

// GetError returns the connection error
func (conn *TCPConnection) GetError() PDNetErr {
	if conn == nil || conn.ptr == nil {
		return NetNoDevice
	}
	return PDNetErr(C.net_tcp_getError(conn.net, conn.ptr))
}

// Open opens the connection
func (conn *TCPConnection) Open() error {
	if conn == nil || conn.ptr == nil {
		return errors.New("connection is nil")
	}
	err := PDNetErr(C.net_tcp_open(conn.net, conn.ptr))
	if err != NetOK {
		return err
	}
	return nil
}

// Close closes the connection
func (conn *TCPConnection) Close() error {
	if conn == nil || conn.ptr == nil {
		return nil
	}
	err := PDNetErr(C.net_tcp_close(conn.net, conn.ptr))
	if err != NetOK {
		return err
	}
	return nil
}

// SetReadTimeout sets the read timeout
func (conn *TCPConnection) SetReadTimeout(ms int) {
	if conn != nil && conn.ptr != nil {
		C.net_tcp_setReadTimeout(conn.net, conn.ptr, C.int(ms))
	}
}

// SetReadBufferSize sets the read buffer size
func (conn *TCPConnection) SetReadBufferSize(bytes int) {
	if conn != nil && conn.ptr != nil {
		C.net_tcp_setReadBufferSize(conn.net, conn.ptr, C.int(bytes))
	}
}

// GetBytesAvailable returns available bytes
func (conn *TCPConnection) GetBytesAvailable() int {
	if conn == nil || conn.ptr == nil {
		return 0
	}
	return int(C.net_tcp_getBytesAvailable(conn.net, conn.ptr))
}

// Read reads data from the connection
func (conn *TCPConnection) Read(buf []byte) (int, error) {
	if conn == nil || conn.ptr == nil {
		return 0, errors.New("connection is nil")
	}
	if len(buf) == 0 {
		return 0, nil
	}

	n := C.net_tcp_read(conn.net, conn.ptr, unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	if n < 0 {
		return 0, PDNetErr(n)
	}
	return int(n), nil
}

// Write writes data to the connection
func (conn *TCPConnection) Write(buf []byte) (int, error) {
	if conn == nil || conn.ptr == nil {
		return 0, errors.New("connection is nil")
	}
	if len(buf) == 0 {
		return 0, nil
	}

	n := C.net_tcp_write(conn.net, conn.ptr, unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	if n < 0 {
		return 0, PDNetErr(n)
	}
	return int(n), nil
}
