//go:build tinygo

// TinyGo implementation of JSON API

package pdgo

// JSON provides access to JSON parsing
type JSON struct{}

func newJSON() *JSON {
	return &JSON{}
}

// JSONDecoder represents a JSON decoder
type JSONDecoder struct {
	ptr uintptr
}

// JSONEncoder represents a JSON encoder
type JSONEncoder struct {
	ptr uintptr
}

// JSONValue represents a JSON value
type JSONValue struct {
	Type  JSONValueType
	Int   int
	Float float32
	Str   string
	Bool  bool
}

// JSONValueType represents JSON value types
type JSONValueType int32

const (
	JSONNull   JSONValueType = 0
	JSONTrue   JSONValueType = 1
	JSONFalse  JSONValueType = 2
	JSONInt    JSONValueType = 3
	JSONFloat  JSONValueType = 4
	JSONString JSONValueType = 5
	JSONArray  JSONValueType = 6
	JSONTable  JSONValueType = 7
)

// DecodeString decodes a JSON string
func (j *JSON) DecodeString(jsonStr string) (*JSONValue, error) {
	// Simplified implementation - full JSON parsing requires callbacks
	if bridgeJSONDecode != nil {
		buf := make([]byte, len(jsonStr)+1)
		copy(buf, jsonStr)
		// This is a simplified stub - real implementation needs callbacks
		_ = buf
	}
	return nil, &jsonError{op: "decode", msg: "simplified implementation"}
}

// EncodeToString encodes a value to JSON string (stub)
func (j *JSON) EncodeToString(value interface{}) (string, error) {
	// Not fully implemented - would need reflection or type switching
	return "", &jsonError{op: "encode", msg: "not implemented"}
}

type jsonError struct {
	op  string
	msg string
}

func (e *jsonError) Error() string {
	return "json " + e.op + ": " + e.msg
}
