package main

import (
	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

// JSONHandler implements pdgo.JSONDecodeHandler
type JSONHandler struct {
	path string
}

func (h *JSONHandler) DecodeError(error string, lineNum int) {
	pd.System.LogToConsole("decode error: " + error)
}

func (h *JSONHandler) WillDecodeSublist(name string, valueType pdgo.JSONValueType) {
	pd.System.LogToConsole(h.path + " willDecodeSublist " + typeName(valueType) + " " + name)
	if name != "" {
		h.path = h.path + "." + name
	}
}

func (h *JSONHandler) ShouldDecodeTableValueForKey(key string) bool {
	pd.System.LogToConsole(h.path + " shouldDecodeTableValueForKey " + key)
	return true
}

func (h *JSONHandler) DidDecodeTableValue(key string, value pdgo.JSONValue) {
	pd.System.LogToConsole(h.path + " didDecodeTableValue " + key + " " + typeName(value.Type))
}

func (h *JSONHandler) ShouldDecodeArrayValueAtIndex(pos int) bool {
	pd.System.LogToConsole(h.path + " shouldDecodeArrayValueAtIndex")
	return true
}

func (h *JSONHandler) DidDecodeArrayValue(pos int, value pdgo.JSONValue) {
	pd.System.LogToConsole(h.path + " didDecodeArrayValue " + typeName(value.Type))
}

func (h *JSONHandler) DidDecodeSublist(name string, valueType pdgo.JSONValueType) {
	pd.System.LogToConsole(h.path + " didDecodeSublist " + typeName(valueType) + " " + name)
	// Pop path
	if name != "" && len(h.path) > len(name)+1 {
		h.path = h.path[:len(h.path)-len(name)-1]
	}
}

func typeName(t pdgo.JSONValueType) string {
	switch t {
	case pdgo.JSONNull:
		return "null"
	case pdgo.JSONTrue:
		return "true"
	case pdgo.JSONFalse:
		return "false"
	case pdgo.JSONInteger:
		return "integer"
	case pdgo.JSONFloat:
		return "float"
	case pdgo.JSONString:
		return "string"
	case pdgo.JSONArray:
		return "array"
	case pdgo.JSONTable:
		return "table"
	default:
		return "???"
	}
}

func initGame() {
	pd.System.LogToConsole("JSON Lowlevel - Playdate JSON API callbacks")
	pd.System.LogToConsole("=========================================")

	file, err := pd.File.Open("assets/test.json", pdgo.FileRead)
	if err != nil {
		pd.System.LogToConsole("Error opening file")
		return
	}

	handler := &JSONHandler{path: ""}
	err = pd.JSON.DecodeFile(file, handler)
	pd.File.Close(file)

	if err != nil {
		pd.System.LogToConsole("Error decoding JSON")
		return
	}

	pd.System.LogToConsole("")
	pd.System.LogToConsole("=========================================")
	pd.System.LogToConsole("JSON parsing complete!")
}

func update() int {
	pd.Graphics.Clear(pdgo.SolidWhite)
	pd.Graphics.DrawText("JSON Lowlevel", 10, 10)
	pd.Graphics.DrawText("Check console for output", 10, 30)
	return 1
}

func main() {}
