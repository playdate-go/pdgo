package main

import (
	"github.com/playdate-go/pdgo"
)

const (
	textWidth  = 86
	textHeight = 16
)

var (
	pd   *pdgo.PlaydateAPI
	font *pdgo.LCDFont

	x = (pdgo.LCDColumns - textWidth) / 2
	y = (pdgo.LCDRows - textHeight) / 2

	frameCount = 0

	// Pre-allocated buffer for formatting (avoid allocations in hot path)
	textBuf [64]byte
)

// formatInt writes integer to buffer, returns length
func formatInt(buf []byte, n int) int {
	if n == 0 {
		buf[0] = '0'
		return 1
	}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return len(buf) - i
}

// formatStats formats stats into textBuf without allocating, returns length
func formatStats(calls, frees, bytes int) int {
	pos := 0

	// "calls="
	copy(textBuf[pos:], "calls=")
	pos += 6
	pos += formatInt(textBuf[pos:], calls)

	// " frees="
	copy(textBuf[pos:], " frees=")
	pos += 7
	pos += formatInt(textBuf[pos:], frees)

	// " bytes="
	copy(textBuf[pos:], " bytes=")
	pos += 7
	pos += formatInt(textBuf[pos:], bytes)

	// null terminate
	textBuf[pos] = 0
	return pos
}

func initGame() {
	// Load font
	var err error
	font, err = pd.Graphics.LoadFont("/System/Fonts/Asheville-Sans-14-Bold.pft")
	if err != nil {
		pd.System.Error("Couldn't load font")
	}

	// Enable realloc debug logging
	pd.Memory.SetReallocDebug(true)
	pd.System.LogToConsole("[DEBUG] realloc logging enabled")
}

func update() int {
	pd.Graphics.Clear(pdgo.SolidWhite)

	if font != nil {
		pd.Graphics.SetFont(font)
	}

	// Print stats every 60 frames (once per second at 60fps)
	frameCount++
	if frameCount%60 == 0 {
		stats, _ := pd.Memory.GetReallocStats()
		// Use pre-allocated formatting to avoid memory allocation
		length := formatStats(stats.Count, stats.FreeCount, int(stats.TotalBytes))
		pd.Graphics.DrawTextBytes(textBuf[:], length, x, y)
	}

	pd.System.DrawFPS(0, 0)
	return 1
}

func main() {}
