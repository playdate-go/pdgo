
<img src="assets/gopher-on-playdate.jpg" alt="Gopher on Playdate" width="300"  >

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/playdate-go/pdgo/main/install.sh | bash
```

This installs **everything** you need:
- `pdgoc` - the build tool
- `TinyGo` with Playdate support (for device builds)
- Configures your PATH automatically

**Build time:** ~25 minutes (mostly LLVM compilation on macOS)

**Options:**
```bash
# Skip TinyGo build (simulator-only, instant install)
SKIP_TINYGO=1 curl -fsSL https://raw.githubusercontent.com/playdate-go/pdgo/main/install.sh | bash

# Custom parallel jobs
JOBS=8 curl -fsSL ... | bash
```

---

## Menu

- [Overview](#overview)
- [Quick Install](#quick-install)
- [Usage](#usage)
- [Internals](#internals)
- [Why Not Go But TinyGo](#why-not-go-but-tinygo)
- [Flow](#flow)
- [How To Start](#how-to-start)
- [API Bindings](#api-bindings)
- [Examples](#examples)
- [Roadmap](#roadmap)
- [Community](#community)
- [Attribution](#attribution)
- [License](#license)
---

## Overview

>[!NOTE]  
> This project is currently under active development, all API covered but not all features have been fully tested or implemented yet, PRs are always welcome.  
> Tested on macOS, tinygo 0.40.1 darwin/arm64, go1.25.6, LLVM 20.1.1, Playdate OS 3.0.2.

>[!NOTE]  
> Playdate SDK >= 3.0.2 is required  
> Golang >= 1.21 is requred

Hi, my name is Roman Bielyi, and I'm developing this project in my spare time as a personal initiative.
This project is an independent effort and is neither endorsed by nor affiliated with [Panic Inc](https://panic.com/).

As a Go developer, I immediately wanted to bring Go to the [Playdate](https://play.date/). It wasn’t straightforward, but I got it working -- hope you’ll enjoy experimenting with it.

>[!IMPORTANT]  
The main objective now is to release a stable 1.0.x version.
To achieve this, we need to rewrite all official Playdate SDK examples from C/Lua into Go and ensure that the Go API bindings are mature, stable, and provide complete coverage of all subsystems.

## Usage

`pdgoc` is a command-line tool that handles **everything** for building Go apps for Playdate, both Simulator and Device builds.

> [!IMPORTANT]  
> **Always use `pdgoc` for building.** Do not try to run `go build` or `tinygo build` directly because `pdgoc` handles all the complexity: SDK paths, CGO flags, temporary files, etc.

> [!TIP]  
> The `sim` and `device` flags can be combined to build for both Simulator and Device simultaneously.

| Flag     | Description                                                   |
|----------|---------------------------------------------------------------|
| `sim`    | Builds project for the Playdate Simulator only                |
| `device` | Builds project for the Playdate console only                  | 
| `run`    | Builds and runs project in the Playdate Simulator             |
| `deploy` | Deploys and runs on connected Playdate device (requires `-device`) | 


| Flag                | Description                                     |
|---------------------|-------------------------------------------------|
| `name`              | Sets the `name` property for pdxinfo            |
| `author`            | Sets the `author` property for pdxinfo          |
| `desc`              | Sets the `description` property for pdxinfo     | 
| `bundle-id`         | Sets the `bundleID` property for pdxinfo        | 
| `version`           | Sets the `version` property for pdxinfo         |
| `build-number`      | Sets the `buildNumber` property for pdxinfo     |
| `image-path`        | Sets the `imagePath` property for pdxinfo       |
| `launch-sound-path` | Sets the `launchSoundPath` property for pdxinfo | 
| `content-warn`      | Sets the `contentWarning` property for pdxinfo  |
| `content-warn2`     | Sets the `contentWarning2` property for pdxinfo |


> [!NOTE]  
> To use the `pdgoc` CLI tool, navigate to the project root directory -- the one containing the `Source` folder with your `.go` source files, `go.mod`, `go.sum`, and any assets.  
> Simply execute `pdgoc` from there. It will detect the 'Source' directory automatically.
>
> **Example:**  
> If your structure looks like this:
> ```
> your-project/
> ├── Source/
> │   ├── main.go
> │   ├── go.mod
> │   ├── go.sum
> │   └── assets/ (images, sounds, etc.)
> └──
> ```  
>
> Then `cd your-project/` and run `pdgoc`.

Example:
```bash
pdgoc -device -sim \
  -name=MyApp \
  -author=YourName \
  -desc="My App" \
  -bundle-id=com.yourname.myapp \
  -version=1.0 \
  -build-number=1
