package main

import (
  "fmt"
  "os"
  "vm/instructions"
  "vm/assembler"
)

const MEMORY_MAX = (1 << 16)
const INSTRUCTION_SIZE = 16
const OPCODE_SIZE = 4
const PARAMETER_SIZE = 12
const CONSTANT_POOL_OFFSET = 1

/**
 * MEMORY
 * =============================================================================
 *
 * The virutal machine has 65536 (2^16) memory locations, each of which can
 * hold a 16-bit value.
 *
 * This means we have a total memory of 128kB
 */
var memory [MEMORY_MAX]uint16

// Writes to the literal memory at the given address
func lit_mem_write(address uint16, value uint16) {
    memory[address] = value
}

// Reads the literal memory at the given address
func lit_mem_read(address uint16) uint16 {
  return memory[address]
}

func map_mem_write(address uint16, value uint16) {
  memory[reg[R_MAR]] = value
}

func map_mem_read(address uint16) uint16 {
  return memory[reg[R_MAR] + address]
}


/**
 * MEMORY LAYOUT
 * =============================================================================
 *
 * The first block is reversed to hold the memory adress where the program
 * starts.
 *
 * The blocks between the first block and the program start are the constant
 * pool. Any constant values used by the program can be stored here. The
 * constant pool is readonly. Constant slots are limited to 512.
 *
 * The rest of the memory is used for the program itself.

 * The memory block after the program's instructions can be used to as a read
 * and write memory. They are mapped by an internal helper so they can be
 * accessed starting at adress 0. The size of the mapped memory is limited to
 * 512 slots.
 *
 * 0: 0x0006 - The program starts at memory adress 6
 * 1: 0x0001 - Setting constant zero to 1
 * 2: 0x0002 - Setting constant one to 2
 * 3: 0x0003 - Setting constant two to 3
 * 4: 0x0004 - Setting constant three to 4
 * 5: 0x0005 - Setting constant four to 5
 * 6: 0x1000 - Program starts here, Loading constant zero into R0
 * 7: 0x0000 - Program ends here, HALT instruction
 * 8: 0x0000 - Read/Write memory starts here
 * ...
 */

/**
 * CONSTANT POOL
 * =============================================================================
 * 
 * The constant pool is the area of memory before the program that holds
 * constants that can be loaded using the LOADC instruction.
 *
 */
func const_read(adress uint16) uint16 {
  return memory[adress + CONSTANT_POOL_OFFSET]
}

/**
 * REGISTERS
 * =============================================================================
 *
 * The virtual machine has 10 total registers. 
 * 8 of them are general purpose registers (R0-R7)
 *
 * Each register is 16 bits wide.
 */
const (
  R_R0 = 0x00
  R_R1 = 0x01
  R_R2 = 0x02
  R_R3 = 0x03
  R_R4 = 0x04
  R_R5 = 0x05
  R_R6 = 0x06
  R_R7 = 0x07
  R_PC = 0x08   /* program counter */
  R_COND = 0x09 /* condition flags */
  R_MAR = 0x0A  /* memory address register */
  R_COUNT = 0x0B
)

var reg [R_COUNT]uint16

/**
 * CONDITION FLAGS
 * =============================================================================
 *
 * The R_COND register stores condition flags. These hold information about the
 * most recent calculation. This allows programs to check for logical
 * conditions.
 */
const (
    FL_POS = 1 << 0 /* Positive */
    FL_ZRO = 1 << 1 /* Zero */
    FL_NEG = 1 << 2 /* Negative */
)

/**
 * UTILITY FUNCTIONS
 * =============================================================================
 */
func sign_extend(x uint16, bit_count int) uint16 {
  if (x >> (bit_count - 1)) & 1 == 1 {
    x |= (0xFFFF << bit_count)
  }
  return x
}

func update_flags(r uint16) {
  if (reg[r] == 0) {
    reg[R_COND] = FL_ZRO
  } else if (reg[r] >> 15) == 1 {
    reg[R_COND] = FL_NEG
  } else {
    reg[R_COND] = FL_POS
  }
}

