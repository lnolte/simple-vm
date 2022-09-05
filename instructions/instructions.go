package instructions

/**
 * INSTRUCTIONS
 * =============================================================================
 *
 * The virtual machine has 16 instructions. Each instruction is 16 bits wide.
 * The first four bits are the "opcode" (or instruction type). The other 12
 * bits hold the parameters.
 *
 * -----------------------------------------------------------------------
 * | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 |
 * -----------------------------------------------------------------------
 * |     OPCODE    |                     PARAMETERS                      |
 * -----------------------------------------------------------------------
 */

 // 1. [x] LOADC – Load Constant into register
 // 2. [x] MOVE - Move value from one register to another
 // 3. [x] LOADM - Load memore value into register
 // 4. [x] STOREM - Store value into memory

 // 5. [x] JUMP - Jump to a new location in memory

 // 6. [x] ADD - Add two registers and store the result in a register
 // 7. [x] SUB - Subtract two registers and store the result in a register
 // 8. [x] MUL - Multiply two registers and store the result in a register
 // 9. [x] DIV - Divide two registers and store the result in a register

 // 10. [x] NOT - Bitwise NOT
 // 11. [x] EQ  - Check if two registers are equal if: continue else: PC++
 // 12. [ ] LT  - Check if register A is less than register B if: continue else: PC++
 // 13. [ ] LE  - Check if register A is less than or equal to register B if: continue else: PC++

const (
    OP_HALT    = 0x0  /* Halt the program */
    OP_LOADC   = 0x1  /* LOADC */
    OP_MOVE    = 0x2  /* MOVE */
    OP_LOADM   = 0x3  /* LOADM */
    OP_STOREM  = 0x4  /* STOREM */
    OP_JUMP    = 0x5  /* JUMP */
    OP_ADD     = 0x6  /* ADD */
    OP_SUB     = 0x7  /* SUB */
    OP_MUL     = 0x8  /* MUL */
    OP_DIV     = 0x9  /* DIV */
    OP_NOT     = 0xA  /* NOT */
    OP_EQ      = 0xB  /* EQ */
    OP_LT      = 0xC  /* LT */
    OP_LE      = 0xD  /* LE */
    OP_DBG     = 0xE  /* DBG */
)