```

The `main.go`:

```go
package main

import (
	"github.com/playdate-go/pdgo"
)

// A global pointer to the Playdate API. 
//Initialized automatically when the game starts. 
//All SDK calls go through this variable: pd.Graphics.DrawText(...), pd.System.DrawFPS(...), etc.
var pd *pdgo.PlaydateAPI


// Called once when the game launches (during kEventInit).
// Use this to load images, sounds, fonts, and initialize your game state. The Playdate API (pd) is fully available here.
func initGame() {
	
}

// The main game loop. Called every frame (~30 FPS by default). Here you:
// Handle input (pd.System.GetButtonState())
// Update game logic
// Draw graphics (pd.Graphics.DrawText(), pd.Graphics.DrawBitmap())
// Return value: 1 to tell Playdate the display was updated and needs refresh. Return 0 if nothing changed (saves battery).
func update() int {
	
}

// Must exist but remains empty. 
//Playdate doesn't use Go's normal main() entry point, instead, the SDK calls eventHandler which is generated by pdgoc
func main() {}

```

# Internals

### Device:

For Playdate hardware `pdgoc` uses a custom [TinyGo](https://tinygo.org/) build with full Playdate hardware support instead of standard Go build tools.

**Custom GC**:  
Custom TinyGo build for Playdate achieves minimal binary sizes through custom Playdate GC:
Instead of including traditional garbage collection logic (mark, sweep, write barriers), we delegate all allocations to the Playdate SDK's `realloc` function. This results in minimal GC wrapper code in the final binary. The Playdate OS already provides a robust allocator optimized for the hardware, and we simply use it!

<details>
<summary>click to see: gc_playdate.go</summary>

```go
//go:noinline
func alloc(size uintptr, layout unsafe.Pointer) unsafe.Pointer {
    size = align(size)
    gcTotalAlloc += uint64(size)
    gcMallocs++

    ptr := _cgo_pd_realloc(nil, size)
    if ptr == nil {
        runtimePanic("out of memory")
    }

    // Zero the allocated memory
    memzero(ptr, size)
    return ptr
}

func GC() {
    // No-op - Playdate SDK manages memory
}

func SetFinalizer(obj interface{}, finalizer interface{}) {
    // No-op
}
```

</details>

**No Static Heap**:  
Standard TinyGo embedded targets reserve heap space in BSS section. Our runtime configuration eliminates this by setting `needsStaticHeap = false`.
As a result, BSS is reduced from approx. 1MB to approx. 300 bytes.

<details>
<summary>click to see: gc_playdate.go</summary>

```go
const needsStaticHeap = false

func initHeap() {
    // No initialization needed - Playdate SDK handles it
}
```

</details>

**Minimal Runtime Configuration**:  
No scheduler, no threading, no dynamic stack management. A fixed stack size of ~60KB is used instead of Go's traditional growable stacks.

<details>
<summary>click to see: playdate.json</summary>

```json
{
    "inherits": ["cortex-m"],
    "llvm-target": "thumbv7em-unknown-unknown-eabihf",
    "cpu": "cortex-m7",
    "features": "+armv7e-m,+dsp,+hwdiv,+thumb-mode,+fp-armv8d16sp,+vfp4d16sp",
    "build-tags": ["playdate", "tinygo", "gc.playdate"],
    "gc": "playdate",
    "scheduler": "none",
    "serial": "none",
    "automatic-stack-size": false,
    "default-stack-size": 61800,
    "cflags": ["-DTARGET_PLAYDATE=1", "-mfloat-abi=hard", "-mfpu=fpv5-sp-d16"]
}
```

</details>

**LLVM Optimization:**
* Target: `thumbv7em-unknown-unknown-eabihf`
* CPU: `cortex-m7` with FPU (`-mfpu=fpv5-sp-d16`, `-mfloat-abi=hard`)
* Features: Thumb-2, DSP, hardware divide, VFP4
* Dead Code Elimination: Unused functions stripped via linker script
* Link-Time Optimization: Cross-module inlining via LLVM

<details>
<summary>click to see: playdate.json & playdate.ld</summary>

**Target configuration:**

```json
{
    "llvm-target": "thumbv7em-unknown-unknown-eabihf",
    "cpu": "cortex-m7",
    "features": "+armv7e-m,+dsp,+hwdiv,+thumb-mode,+fp-armv8d16sp,+vfp4d16sp",
    "cflags": ["-DTARGET_PLAYDATE=1", "-mfloat-abi=hard", "-mfpu=fpv5-sp-d16"]
}
```

**Linker script (dead code elimination):**

```ld
ENTRY(eventHandlerShim)

