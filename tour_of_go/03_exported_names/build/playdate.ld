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
    .data : ALIGN(4) {
        __data_start__ = .; *(.data) *(.data.*) . = ALIGN(4); __data_end__ = .;
    }
    .bss (NOLOAD) : ALIGN(4) {
        __bss_start__ = .; _sbss = .; *(.bss) *(.bss.*) *(COMMON) . = ALIGN(4); __bss_end__ = .; _ebss = .;
    }
    /DISCARD/ : { *(.ARM.exidx*) *(.ARM.extab*) }
    _sidata = LOADADDR(.data);
    _sdata = __data_start__; _edata = __data_end__;
    _globals_start = __data_start__; _globals_end = __bss_end__;
    _stack_top = __bss_end__;
}
