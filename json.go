//go:build !tinygo

package pdgo

/*
#include "pd_api.h"
#include <stdlib.h>
#include <string.h>

// JSON API helper functions

static void json_initEncoder(const struct playdate_json* json, json_encoder* encoder, json_writeFunc* write, void* userdata, int pretty) {
    json->initEncoder(encoder, write, userdata, pretty);
}

static int json_decode(const struct playdate_json* json, struct json_decoder* functions, json_reader reader, json_value* outval) {
    return json->decode(functions, reader, outval);
}

static int json_decodeString(const struct playdate_json* json, struct json_decoder* functions, const char* jsonString, json_value* outval) {
    return json->decodeString(functions, jsonString, outval);
}

// JSON encoder wrapper functions
typedef struct {
    char* buffer;
    size_t size;
    size_t capacity;
} StringBuffer;

static void stringBufferWrite(void* userdata, const char* str, int len) {
    StringBuffer* buf = (StringBuffer*)userdata;
    if (buf->size + len >= buf->capacity) {
        size_t newCapacity = (buf->capacity == 0) ? 256 : buf->capacity * 2;
        while (newCapacity < buf->size + len + 1) {
            newCapacity *= 2;
        }
        char* newBuffer = realloc(buf->buffer, newCapacity);
        if (newBuffer == NULL) return;
        buf->buffer = newBuffer;
        buf->capacity = newCapacity;
    }
    memcpy(buf->buffer + buf->size, str, len);
    buf->size += len;
    buf->buffer[buf->size] = '\0';
}

static json_encoder* json_newEncoder(const struct playdate_json* json, StringBuffer* buf, int pretty) {
    json_encoder* encoder = malloc(sizeof(json_encoder));
    if (encoder == NULL) return NULL;
    json->initEncoder(encoder, stringBufferWrite, buf, pretty);
    return encoder;
}

static void json_freeEncoder(json_encoder* encoder) {
    free(encoder);
}

static void json_encoder_startArray(json_encoder* encoder) {
    encoder->startArray(encoder);
}

static void json_encoder_addArrayMember(json_encoder* encoder) {
    encoder->addArrayMember(encoder);
}

static void json_encoder_endArray(json_encoder* encoder) {
    encoder->endArray(encoder);
}

static void json_encoder_startTable(json_encoder* encoder) {
    encoder->startTable(encoder);
}

static void json_encoder_addTableMember(json_encoder* encoder, const char* name, int len) {
    encoder->addTableMember(encoder, name, len);
}

static void json_encoder_endTable(json_encoder* encoder) {
    encoder->endTable(encoder);
}

static void json_encoder_writeNull(json_encoder* encoder) {
    encoder->writeNull(encoder);
}

static void json_encoder_writeFalse(json_encoder* encoder) {
    encoder->writeFalse(encoder);
}

static void json_encoder_writeTrue(json_encoder* encoder) {
    encoder->writeTrue(encoder);
}

static void json_encoder_writeInt(json_encoder* encoder, int num) {
    encoder->writeInt(encoder, num);
}

static void json_encoder_writeDouble(json_encoder* encoder, double num) {
    encoder->writeDouble(encoder, num);
}

static void json_encoder_writeString(json_encoder* encoder, const char* str, int len) {
    encoder->writeString(encoder, str, len);
}

static StringBuffer* stringBuffer_new() {
    StringBuffer* buf = malloc(sizeof(StringBuffer));
    if (buf == NULL) return NULL;
    buf->buffer = NULL;
    buf->size = 0;
    buf->capacity = 0;
    return buf;
}

static void stringBuffer_free(StringBuffer* buf) {
    if (buf != NULL) {
        free(buf->buffer);
        free(buf);
    }
}

static const char* stringBuffer_getString(StringBuffer* buf) {
    return buf->buffer;
}

static size_t stringBuffer_getSize(StringBuffer* buf) {
    return buf->size;
}
*/
import "C"
import "unsafe"

// JSONValueType represents JSON value types
type JSONValueType int

const (
	JSONNull    JSONValueType = C.kJSONNull
	JSONTrue    JSONValueType = C.kJSONTrue
	JSONFalse   JSONValueType = C.kJSONFalse
	JSONInteger JSONValueType = C.kJSONInteger
	JSONFloat   JSONValueType = C.kJSONFloat
	JSONString  JSONValueType = C.kJSONString
	JSONArray   JSONValueType = C.kJSONArray
	JSONTable   JSONValueType = C.kJSONTable
)

// JSONValue represents a JSON value
type JSONValue struct {
	Type        JSONValueType
	IntValue    int
	FloatValue  float32
	StringValue string
}

// JSON wraps the playdate_json API
type JSON struct {
	ptr *C.struct_playdate_json
}

func newJSON(ptr *C.struct_playdate_json) *JSON {
	return &JSON{ptr: ptr}
}

// JSONEncoder wraps a JSON encoder
type JSONEncoder struct {
	json    *C.struct_playdate_json
	encoder *C.json_encoder
	buffer  *C.StringBuffer
}