SECTIONS
{
    .text : ALIGN(4) {
        KEEP(*(.text.eventHandlerShim))
        KEEP(*(.text.eventHandler))
        KEEP(*(.text.updateCallback))
        KEEP(*(.text.runtime_init))
        *(.text) *(.text.*) *(.rodata) *(.rodata.*)
        KEEP(*(.init)) KEEP(*(.fini))
        . = ALIGN(4);
    }
    ...
    /DISCARD/ : { *(.ARM.exidx*) *(.ARM.extab*) }
}
```

`/DISCARD/` removes unused ARM exception sections, `KEEP()` preserves only critical entry points.

</details>

### Simulator:

`pdgoc` uses Go's native build tools to compile apps for the Playdate Simulator.

**Under the hood, it automatically runs:**
```bash
go build -ldflags="-w -s" -gcflags="all=-l" \
  -trimpath -buildvcs=false -race=false \
  -buildmode=c-shared \
  -o "some/output" "some/input"
```

All flags are optimized: stripping debug info (`-w -s`), disabling race detector, and producing a C-shared library with `-buildmode=c-shared` needed for Simulator instead of binary executable.
In Unix systems it's `.so` and in macOS  it's `.dylib`

| Flag                  | Purpose                                                                                  |
|-----------------------|------------------------------------------------------------------------------------------|
| `-ldflags="-w -s"`    | **`-w`**: Strip debug info (DWARF). **`-s`**: Strip symbol table. Shrinks binary ~30-50% |
| `-gcflags="all=-l"`   | Disable function inlining & optimizations for simulator compatibility                    |
| `-trimpath`           | Remove local filesystem paths from binaries (security/portability)                       |
| `-buildvcs=false`     | Skip embedding VCS data (git info) - faster builds                                       |
| `-race=false`         | Explicitly disable race detector (already off by default)                                |
| `-buildmode=c-shared` | **Key**: Build as C-shared library (`.dylib`/`.so`) for Playdate Simulator               |


## Why Not Go But TinyGo

No Bare-Metal ARM Support:  
Standard Go compiler (gc) only supports these targets:

| Flag    | Purpose                     |
|---------|-----------------------------|
| linux   | amd64, arm64, arm, 386, ... |
| darwin  | amd64, arm64                |
| windows | amd64, arm64, 386           |

Playdate requires: thumbv7em-none-eabihf (ARM Cortex-M7, no OS), and this is simply impossible:  
`GOOS=none GOARCH=thumbv7em go build  # not supported`

Size:  
Standard Go runtime includes Garbage Collector, Goroutine Scheduler, Stack Management and Reflection, binary size approx. 2-5 MB minimum.
Playdate constraints are 16 MB total RAM (shared with game data, graphics, sound), games typically 50 KB - 2 MB.

| Feature            | Standard Go    | TinyGo                      |
|--------------------|----------------|-----------------------------|
| Bare-metal support | No             | Yes                         |
| GOOS= not required | No             | Yes                         |
| ARM Cortex-M       | No             | Yes thumbv7em target        |
| Minimal runtime    | No approx. 2MB | Yes approx. 1-4 KB          |
| Custom GC          | No             | Yes pluggable (gc.playdate) |
| No OS required     | No             | Yes                         |
| Relocatable code   | No             | Yes via LLVM                |
| CGO on bare-metal  | No             | Yes (with custom runtime)   |

In short:

```
Go Source -> TinyGo Frontend -> LLVM IR -> LLVM Backend -> ARM Thumb-2 ELF
                                              |
                              Cortex-M7 optimizations
                              Position-independent code
                              Dead code elimination
```

Summary: Standard Go is designed for desktop/server environments, full operating systems, abundant memory (GB).
Playdate requires: bare-metal ARM Cortex-M7, no operating system, tiny runtime, SDK-managed memory.

