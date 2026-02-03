// pdgo JSON API - Playdate SDK JSON encoder/decoder

package pdgo

/*
// Decoder
int pd_json_decodeString(const char* jsonString);
int pd_json_decodeFile(void* file);

// Encoder
void pd_json_encoder_init(int pretty);
void pd_json_encoder_startArray(void);
void pd_json_encoder_addArrayMember(void);
void pd_json_encoder_endArray(void);
void pd_json_encoder_startTable(void);
void pd_json_encoder_addTableMember(const char* name, int len);
void pd_json_encoder_endTable(void);
void pd_json_encoder_writeNull(void);
void pd_json_encoder_writeFalse(void);
void pd_json_encoder_writeTrue(void);
void pd_json_encoder_writeInt(int num);
void pd_json_encoder_writeDouble(double num);
void pd_json_encoder_writeString(const char* str, int len);
const char* pd_json_encoder_getOutput(void);
int pd_json_encoder_getOutputLen(void);
void pd_json_encoder_free(void);
*/
import "C"
import "unsafe"

// JSON provides access to Playdate JSON parsing
type JSON struct{}

func newJSON() *JSON {
	return &JSON{}
}

// JSONValueType represents JSON value types (matches Playdate SDK)
type JSONValueType int

const (
	JSONNull    JSONValueType = 0
	JSONTrue    JSONValueType = 1
	JSONFalse   JSONValueType = 2
	JSONInteger JSONValueType = 3
	JSONFloat   JSONValueType = 4
	JSONString  JSONValueType = 5
	JSONArray   JSONValueType = 6
	JSONTable   JSONValueType = 7
)

// JSONValue represents a decoded JSON value
type JSONValue struct {
	Type   JSONValueType
	Int    int
	Float  float32
	String string
}

// JSONDecodeHandler handles JSON decode events
type JSONDecodeHandler interface {
	// DecodeError is called when a parse error occurs
	DecodeError(error string, lineNum int)

	// WillDecodeSublist is called before decoding an array or table
	WillDecodeSublist(name string, valueType JSONValueType)

	// ShouldDecodeTableValueForKey returns true to decode this key's value
	ShouldDecodeTableValueForKey(key string) bool

	// DidDecodeTableValue is called after decoding a table value
	DidDecodeTableValue(key string, value JSONValue)

	// ShouldDecodeArrayValueAtIndex returns true to decode this index's value
	ShouldDecodeArrayValueAtIndex(pos int) bool

	// DidDecodeArrayValue is called after decoding an array value
	DidDecodeArrayValue(pos int, value JSONValue)

	// DidDecodeSublist is called after finishing an array or table
	DidDecodeSublist(name string, valueType JSONValueType)
}

// Global handler for callbacks
var jsonHandler JSONDecodeHandler

// DecodeString decodes a JSON string using the provided handler
func (j *JSON) DecodeString(jsonStr string, handler JSONDecodeHandler) error {
	jsonHandler = handler
	cstr := make([]byte, len(jsonStr)+1)
	copy(cstr, jsonStr)
	result := C.pd_json_decodeString((*C.char)(unsafe.Pointer(&cstr[0])))
	jsonHandler = nil
	if result == 0 {
		return &jsonError{op: "decode", msg: "parse failed"}
	}
	return nil
}

// DecodeFile decodes JSON from an open file using the provided handler
func (j *JSON) DecodeFile(file *SDFile, handler JSONDecodeHandler) error {
	if file == nil || file.ptr == nil {
		return &jsonError{op: "decode", msg: "invalid file"}
	}
	jsonHandler = handler
	result := C.pd_json_decodeFile(file.ptr)
	jsonHandler = nil
	if result == 0 {
		return &jsonError{op: "decode", msg: "parse failed"}
	}
	return nil
}

// Go callbacks called from C

//export pdgo_json_decodeError
func pdgo_json_decodeError(cerror *C.char, linenum C.int) {
	if jsonHandler != nil {
		jsonHandler.DecodeError(C.GoString(cerror), int(linenum))
	}
}

//export pdgo_json_willDecodeSublist
func pdgo_json_willDecodeSublist(cname *C.char, ctype C.int) {
	if jsonHandler != nil {
		jsonHandler.WillDecodeSublist(C.GoString(cname), JSONValueType(ctype))
	}
}

//export pdgo_json_shouldDecodeTableValueForKey
func pdgo_json_shouldDecodeTableValueForKey(ckey *C.char) C.int {
	if jsonHandler != nil {
		if jsonHandler.ShouldDecodeTableValueForKey(C.GoString(ckey)) {
			return 1
		}
	}
	return 0
}

//export pdgo_json_didDecodeTableValue
func pdgo_json_didDecodeTableValue(ckey *C.char, ctype C.int, intVal C.int, floatVal C.float, strVal *C.char) {
	if jsonHandler != nil {
		v := JSONValue{Type: JSONValueType(ctype)}
		switch v.Type {
		case JSONInteger:
			v.Int = int(intVal)
		case JSONFloat:
			v.Float = float32(floatVal)
		case JSONString:
			if strVal != nil {
				v.String = C.GoString(strVal)
			}
		case JSONTrue:
			v.Int = 1
		case JSONFalse:
			v.Int = 0
		}
		jsonHandler.DidDecodeTableValue(C.GoString(ckey), v)
	}
}