// NewEncoder creates a new JSON encoder
func (j *JSON) NewEncoder(pretty bool) *JSONEncoder {
	buf := C.stringBuffer_new()
	if buf == nil {
		return nil
	}

	p := 0
	if pretty {
		p = 1
	}

	encoder := C.json_newEncoder(j.ptr, buf, C.int(p))
	if encoder == nil {
		C.stringBuffer_free(buf)
		return nil
	}

	return &JSONEncoder{
		json:    j.ptr,
		encoder: encoder,
		buffer:  buf,
	}
}

// Free frees the encoder
func (e *JSONEncoder) Free() {
	if e.encoder != nil {
		C.json_freeEncoder(e.encoder)
		e.encoder = nil
	}
	if e.buffer != nil {
		C.stringBuffer_free(e.buffer)
		e.buffer = nil
	}
}

// GetString returns the encoded JSON string
func (e *JSONEncoder) GetString() string {
	if e.buffer == nil {
		return ""
	}
	s := C.stringBuffer_getString(e.buffer)
	if s == nil {
		return ""
	}
	return C.GoStringN(s, C.int(C.stringBuffer_getSize(e.buffer)))
}

// StartArray starts encoding an array
func (e *JSONEncoder) StartArray() {
	if e.encoder != nil {
		C.json_encoder_startArray(e.encoder)
	}
}

// AddArrayMember adds an array member separator
func (e *JSONEncoder) AddArrayMember() {
	if e.encoder != nil {
		C.json_encoder_addArrayMember(e.encoder)
	}
}

// EndArray ends encoding an array
func (e *JSONEncoder) EndArray() {
	if e.encoder != nil {
		C.json_encoder_endArray(e.encoder)
	}
}

// StartTable starts encoding a table/object
func (e *JSONEncoder) StartTable() {
	if e.encoder != nil {
		C.json_encoder_startTable(e.encoder)
	}
}

// AddTableMember adds a table member
func (e *JSONEncoder) AddTableMember(name string) {
	if e.encoder != nil {
		cname := cString(name)
		defer freeCString(cname)
		C.json_encoder_addTableMember(e.encoder, cname, C.int(len(name)))
	}
}

// EndTable ends encoding a table/object
func (e *JSONEncoder) EndTable() {
	if e.encoder != nil {
		C.json_encoder_endTable(e.encoder)
	}
}

// WriteNull writes a null value
func (e *JSONEncoder) WriteNull() {
	if e.encoder != nil {
		C.json_encoder_writeNull(e.encoder)
	}
}

// WriteBool writes a boolean value
func (e *JSONEncoder) WriteBool(val bool) {
	if e.encoder != nil {
		if val {
			C.json_encoder_writeTrue(e.encoder)
		} else {
			C.json_encoder_writeFalse(e.encoder)
		}
	}
}

// WriteInt writes an integer value
func (e *JSONEncoder) WriteInt(val int) {
	if e.encoder != nil {
		C.json_encoder_writeInt(e.encoder, C.int(val))
	}
}

// WriteFloat writes a float value
func (e *JSONEncoder) WriteFloat(val float64) {
	if e.encoder != nil {
		C.json_encoder_writeDouble(e.encoder, C.double(val))
	}
}

// WriteString writes a string value
func (e *JSONEncoder) WriteString(val string) {
	if e.encoder != nil {
		cstr := cString(val)
		defer freeCString(cstr)
		C.json_encoder_writeString(e.encoder, cstr, C.int(len(val)))
	}
}

// DecodeString decodes a JSON string
// Note: This is a simplified implementation. Full JSON decoding requires
// setting up decoder callbacks which is complex in Go-C interop.
func (j *JSON) DecodeString(jsonString string) (*JSONValue, error) {
	// For complex JSON decoding, users should implement their own decoder
	// using the C callbacks. This method provides basic functionality.
	cstr := cString(jsonString)
	defer freeCString(cstr)

	var outval C.json_value
	var decoder C.struct_json_decoder
	decoder.decodeError = nil
	decoder.willDecodeSublist = nil
	decoder.shouldDecodeTableValueForKey = nil
	decoder.didDecodeTableValue = nil
	decoder.shouldDecodeArrayValueAtIndex = nil
	decoder.didDecodeArrayValue = nil
	decoder.didDecodeSublist = nil
	decoder.userdata = nil
	decoder.returnString = 0
	decoder.path = nil

	result := C.json_decodeString(j.ptr, &decoder, cstr, &outval)
	if result == 0 {
		return nil, nil
	}

	val := &JSONValue{
		Type: JSONValueType(outval._type),
	}

	switch val.Type {
	case JSONInteger:
		val.IntValue = int(*(*C.int)(unsafe.Pointer(&outval.data)))
	case JSONFloat:
		val.FloatValue = float32(*(*C.float)(unsafe.Pointer(&outval.data)))
	case JSONString:
		s := *(**C.char)(unsafe.Pointer(&outval.data))
		if s != nil {
			val.StringValue = goString(s)
		}
	case JSONTrue:
		val.IntValue = 1
	case JSONFalse:
		val.IntValue = 0
	}

	return val, nil
}