TinyGo bridges this gap by reimplementing Go compilation targeting embedded systems with LLVM backend. Our custom TinyGo build supports CGO on bare-metal Playdate hardware through a unified C wrapper layer (`pd_cgo.c`).

## Flow:

### Device Build

```
┌─────────────────────────────────────────────────────────────┐
│  pdgoc -device                                              │
├─────────────────────────────────────────────────────────────┤
│  1. Copy pd_cgo.c from pdgo module to build/pd_runtime.c    │
│  2. Create Source/main_tinygo.go                            │
│  3. Run go mod tidy                                         │
│  4. Create /tmp/device-build-*.sh                           │
│  5. Execute build script:                                   │
│     ├── Compile pd_runtime.c -> pd_runtime.o -> libpd.a     │
│     ├── Create build/playdate.ld                            │
│     ├── Create ~/tinygo-playdate/targets/playdate.json      │
│     ├── TinyGo build -> pdex.elf                            │
│     ├── pdc -> GameName.pdx/                                │
│     └── Delete build/ directory                             │
│  6. Delete Source/main_tinygo.go                            │
│  7. Delete Source/pdxinfo                                   │
└─────────────────────────────────────────────────────────────┘
```
### Simulator Build

```
┌─────────────────────────────────────────────────────────────┐
│  pdgoc -sim                                                 │
├─────────────────────────────────────────────────────────────┤
│  1. Create Source/main_cgo.go                               │
│  2. go build -buildmode=c-shared -> pdex.dylib + pdex.h     │
│  3. Delete Source/pdex.h                                    │
│  4. Delete Source/main_cgo.go                               │
│  5. pdc -> GameName_sim.pdx/                                │
│  6. Delete Source/pdex.dylib                                │
│  7. Delete Source/pdxinfo                                   │
└─────────────────────────────────────────────────────────────┘
```

Temporary files created by `pdgoc` during build:

**Device Build Files**

`pd_runtime.c` - Copied from `pd_cgo.c` in the pdgo module. This C file provides all Playdate SDK wrappers that Go code calls via CGO. It includes TinyGo runtime support functions (`_cgo_pd_realloc`, `_cgo_pd_logToConsole`, `_cgo_pd_getCurrentTimeMS`) and the `eventHandler` entry point. Compiled with `-DTARGET_PLAYDATE=1` to enable device-specific code.

`main_tinygo.go` - Contains the `//export go_init` and `//export go_update` directives that tell TinyGo to expose these functions as C-callable symbols. The C runtime calls these functions to initialize the game and run the update loop. This file is separate from the user's main.go to avoid polluting their code with build-specific exports.

`playdate.ld` - The linker script tells the ARM linker how to arrange code and data in memory. It defines the entry point (eventHandlerShim), ensures critical functions appear at the beginning of the binary, and sets up BSS/data sections.

`playdate.json` - TinyGo's target configuration file. It specifies the CPU architecture (Cortex-M7), compiler flags, which garbage collector to use (gc.playdate), and links to the linker script. This file tells TinyGo exactly how to compile for Playdate hardware.

`libpd.a` - A static library compiled from pd_runtime.c. TinyGo links against this library to resolve the C function references. Static linking ensures all SDK wrapper code is embedded directly in the final binary.


| File                 | Location                     | Purpose                          | Cleanup                   |
|----------------------|------------------------------|----------------------------------|---------------------------|
| `pd_runtime.c`       | `build/`                     | C wrappers (copy of pd_cgo.c)    | Deleted with `build/` dir |
| `main_tinygo.go`     | `Source/`                    | TinyGo entry points (`//export`) | Deleted after build       |
| `device-build-*.sh`  | `/tmp/`                      | Embedded build script            | Deleted after build       |
| `playdate.ld`        | `build/`                     | Linker script                    | Deleted with `build/` dir |
| `playdate.json`      | `~/tinygo-playdate/targets/` | TinyGo target config             | Overwritten each build    |
| `pdex.elf`           | `build/`                     | Compiled ELF binary              | Deleted with `build/` dir |
| `pd_runtime.o`       | `build/`                     | Compiled C object                | Deleted with `build/` dir |
| `libpd.a`            | `build/`                     | Static C library                 | Deleted with `build/` dir |
| `pdxinfo`            | `Source/`                    | Game metadata                    | Deleted after build       |


**Simulator Build Files**

The simulator build uses the same `pd_cgo.c` from the pdgo module as the device build, but compiled for the host architecture. The C wrappers are linked directly into the shared library via standard Go CGO.