func load_into_memory(program []uint16) {
  // Load the program into memory
  for i, instruction := range program {
    memory[uint16(i)] = instruction
  }
}

/**
 * MAIN LOOP
 * =============================================================================
 */
func main() {

  if len(os.Args) < 2 {
    fmt.Println("vm [program file]")
    os.Exit(2)
  }

  prog_file := os.Args[1]
  fmt.Println("Loading program from", prog_file)

  data, err := os.ReadFile(prog_file)
  if err != nil {
    fmt.Println("Error reading file", err)
    os.Exit(1)
  }

  load_into_memory(assembler.Assemble(string(data)))

  reg[R_COND] = FL_ZRO
  reg[R_PC] = memory[0x0000]
  reg[R_MAR] = 0x1007 // TODO: This should be set by the VM after loading a ROM

  var running bool = true

  for running {
    var instr uint16 = lit_mem_read(reg[R_PC])

    var op uint16 = instr >> PARAMETER_SIZE
    reg[R_PC]++

    switch op {
      case instructions.OP_HALT:
        running = false
        break;

      case instructions.OP_LOADC:

        // PUSH INSTRUCTION
        // 
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |     OP_LOADC      |     REG     |               CONST9              |
        // -----------------------------------------------------------------------

        r1 := (instr >> 9) & 0x7
        c_offset := sign_extend(instr & 0x1FF, 9)

        reg[r1] = const_read(c_offset)
        break

      case instructions.OP_MOVE:

        // MOVE INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |     OP_MOVE       |     REG     |    REG    |                       |
        // -----------------------------------------------------------------------

        r1 := (instr >> 9) & 0x7
        r2 := (instr >> 6) & 0x7

        reg[r2] = reg[r1]
        break

      case instructions.OP_LOADM:

        // LOADM INSTRUCTION
        // 
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |     OP_LOADC      |     REG     |               CONST9              |
        // -----------------------------------------------------------------------

        r1 := (instr >> 9) & 0x7
        m_offset := sign_extend(instr & 0x1FF, 9)

        reg[r1] = map_mem_read(m_offset)
        break

      case instructions.OP_STOREM:

        // STOREM INSTRUCTION
        // 
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |     OP_LOADC      |     REG     |               CONST9              |
        // -----------------------------------------------------------------------


        r1 := (instr >> 9) & 0x7
        m_offset := sign_extend(instr & 0x1FF, 9)

        map_mem_write(m_offset, reg[r1])
        break

      case instructions.OP_JUMP:

        // JUMP INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |     OP_JUMP       | si |                   OFFSET                   |
        // -----------------------------------------------------------------------

        si := (instr >> 11) & 0x1
        offset := instr & 0x3FF

        if si == 0 {
          reg[R_PC] = reg[R_PC] + offset
        } else {
          reg[R_PC] = reg[R_PC] - offset
        }

        break

      case instructions.OP_ADD:

        // ADD INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 0 |    REG    |       |    REG    |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 1 |    REG    |        IMM5       |
        // -----------------------------------------------------------------------

        dr := (instr >> 9) & 0x7
        r1 := (instr >> 5) & 0x7
        imm_flag := (instr >> 8) & 0x1

        if imm_flag == 1 {
          imm5 := sign_extend(instr & 0x1F, 5)
          reg[dr] = reg[r1] + imm5
        } else {
          r2 := instr & 0x7
          reg[dr] = reg[r1] + reg[r2]
        }
        break

      case instructions.OP_SUB:

        // SUB INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 0 |    REG    |       |    REG    |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 1 |    REG    |        IMM5       |
        // -----------------------------------------------------------------------

        dr := (instr >> 9) & 0x7
        r1 := (instr >> 5) & 0x7
        imm_flag := (instr >> 8) & 0x1

        if imm_flag == 1 {
          imm5 := sign_extend(instr & 0x1F, 5)
          reg[dr] = reg[r1] - imm5
        } else {
          r2 := instr & 0x7
          reg[dr] = reg[r1] - reg[r2]
        }

        break

      case instructions.OP_MUL:

        // MUL INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 0 |    REG    |       |    REG    |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 1 |    REG    |        IMM5       |
        // -----------------------------------------------------------------------

        dr := (instr >> 9) & 0x7
        r1 := (instr >> 5) & 0x7
        imm_flag := (instr >> 8) & 0x1

        if imm_flag == 1 {
          imm5 := sign_extend(instr & 0x1F, 5)
          reg[dr] = reg[r1] * imm5
        } else {
          r2 := instr & 0x7
          reg[dr] = reg[r1] * reg[r2]
        }
        break

      case instructions.OP_DIV:

        // DIV INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 0 |    REG    |       |    REG    |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 1 |    REG    |        IMM5       |
        // -----------------------------------------------------------------------

        dr := (instr >> 9) & 0x7
        r1 := (instr >> 5) & 0x7
        imm_flag := (instr >> 8) & 0x1

        if imm_flag == 1 {
          imm5 := sign_extend(instr & 0x1F, 5)
          reg[dr] = reg[r1] / imm5
        } else {
          r2 := instr & 0x7
          reg[dr] = reg[r1] / reg[r2]
        }
        break

      case instructions.OP_NOT:

        // NOT INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_NOT       |     REG     |    REG    |                       |
        // -----------------------------------------------------------------------

        dr := (instr >> 9) & 0x7
        sr := (instr >> 6) & 0x7

        reg[dr] = ^reg[sr]

        break

      case instructions.OP_EQ:

        // EQ INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 0 |                   |    REG    |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 1 |             IMM8              |
        // -----------------------------------------------------------------------

        r1 := (instr >> 9) & 0x7
        imm_flag := (instr >> 8) & 0x1

        var is_eq bool

        if imm_flag == 1 {
          imm8 := sign_extend(instr & 0xFF, 8)
          is_eq = reg[r1] == imm8
        } else {
          r2 := instr & 0x7
          is_eq = reg[r1] == reg[r2]
        }

        if !is_eq {
          reg[R_PC]++
        }

        break

      case instructions.OP_LT:

        // LT INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 0 |                   |    REG    |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 1 |             IMM8              |
        // -----------------------------------------------------------------------

        r1 := (instr >> 9) & 0x7
        imm_flag := (instr >> 8) & 0x1

        var is_eq bool

        if imm_flag == 1 {
          imm8 := sign_extend(instr & 0xFF, 8)
          is_eq = reg[r1] < imm8
        } else {
          r2 := instr & 0x7
          is_eq = reg[r1] < reg[r2]
        }

        if !is_eq {
          reg[R_PC]++
        }

        break

      case instructions.OP_LE:

        // LE INSTRUCTION
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 0 |                   |    REG    |
        // -----------------------------------------------------------------------
        // |      OP_ADD       |     REG     | 1 |             IMM8              |
        // -----------------------------------------------------------------------

        r1 := (instr >> 9) & 0x7
        imm_flag := (instr >> 8) & 0x1

        var is_eq bool

        if imm_flag == 1 {
          imm8 := sign_extend(instr & 0xFF, 8)
          is_eq = reg[r1] < imm8
        } else {
          r2 := instr & 0x7
          is_eq = reg[r1] < reg[r2]
        }

        if !is_eq {
          reg[R_PC]++
        }

        break

      case instructions.OP_DBG:

        // DBG INSTRUCTION
        // 
        // Prints the current value of general purpose registers to stdout
        //
        // -----------------------------------------------------------------------
        // | 15 | 14 | 13 | 12 | 11 | 10 | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
        // -----------------------------------------------------------------------
        // |      OP_DBG       |                                                 |
        // -----------------------------------------------------------------------

        fmt.Printf("PC\tR0\tR1\tR2\tR3\tR4\tR5\tR6\tR7\n")
        fmt.Println("--------------------------------------------------------------------")
        fmt.Printf("%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", reg[R_PC], reg[0], reg[1], reg[2], reg[3], reg[4], reg[5], reg[6], reg[7])
        fmt.Println("--------------------------------------------------------------------")
        fmt.Println("")
        break;

      default:
        fmt.Printf("Unknown instruction: %x\n", op)
        running = false
        os.Exit(1)
        break
    }
  }
}