//export pdgo_json_shouldDecodeArrayValueAtIndex
func pdgo_json_shouldDecodeArrayValueAtIndex(pos C.int) C.int {
	if jsonHandler != nil {
		if jsonHandler.ShouldDecodeArrayValueAtIndex(int(pos)) {
			return 1
		}
	}
	return 0
}

//export pdgo_json_didDecodeArrayValue
func pdgo_json_didDecodeArrayValue(pos C.int, ctype C.int, intVal C.int, floatVal C.float, strVal *C.char) {
	if jsonHandler != nil {
		v := JSONValue{Type: JSONValueType(ctype)}
		switch v.Type {
		case JSONInteger:
			v.Int = int(intVal)
		case JSONFloat:
			v.Float = float32(floatVal)
		case JSONString:
			if strVal != nil {
				v.String = C.GoString(strVal)
			}
		case JSONTrue:
			v.Int = 1
		case JSONFalse:
			v.Int = 0
		}
		jsonHandler.DidDecodeArrayValue(int(pos), v)
	}
}

//export pdgo_json_didDecodeSublist
func pdgo_json_didDecodeSublist(cname *C.char, ctype C.int) {
	if jsonHandler != nil {
		jsonHandler.DidDecodeSublist(C.GoString(cname), JSONValueType(ctype))
	}
}

type jsonError struct {
	op  string
	msg string
}

func (e *jsonError) Error() string {
	return "json " + e.op + ": " + e.msg
}

// ============== High-level API ==============

// JSONNode represents a JSON value in a tree structure
type JSONNode struct {
	Type     JSONValueType
	IntVal   int
	FloatVal float32
	StrVal   string
	Array    []*JSONNode
	Object   map[string]*JSONNode
}

// Get returns child node by key (for objects)
func (n *JSONNode) Get(key string) *JSONNode {
	if n != nil && n.Type == JSONTable && n.Object != nil {
		return n.Object[key]
	}
	return nil
}

// At returns child node by index (for arrays)
func (n *JSONNode) At(index int) *JSONNode {
	if n != nil && n.Type == JSONArray && index >= 0 && index < len(n.Array) {
		return n.Array[index]
	}
	return nil
}

// Len returns length of array or object
func (n *JSONNode) Len() int {
	if n == nil {
		return 0
	}
	if n.Type == JSONArray {
		return len(n.Array)
	}
	if n.Type == JSONTable {
		return len(n.Object)
	}
	return 0
}

// GetString returns string value
func (n *JSONNode) GetString() string {
	if n != nil && n.Type == JSONString {
		return n.StrVal
	}
	return ""
}

// GetInt returns integer value
func (n *JSONNode) GetInt() int {
	if n != nil && n.Type == JSONInteger {
		return n.IntVal
	}
	return 0
}

// GetFloat returns float value
func (n *JSONNode) GetFloat() float32 {
	if n != nil && n.Type == JSONFloat {
		return n.FloatVal
	}
	return 0
}

// GetBool returns boolean value
func (n *JSONNode) GetBool() bool {
	if n != nil {
		if n.Type == JSONTrue {
			return true
		}
	}
	return false
}

// IsNull returns true if node is null
func (n *JSONNode) IsNull() bool {
	return n == nil || n.Type == JSONNull
}

// treeBuilder builds JSONNode tree from callbacks
type treeBuilder struct {
	root    *JSONNode
	stack   []*JSONNode
	current *JSONNode
	err     string
}

func (b *treeBuilder) DecodeError(error string, lineNum int) {
	b.err = error
}

func (b *treeBuilder) WillDecodeSublist(name string, valueType JSONValueType) {
	node := &JSONNode{Type: valueType}
	if valueType == JSONTable {
		node.Object = make(map[string]*JSONNode)
	} else if valueType == JSONArray {
		node.Array = make([]*JSONNode, 0)
	}

	if b.current != nil {
		b.stack = append(b.stack, b.current)
	}
	b.current = node

	if b.root == nil {
		b.root = node
	}
}

func (b *treeBuilder) ShouldDecodeTableValueForKey(key string) bool {
	return true
}

func (b *treeBuilder) DidDecodeTableValue(key string, value JSONValue) {
	if b.current == nil || b.current.Object == nil {
		return
	}
	// Skip if value is table/array - it will be added in DidDecodeSublist
	if value.Type == JSONTable || value.Type == JSONArray {
		return
	}
	node := &JSONNode{
		Type:     value.Type,
		IntVal:   value.Int,
		FloatVal: value.Float,
		StrVal:   value.String,
	}
	b.current.Object[key] = node
}

func (b *treeBuilder) ShouldDecodeArrayValueAtIndex(pos int) bool {
	return true
}