`main_cgo.go` - Contains `import "C"` and `//export eventHandler` directive that tells the standard Go compiler to generate a C-callable entry point. The simulator runs on your host machine (macOS/Linux), where CGO is fully supported.

`pdex.h` - Automatically generated by go build -buildmode=c-shared. This header file contains C function declarations for all exported Go functions. We immediately delete it since Playdate doesn't need it - the SDK already knows the expected function signatures.

`pdex.dylib` / `pdex.so` - The compiled shared library containing your Go game code and the C wrappers. The Playdate Simulator dynamically loads this library at runtime and calls eventHandler when your game starts. This file is moved into the .pdx bundle by pdc.

| File                     | Location  | Purpose                       | Cleanup             |
|--------------------------|-----------|-------------------------------|---------------------|
| `main_cgo.go`            | `Source/` | CGO entry points (`//export`) | Deleted after build |
| `pdex.h`                 | `Source/` | CGO header (auto-generated)   | Deleted after build |
| `pdex.dylib` / `pdex.so` | `Source/` | Compiled shared library       | Deleted after build |
| `pdxinfo`                | `Source/` | Game metadata                 | Deleted after build |


## How To Start

You must build a custom TinyGo compiler with Playdate hardware support before using the pdgoc CLI tool.
Thankfully it's easy to do using just one `scripts/build-tinygo-playdate.sh`.

### Step 1: Install LLVM (Linux only)

> [!IMPORTANT]
> **macOS users**: Skip this step! Homebrew LLVM does not include static LLD libraries required to build TinyGo.
> The build script will automatically compile LLVM from source (~25-30 minutes, only needed once).

Pre-installing LLVM on Linux reduces build time from **~25-30 minutes** to **~1 minute**.

#### Linux (Ubuntu/Debian)

```bash
# Add LLVM apt repository
wget https://apt.llvm.org/llvm.sh
chmod +x llvm.sh
sudo ./llvm.sh 20

# Install ALL required dev packages
sudo apt install clang-20 llvm-20-dev lld-20 libclang-20-dev
```

#### Linux (Fedora)

```bash
sudo dnf install llvm20-devel clang20-devel lld20-devel
```

#### Linux (Arch)

```bash
sudo pacman -S llvm clang lld
```

### Step 2: Build TinyGo

> [!NOTE]
> The build script checks for missing dependencies (go, cmake, ninja, git) and tells you how to install them.
>
> **ARM Toolchain:**
> - **macOS**: Included with Playdate SDK - no separate installation needed
> - **Linux**: Install `gcc-arm-none-eabi` via your package manager

```bash
# Run the build script (auto-detects system installed LLVM)
./cmd/pdgoc/scripts/build-tinygo-playdate.sh
```

**Build script options:**

| Flag                   | Description                                         |
|------------------------|-----------------------------------------------------|
| `--tinygo-version VER` | TinyGo version to build (default: 0.40.1)           |
| `--use-system-llvm`    | Force using system LLVM (fails if not found)        |
| `--build-llvm`         | Force building LLVM from source                     |
| `--jobs N`             | Number of parallel jobs                             |

**What `build-tinygo-playdate.sh` does:**

The script builds a custom TinyGo compiler with Playdate support:

1. **Clone TinyGo** - downloads TinyGo v0.40.1 (or specified version) to `~/tinygo-playdate`
2. **Setup LLVM** - uses system installed LLVM if available, otherwise builds from source
3. **Add Playdate Files** - injects custom files into TinyGo:
   - `playdate.json` - target config (Cortex-M7, custom GC, no scheduler)
   - `playdate.ld` - linker script (memory layout, entry point)
   - `runtime_playdate.go` - platform runtime (time, console output via SDK)
   - `gc_playdate.go` - custom GC that delegates to Playdate SDK's realloc
4. **Build TinyGo** - compiles the TinyGo binary with Playdate support

Result: `~/tinygo-playdate/build/tinygo` - a TinyGo compiler that accepts `-target=playdate` and produces minimal ARM binaries (~4-5 KB) for Playdate hardware.

If all went ok - you will see this (versions may differ in your environment):

