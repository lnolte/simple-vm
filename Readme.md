# A Simple Virtual Machine

A very much work in progress 16-bit register based virtual machine implemented in Go. Mostly to learn about Go as well as VMs in general.

## Example Code

The vm can be programmed in a simple assembler language. The assembler isn't yet
ready to support the full feature set of the VM.

```asm
CONST 5
CONST 6
START
LOADC 0 0
LOADC 1 1
DBG
ADD 2 0 1
DBG
ADD 2 2 2
DBG
ADD 0 1 2
HALT
```

## Architecture

### Memory

Each memory slot is 16 bits wide and can hold one instruction or data at 16 bits
width. With memory being addressable using 16 bit numbers it has a total of 65536
memory slots, giving the machine 128kB of memory.

### Memory Layout

The first block is reserved to hold the memory adress where the program starts.

The blocks between the first block and the program start are the constant pool. Any constant values used by the program can be stored here. The constant pool is readonly. Constant slots are limited to 512.

The rest of the memory is used for the program itself.

The memory block after the program's instructions can be used to as a read and write memory. They are mapped by an internal helper so they can be accessed starting at adress 0. The size of the mapped memory is limited to 512 slots.

```
0: 0x0006 - The program starts at memory adress 6
1: 0x0001 - Setting constant zero to 1
2: 0x0002 - Setting constant one to 2
3: 0x0003 - Setting constant two to 3
4: 0x0004 - Setting constant three to 4
5: 0x0005 - Setting constant four to 5
6: 0x1000 - Program starts here, Loading constant zero into R0
7: 0x0000 - Program ends here, HALT instruction
8: 0x0000 - Read/Write memory starts here
```

### Registers

Registers are addressed using 3 bits, which yields a total 8 general purpose
registers (R0–R7).

### Operations

The first four bits of an instruction hold the opcode. This allows for a maximum
of 16 instructions. The remaining 12 bits of the instruction are used for the
instructions' parameters.

```
-----------------------------------------------------------------------
| 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 |
-----------------------------------------------------------------------
|     OPCODE    |                     PARAMETERS                      |
-----------------------------------------------------------------------
```

### Why?

Mostly because I like rabbit holes. And who knows, I've been playing with building a
[lisp like language](https://github.com/lnolte/Sol) for generative design. Maybe
at some point in the future it could be run on this VM. Maybe this VM could also
be optimized to do vector based creative coding. Who knows, so many
possiblities...

## Resources

- [Write your own virtual machine](https://www.jmeiners.com/lc3-vm/)
- [A 16-bit VM in JavaScript: LowLevel Javascript](https://www.youtube.com/watch?v=fTBwD3sb5mw&list=PLP29wDx6QmW5DdwpdwHCRJsEubS5NrQ9b)