func (b *treeBuilder) DidDecodeArrayValue(pos int, value JSONValue) {
	if b.current == nil {
		return
	}
	// Skip if value is table/array - it will be added in DidDecodeSublist
	if value.Type == JSONTable || value.Type == JSONArray {
		return
	}
	node := &JSONNode{
		Type:     value.Type,
		IntVal:   value.Int,
		FloatVal: value.Float,
		StrVal:   value.String,
	}
	b.current.Array = append(b.current.Array, node)
}

func (b *treeBuilder) DidDecodeSublist(name string, valueType JSONValueType) {
	finished := b.current

	if len(b.stack) > 0 {
		parent := b.stack[len(b.stack)-1]
		b.stack = b.stack[:len(b.stack)-1]

		// Attach finished node to parent using name from callback
		if parent.Type == JSONTable && name != "" && name != "_root" {
			parent.Object[name] = finished
		} else if parent.Type == JSONArray {
			parent.Array = append(parent.Array, finished)
		}
		b.current = parent
	}
}

// Parse parses JSON string and returns tree
func (j *JSON) Parse(jsonStr string) (*JSONNode, error) {
	builder := &treeBuilder{}
	err := j.DecodeString(jsonStr, builder)
	if err != nil {
		return nil, err
	}
	if builder.err != "" {
		return nil, &jsonError{op: "parse", msg: builder.err}
	}
	return unwrapRoot(builder.root), nil
}

// ParseFile parses JSON from file and returns tree
func (j *JSON) ParseFile(file *SDFile) (*JSONNode, error) {
	builder := &treeBuilder{}
	err := j.DecodeFile(file, builder)
	if err != nil {
		return nil, err
	}
	if builder.err != "" {
		return nil, &jsonError{op: "parse", msg: builder.err}
	}
	return unwrapRoot(builder.root), nil
}

// unwrapRoot removes Playdate's _root wrapper if present
// Returns the actual JSON content (not _root itself)
func unwrapRoot(root *JSONNode) *JSONNode {
	// Root from Playdate is always _root table containing actual JSON
	// We return root as-is since _root IS the actual parsed content
	return root
}

// ============== JSON Encoder ==============

// JSONEncoder encodes data to JSON format
type JSONEncoder struct {
	pretty bool
}

// NewEncoder creates a new JSON encoder
// If pretty is true, output will be formatted with indentation
func (j *JSON) NewEncoder(pretty bool) *JSONEncoder {
	C.pd_json_encoder_init(boolToInt(pretty))
	return &JSONEncoder{pretty: pretty}
}

// StartObject starts a JSON object {...}
func (e *JSONEncoder) StartObject() {
	C.pd_json_encoder_startTable()
}

// EndObject ends a JSON object
func (e *JSONEncoder) EndObject() {
	C.pd_json_encoder_endTable()
}

// StartArray starts a JSON array [...]
func (e *JSONEncoder) StartArray() {
	C.pd_json_encoder_startArray()
}

// EndArray ends a JSON array
func (e *JSONEncoder) EndArray() {
	C.pd_json_encoder_endArray()
}

// WriteKey writes an object key (call before writing value)
func (e *JSONEncoder) WriteKey(key string) {
	ckey := make([]byte, len(key)+1)
	copy(ckey, key)
	C.pd_json_encoder_addTableMember((*C.char)(unsafe.Pointer(&ckey[0])), C.int(len(key)))
}

// AddArrayElement prepares for next array element
func (e *JSONEncoder) AddArrayElement() {
	C.pd_json_encoder_addArrayMember()
}

// WriteNull writes null value
func (e *JSONEncoder) WriteNull() {
	C.pd_json_encoder_writeNull()
}

// WriteBool writes boolean value
func (e *JSONEncoder) WriteBool(v bool) {
	if v {
		C.pd_json_encoder_writeTrue()
	} else {
		C.pd_json_encoder_writeFalse()
	}
}

// WriteInt writes integer value
func (e *JSONEncoder) WriteInt(v int) {
	C.pd_json_encoder_writeInt(C.int(v))
}

// WriteFloat writes float value
func (e *JSONEncoder) WriteFloat(v float64) {
	C.pd_json_encoder_writeDouble(C.double(v))
}

// WriteString writes string value
func (e *JSONEncoder) WriteString(s string) {
	cstr := make([]byte, len(s)+1)
	copy(cstr, s)
	C.pd_json_encoder_writeString((*C.char)(unsafe.Pointer(&cstr[0])), C.int(len(s)))
}

// String returns the encoded JSON string
func (e *JSONEncoder) String() string {
	cstr := C.pd_json_encoder_getOutput()
	return C.GoString(cstr)
}

// Bytes returns the encoded JSON as bytes
func (e *JSONEncoder) Bytes() []byte {
	cstr := C.pd_json_encoder_getOutput()
	length := int(C.pd_json_encoder_getOutputLen())
	if length == 0 {
		return nil
	}
	result := make([]byte, length)
	copy(result, C.GoStringN(cstr, C.int(length)))
	return result
}

// Free releases encoder resources
func (e *JSONEncoder) Free() {
	C.pd_json_encoder_free()
}

func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