```
╔══════════════════════════════════════════════════════════╗
║                    Build Complete!                       ║
╚══════════════════════════════════════════════════════════╝

TinyGo binary: /Users/SomeUser/tinygo-playdate/build/tinygo

To use:
  export PATH="/Users/SomeUser/tinygo-playdate/build:$PATH"
  tinygo build -target=playdate -o game.elf .

tinygo version 0.40.1 darwin/arm64 (using go version go1.25.6 and LLVM version 20.1.1)
Playdate target added
```

>[!NOTE]
> With pre-installed LLVM: ~1 minute. Without: ~25 minutes.

Once you have TinyGo, you need install `pdgoc`

```bash
go install https://github.com/playdate-go/cmd/pdgoc

# check installation 
pdgoc version
```

That's it! and you may forget about TinyGo (`pdgoc` will use it automatically under the hood)

>[!WARNING]
> We don't support Windows OS currently, supported: Linux, macOS.  
> You may need to wait for Windows support, or you can make a contribution.

## API Bindings
The latest full documentation for API bindings is hosted here:
https://pkg.go.dev/github.com/playdate-go/pdgo#section-documentation

## Examples

> [!NOTE]
> We will add more complex examples as the project progresses

To build all examples please do this:
```bash
# in project repo root
chmod +x examples/build_all.sh 
chmod +x examples/*/build.sh
./examples/build_all.sh
```

Each example includes a `build.sh` script that runs `pdgoc` with all necessary flags.

**Bach MIDI** -- [examples/bach_midi](examples/bach_midi)

**3D Library** -- [examples/3d_library](examples/3d_library)

**Sprite Game** -- [examples/spritegame](examples/spritegame)

**Conway's Game of Life** -- [examples/life](examples/life)

**Bouncing Square** -- [examples/bouncing_square](examples/bouncing_square)

**Go Logo** -- [examples/go_logo](examples/go_logo)

**Hello World** -- [examples/hello_world](examples/hello_world)


## Roadmap

- [ ] Add more own complex code examples to cover and test all API subsystems
- [ ] Rewrite to Go all official examples from SDK
    - [x] Hello World
    - [x] Life
    - [ ] Tilemap
    - [x] Sprite Game
    - [ ] Sprite Collisions
    - [ ] Particles
    - [ ] Networking
    - [ ] JSON
    - [ ] Exposure
    - [x] Bach.mid
    - [ ] Array
    - [x] 3D Library
    - [ ] 2020
    - [ ] Accelerometer Test
    - [ ] Asteroids
    - [ ] ControllerTest
    - [ ] Drum Machine
    - [ ] Flippy Fish
    - [ ] Game Template
    - [ ] Hammer Down
    - [ ] Level 1-1
    - [ ] MIDI Player
    - [ ] Node7Driver
    - [ ] Networking
    - [ ] Pathfinder
    - [ ] Single File Examples
    - [ ] Sprite Collisions Masks
- [ ] Make sure Lua interoperability works
- [ ] Make sure C interoperability works
- [X] Write documentation for API bindings
- [ ] Create unit tests for `pdgoc` and API bindings
- [ ] Add support for Windows OS

## Community
Using these links and places, you can discuss the PdGo project with each other:

Slack
1) https://gophers.slack.com/archives/C029RQSEE/p1769119174451979
2) https://gophers.slack.com/archives/CDJD3SUP6/p1769119574841489

Reddit:
1) https://www.reddit.com/r/golang/comments/1qk1ec9/golang_support_for_playdate_handheld_compiler_sdk/
2) https://www.reddit.com/r/PlaydateDeveloper/comments/1qk0r60/golang_support_for_playdate_handheld_compiler_sdk/
3) https://www.reddit.com/r/programming/comments/1qk19kb/playdate_supports_go_language_compiler_sdk/
4) https://www.reddit.com/r/PlaydateConsole/comments/1qk0wy0/golang_support_for_playdate_handheld_compiler_sdk/

**Discord**:
1) https://discord.com/channels/118456055842734083/1464001888243548181
2) https://discord.com/channels/675983554655551509/1464004567476867247

**Playdate Development Forum (this is the main place to discuss)**:
https://devforum.play.date/t/golang-support-for-playdate-compiler-sdk-bindings-tools-and-examples/24919

## Attribution

The Go Gopher was designed by [Renee French](https://reneefrench.blogspot.com/)
and is licensed under [Creative Commons 4.0 Attribution License](https://creativecommons.org/licenses/by/4.0/).

## License

MIT License

Copyright (c) 2026 Roman Bielyi and PdGo contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
