package main

import (
	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

func initGame() {
	pd.System.LogToConsole("JSON Example - High-level API")
	pd.System.LogToConsole("=========================================")

	file, err := pd.File.Open("assets/test.json", pdgo.FileRead)
	if err != nil {
		pd.System.LogToConsole("Error opening file")
		return
	}

	// Parse JSON into tree
	root, err := pd.JSON.ParseFile(file)
	pd.File.Close(file)

	if err != nil {
		pd.System.LogToConsole("Error parsing JSON")
		return
	}

	pd.System.LogToConsole("")
	pd.System.LogToConsole("=== Parsed JSON Structure ===")

	if root == nil {
		pd.System.LogToConsole("Error: root is nil")
		return
	}

	// Debug: check root
	pd.System.LogToConsole("root ok, checking web-app...")

	webApp := root.Get("web-app")
	if webApp == nil {
		pd.System.LogToConsole("Error: web-app not found")
		// Show keys
		if root.Object != nil {
			for k := range root.Object {
				pd.System.LogToConsole("key: " + k)
			}
		}
		return
	}

	pd.System.LogToConsole("web-app found!")

	servlets := webApp.Get("servlet")
	if servlets != nil {
		pd.System.LogToConsole("servlets found")

		for i := 0; i < servlets.Len(); i++ {
			servlet := servlets.At(i)
			if servlet != nil {
				name := servlet.Get("servlet-name")
				if name != nil {
					pd.System.LogToConsole("  - " + name.GetString())
				}
			}
		}
	}

	mappings := webApp.Get("servlet-mapping")
	if mappings != nil {
		pd.System.LogToConsole("servlet-mapping found")
	}

	taglib := webApp.Get("taglib")
	if taglib != nil {
		pd.System.LogToConsole("taglib found")
		uri := taglib.Get("taglib-uri")
		if uri != nil {
			pd.System.LogToConsole("  uri: " + uri.GetString())
		}
	}

	pd.System.LogToConsole("")
	pd.System.LogToConsole("=========================================")
	pd.System.LogToConsole("JSON parsing complete!")

	// ---- JSON Encoder Example ----
	pd.System.LogToConsole("")
	pd.System.LogToConsole("=== JSON Encoder ===")

	enc := pd.JSON.NewEncoder(true) // pretty print

	enc.StartObject()

	enc.WriteKey("name")
	enc.WriteString("Test Game")

	enc.WriteKey("version")
	enc.WriteInt(1)

	enc.WriteKey("enabled")
	enc.WriteBool(true)

	enc.WriteKey("scores")
	enc.StartArray()
	enc.AddArrayElement()
	enc.WriteInt(100)
	enc.AddArrayElement()
	enc.WriteInt(200)
	enc.AddArrayElement()
	enc.WriteInt(300)
	enc.EndArray()

	enc.WriteKey("settings")
	enc.StartObject()
	enc.WriteKey("sound")
	enc.WriteBool(true)
	enc.WriteKey("difficulty")
	enc.WriteString("normal")
	enc.EndObject()

	enc.EndObject()

	result := enc.String()
	pd.System.LogToConsole("Encoded JSON:")
	pd.System.LogToConsole(result)

	enc.Free()

	pd.System.LogToConsole("")
	pd.System.LogToConsole("=========================================")
}

func update() int {
	pd.Graphics.Clear(pdgo.SolidWhite)
	pd.Graphics.DrawText("JSON Example", 10, 10)
	pd.Graphics.DrawText("Check console for output", 10, 30)
	return 1
}

func main() {}
