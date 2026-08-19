// pdgo replacement for TinyGo's stock src/runtime/asm_arm.S (device builds).
//
// Fixes two bugs in the stock tinygo_scanCurrentStack that break the gc.playdate
// collector on-device:
//
//  1. The epilogue ("add sp, #32; pop {pc}") pops the saved registers without
//     restoring r4-r11. Every GC stack scan therefore corrupts the caller's
//     callee-saved registers — LLVM keeps live values there across the
//     runtime.GC() call site, so this silently poisons arbitrary Go code.
//
//  2. The 9-word push leaves SP 4 mod 8 at the "bl tinygo_scanstack", an AAPCS
//     violation under the hard-float ABI (this target builds with
//     -mfpu=fpv5-sp-d16). "sub sp, #4" pads to 8-byte alignment.
//
// The whole routine is restricted to 16-bit Thumb encodings so GNU as
// (arm-none-eabi-gcc) can assemble it — the stock __thumb2__ branch's
// 9-register push needs the 32-bit encoding GNU as rejects. LLVM's integrated
// assembler handles this file identically.

// Only generate .debug_frame, don't generate .eh_frame.
.cfi_sections .debug_frame

.section .text.tinygo_scanCurrentStack
.global  tinygo_scanCurrentStack
.type    tinygo_scanCurrentStack, %function
tinygo_scanCurrentStack:
    .cfi_startproc
    // Save callee-saved registers onto the stack. r4-r7 and lr first, then a
    // 4-byte pad (9 saved words would leave SP 4 mod 8), then the high
    // registers staged through r0-r3. Total 10 words = 40 bytes: 8-byte
    // aligned at the call, and every saved register lies inside the
    // [sp, pdStackTop) range tinygo_scanstack scans.
    push {r4-r7, lr}
    .cfi_def_cfa_offset 5*4
    sub  sp, #4
    .cfi_def_cfa_offset 6*4
    mov  r0, r8
    mov  r1, r9
    mov  r2, r10
    mov  r3, r11
    push {r0-r3}
    .cfi_def_cfa_offset 10*4

    // Scan the stack.
    mov r0, sp
    bl tinygo_scanstack

    // Restore stack state and return. The staged r0-r3 are reloaded from
    // memory, so tinygo_scanstack clobbering them is fine.
    pop {r0-r3}
    .cfi_def_cfa_offset 6*4
    mov  r8, r0
    mov  r9, r1
    mov  r10, r2
    mov  r11, r3
    add  sp, #4
    .cfi_def_cfa_offset 5*4
    pop {r4-r7, pc}
    .cfi_endproc
.size tinygo_scanCurrentStack, .-tinygo_scanCurrentStack


.section .text.tinygo_longjmp
.global tinygo_longjmp
.type   tinygo_longjmp, %function
tinygo_longjmp:
    .cfi_startproc
    // Note: the code we jump to assumes r0 is set to a non-zero value if we
    // jump from here (which is conveniently already the case).

    ldm r0, {r0, r1}
    mov sp, r0 // jumpSP
    mov pc, r1 // jumpPC
    .cfi_endproc
.size tinygo_longjmp, .-tinygo_longjmp
