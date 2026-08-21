// pdgo retention registries - keeps Go wrappers alive while the C SDK stores
// their pointers, so finalizers cannot free objects the SDK still uses.
//
// Every API that hands a wrapper's pointer into SDK state retains the wrapper
// here (package-level maps are GC roots) and releases it at the matching
// removal point: the corresponding Remove*/Free* call, the next call that
// replaces the slot, or the owner's finalizer. Game code therefore never
// needs to keep references to SDK-owned objects.

package pdgo

import "unsafe"

var (
	// displaySprites: sprites currently on the SDK display list,
	// keyed by sprite C pointer.
	displaySprites = make(map[unsafe.Pointer]*LCDSprite)

	// Slot state: one retained wrapper per owner, keyed by the owner's
	// C pointer. Slots are replaced by the next Set* call and released
	// when the owner is freed.
	spriteImages        = make(map[unsafe.Pointer]*LCDBitmap)      // SetImage
	tilemapTables       = make(map[unsafe.Pointer]*LCDBitmapTable) // SetImageTable
	synthSamples        = make(map[unsafe.Pointer]*AudioSample)    // SetSample
	samplePlayerSamples = make(map[unsafe.Pointer]*AudioSample)    // SetSamplePlayerSample

	// Single global slots, replaced by the next call that sets them.
	currentFont    *LCDFont   // SetFont
	currentStencil *LCDBitmap // SetStencilImage
	patternBitmap  *LCDBitmap // SetColorToPattern
	menuImage      *LCDBitmap // SetMenuImage (never unset; kept for program lifetime)

	// contextStack: bitmaps pushed via PushContext. Contexts nest in the
	// SDK, and a nil target (display framebuffer) is pushed as-is so the
	// stack stays aligned with PopContext calls.
	contextStack []*LCDBitmap

	// Sound objects the SDK holds with no removal API available:
	// retained until an explicit FreeInstrument/FreeSynth.
	retainedInstruments = make(map[unsafe.Pointer]*PDSynthInstrument)          // AddInstrumentAsSource, SetInstrument, AddVoice
	instrumentVoices    = make(map[unsafe.Pointer]map[unsafe.Pointer]*PDSynth) // instrument -> synth -> synth
)

// ---- display list ----

func retainDisplaySprite(sp *LCDSprite) {
	if sp != nil && sp.ptr != nil {
		displaySprites[sp.ptr] = sp
	}
}

func releaseDisplaySprite(ptr unsafe.Pointer) {
	delete(displaySprites, ptr)
}

func clearDisplaySprites() {
	for ptr := range displaySprites {
		delete(displaySprites, ptr)
	}
}

// ---- per-owner slots ----

func retainSpriteImage(spritePtr unsafe.Pointer, b *LCDBitmap) {
	if b == nil {
		delete(spriteImages, spritePtr)
		return
	}
	spriteImages[spritePtr] = b
}

func releaseSpriteImage(spritePtr unsafe.Pointer) {
	delete(spriteImages, spritePtr)
}

func retainTilemapTable(tilemapPtr unsafe.Pointer, t *LCDBitmapTable) {
	if t == nil {
		delete(tilemapTables, tilemapPtr)
		return
	}
	tilemapTables[tilemapPtr] = t
}

func releaseTilemapTable(tilemapPtr unsafe.Pointer) {
	delete(tilemapTables, tilemapPtr)
}

func retainSynthSample(synthPtr unsafe.Pointer, s *AudioSample) {
	if s == nil {
		delete(synthSamples, synthPtr)
		return
	}
	synthSamples[synthPtr] = s
}

func releaseSynthSample(synthPtr unsafe.Pointer) {
	delete(synthSamples, synthPtr)
}

func retainPlayerSample(playerPtr unsafe.Pointer, s *AudioSample) {
	if s == nil {
		delete(samplePlayerSamples, playerPtr)
		return
	}
	samplePlayerSamples[playerPtr] = s
}

func releasePlayerSample(playerPtr unsafe.Pointer) {
	delete(samplePlayerSamples, playerPtr)
}

// ---- sound instruments and voices ----

func retainInstrument(inst *PDSynthInstrument) {
	if inst != nil && inst.ptr != nil {
		retainedInstruments[inst.ptr] = inst
	}
}

// releaseInstrument drops the instrument and its retained voices.
func releaseInstrument(instPtr unsafe.Pointer) {
	delete(retainedInstruments, instPtr)
	delete(instrumentVoices, instPtr)
}

func retainInstrumentVoice(instPtr unsafe.Pointer, synth *PDSynth) {
	if synth == nil || synth.ptr == nil {
		return
	}
	voices, ok := instrumentVoices[instPtr]
	if !ok {
		voices = make(map[unsafe.Pointer]*PDSynth)
		instrumentVoices[instPtr] = voices
	}
	voices[synth.ptr] = synth
}

// releaseSynthVoices removes the synth from the voice set of every
// instrument (used when a synth is explicitly freed).
func releaseSynthVoices(synthPtr unsafe.Pointer) {
	for instPtr, voices := range instrumentVoices {
		delete(voices, synthPtr)
		if len(voices) == 0 {
			delete(instrumentVoices, instPtr)
		}
	}
}

// ---- context stack ----

func pushRetainedContext(b *LCDBitmap) {
	contextStack = append(contextStack, b)
}

func popRetainedContext() *LCDBitmap {
	n := len(contextStack)
	if n == 0 {
		return nil
	}
	b := contextStack[n-1]
	contextStack = contextStack[:n-1]
	return b
}
