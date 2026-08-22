// gc_pause_benchmark - per-frame GC pause benchmark
// Output: CSV to console: frame,NumGC,HeapAlloc,LastPauseNs,LiveObjects

package main

import (
	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

const numParticles = 1000

type particle struct {
	x, y   float32
	vx, vy float32
	life   int
}

var particles []particle
var frame int
var done bool

func initGame() {
	particles = make([]particle, 0, numParticles)
	pd.System.LogToConsole("frame,NumGC,HeapAlloc,LastPauseNs,LiveObjects")
}

func update() int {
	pd.Graphics.Clear(pdgo.SolidWhite)

	// Spawn particles.
	for i := 0; i < 20; i++ {
		particles = append(particles, particle{
			x:    200,
			y:    120,
			vx:   float32((frame + i) % 10),
			vy:   float32(i),
			life: 60,
		})
	}

	// Update + cull dead particles (reuses backing array).
	live := particles[:0]
	for _, p := range particles {
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life > 0 {
			live = append(live, p)
		}
	}
	particles = live

	// Simulate per-frame game garbage (temp allocations that become
	// unreachable each frame — physics scratch, collision results, etc.).
	for i := 0; i < 50; i++ {
		scratch := make([]byte, 64)
		scratch[0] = byte(i)
		_ = scratch[0]
	}

	// Per-frame CSV sample.
	stats := pd.Memory.Stats()
	logCSVRow(frame, stats.NumGC, stats.HeapAlloc, stats.LastPauseNs, stats.LiveObjects)

	frame++
	if frame >= 60 && !done {
		done = true
		pd.System.LogToConsole("=== benchmark complete ===")
		// The update callback is called every frame for the lifetime of
		// the game — returning 0 only skips the display refresh, it does
		// not stop the loop. So: log once, then idle quietly.
	}
	return 1
}

// logCSVRow emits "frame,NumGC,HeapAlloc,LastPauseNs,LiveObjects" without fmt.
func logCSVRow(frame int, numGC uint32, heap uint64, pause uint64, live uint32) {
	var buf [128]byte
	n := 0
	n += writeInt(buf[n:], frame)
	buf[n] = ','
	n++
	n += writeUint64(buf[n:], uint64(numGC))
	buf[n] = ','
	n++
	n += writeUint64(buf[n:], heap)
	buf[n] = ','
	n++
	n += writeUint64(buf[n:], pause)
	buf[n] = ','
	n++
	n += writeUint64(buf[n:], uint64(live))
	pd.System.LogToConsole(string(buf[:n]))
}

func writeInt(buf []byte, v int) int {
	if v == 0 {
		buf[0] = '0'
		return 1
	}
	tmp := [16]byte{}
	n := 0
	for v > 0 {
		tmp[n] = byte('0' + v%10)
		v /= 10
		n++
	}
	for i := 0; i < n; i++ {
		buf[i] = tmp[n-1-i]
	}
	return n
}

func writeUint64(buf []byte, v uint64) int {
	if v == 0 {
		buf[0] = '0'
		return 1
	}
	tmp := [24]byte{}
	n := 0
	for v > 0 {
		tmp[n] = byte('0' + v%10)
		v /= 10
		n++
	}
	for i := 0; i < n; i++ {
		buf[i] = tmp[n-1-i]
	}
	return n
}

func main() {}
