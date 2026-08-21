package pdgo

import (
	"testing"
	"unsafe"
)

// testPtr returns distinct, vet-clean fake pointers for map keys.
var testPtrSlots [16]byte

func testPtr(n int) unsafe.Pointer { return unsafe.Pointer(&testPtrSlots[n%len(testPtrSlots)]) }

func resetRetainState() {
	displaySprites = make(map[unsafe.Pointer]*LCDSprite)
	spriteImages = make(map[unsafe.Pointer]*LCDBitmap)
	tilemapTables = make(map[unsafe.Pointer]*LCDBitmapTable)
	synthSamples = make(map[unsafe.Pointer]*AudioSample)
	samplePlayerSamples = make(map[unsafe.Pointer]*AudioSample)
	retainedInstruments = make(map[unsafe.Pointer]*PDSynthInstrument)
	instrumentVoices = make(map[unsafe.Pointer]map[unsafe.Pointer]*PDSynth)
	contextStack = nil
	currentFont = nil
	currentStencil = nil
	patternBitmap = nil
	menuImage = nil
}

func TestDisplaySpriteRetainRelease(t *testing.T) {
	resetRetainState()
	sp := &LCDSprite{ptr: testPtr(0x100)}
	retainDisplaySprite(sp)
	if displaySprites[testPtr(0x100)] != sp {
		t.Fatal("sprite not retained")
	}
	// Idempotent re-retain.
	retainDisplaySprite(sp)
	if len(displaySprites) != 1 {
		t.Fatalf("want 1 entry, got %d", len(displaySprites))
	}
	releaseDisplaySprite(testPtr(0x100))
	if len(displaySprites) != 0 {
		t.Fatal("sprite not released")
	}
	// Release of unknown ptr is a no-op.
	releaseDisplaySprite(testPtr(0x999))
}

func TestClearDisplaySprites(t *testing.T) {
	resetRetainState()
	retainDisplaySprite(&LCDSprite{ptr: testPtr(1)})
	retainDisplaySprite(&LCDSprite{ptr: testPtr(2)})
	clearDisplaySprites()
	if len(displaySprites) != 0 {
		t.Fatalf("want empty display list, got %d", len(displaySprites))
	}
}

func TestSpriteImageSlot(t *testing.T) {
	resetRetainState()
	b1 := &LCDBitmap{ptr: testPtr(0xA1)}
	b2 := &LCDBitmap{ptr: testPtr(0xA2)}
	retainSpriteImage(testPtr(0x100), b1)
	if spriteImages[testPtr(0x100)] != b1 {
		t.Fatal("image not retained")
	}
	// Replace.
	retainSpriteImage(testPtr(0x100), b2)
	if spriteImages[testPtr(0x100)] != b2 {
		t.Fatal("image not replaced")
	}
	// nil clears the slot.
	retainSpriteImage(testPtr(0x100), nil)
	if _, ok := spriteImages[testPtr(0x100)]; ok {
		t.Fatal("nil image did not clear slot")
	}
	// Owner release.
	retainSpriteImage(testPtr(0x100), b1)
	releaseSpriteImage(testPtr(0x100))
	if len(spriteImages) != 0 {
		t.Fatal("image slot not released with owner")
	}
}

func TestTilemapTableSlot(t *testing.T) {
	resetRetainState()
	tbl := &LCDBitmapTable{ptr: testPtr(0xB1)}
	retainTilemapTable(testPtr(0x200), tbl)
	if tilemapTables[testPtr(0x200)] != tbl {
		t.Fatal("table not retained")
	}
	releaseTilemapTable(testPtr(0x200))
	if len(tilemapTables) != 0 {
		t.Fatal("table slot not released")
	}
}

func TestSampleSlots(t *testing.T) {
	resetRetainState()
	s1 := &AudioSample{ptr: testPtr(0xC1)}
	retainSynthSample(testPtr(0x300), s1)
	if synthSamples[testPtr(0x300)] != s1 {
		t.Fatal("synth sample not retained")
	}
	retainSynthSample(testPtr(0x300), nil) // SetSample(synth, nil)
	if _, ok := synthSamples[testPtr(0x300)]; ok {
		t.Fatal("nil sample did not clear synth slot")
	}
	retainSynthSample(testPtr(0x300), s1)
	releaseSynthSample(testPtr(0x300))
	if len(synthSamples) != 0 {
		t.Fatal("synth slot not released")
	}

	p1 := &AudioSample{ptr: testPtr(0xC2)}
	retainPlayerSample(testPtr(0x400), p1)
	if samplePlayerSamples[testPtr(0x400)] != p1 {
		t.Fatal("player sample not retained")
	}
	retainPlayerSample(testPtr(0x400), nil)
	if _, ok := samplePlayerSamples[testPtr(0x400)]; ok {
		t.Fatal("nil sample did not clear player slot")
	}
	retainPlayerSample(testPtr(0x400), p1)
	releasePlayerSample(testPtr(0x400))
	if len(samplePlayerSamples) != 0 {
		t.Fatal("player slot not released")
	}
}

func TestInstrumentRetention(t *testing.T) {
	resetRetainState()
	inst := &PDSynthInstrument{ptr: testPtr(0x500)}
	synth := &PDSynth{ptr: testPtr(0x501)}
	retainInstrument(inst)
	if retainedInstruments[testPtr(0x500)] != inst {
		t.Fatal("instrument not retained")
	}
	retainInstrumentVoice(testPtr(0x500), synth)
	if instrumentVoices[testPtr(0x500)][testPtr(0x501)] != synth {
		t.Fatal("voice not retained")
	}
	// Re-retain idempotent.
	retainInstrumentVoice(testPtr(0x500), synth)
	if len(instrumentVoices[testPtr(0x500)]) != 1 {
		t.Fatalf("want 1 voice, got %d", len(instrumentVoices[testPtr(0x500)]))
	}
	// releaseSynthVoices removes the synth from every instrument.
	releaseSynthVoices(testPtr(0x501))
	if len(instrumentVoices[testPtr(0x500)]) != 0 {
		t.Fatal("voice not released")
	}
	// releaseInstrument drops instrument + its voice map.
	retainInstrumentVoice(testPtr(0x500), synth)
	releaseInstrument(testPtr(0x500))
	if _, ok := retainedInstruments[testPtr(0x500)]; ok {
		t.Fatal("instrument not released")
	}
	if _, ok := instrumentVoices[testPtr(0x500)]; ok {
		t.Fatal("instrument voice map not released")
	}
}

func TestContextStack(t *testing.T) {
	resetRetainState()
	b := &LCDBitmap{ptr: testPtr(0xD1)}
	pushRetainedContext(b)
	pushRetainedContext(nil) // display target
	if len(contextStack) != 2 || contextStack[0] != b || contextStack[1] != nil {
		t.Fatalf("stack = %v", contextStack)
	}
	if popRetainedContext() != nil {
		t.Fatal("expected nil (display target) on top")
	}
	if popRetainedContext() != b {
		t.Fatal("expected bitmap on second pop")
	}
	if len(contextStack) != 0 {
		t.Fatal("stack not empty")
	}
	// Pop on empty must not panic.
	popRetainedContext()
}
