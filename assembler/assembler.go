package assembler

import (
  "fmt"
  "strings"
  "strconv"
  "vm/instructions"
)

// TODO: The assembler is far from complete
func Assemble(code string) []uint16 {
  output := []uint16{}
  lines := strings.Split(code, "\n")

  var prog_start int = 0

  for i, line := range lines {
    tokens := strings.Split(line, " ")
    instr := tokens[0]

    if (len(line) == 0) {
      continue
    }

    switch instr {
      case "CONST":
        i, _ := strconv.Atoi(tokens[1])
        output = append(output, 0x0000 | uint16(i))
        break;
      case "START":
        prog_start = i + 1
      case "LOADC":
        reg, _ := strconv.Atoi(tokens[1])
        val, _ := strconv.Atoi(tokens[2])
        op := instructions.OP_LOADC << 12 | reg << 9 | val
        output = append(output, uint16(op))
        break;
      case "ADD":
        reg0, _ := strconv.Atoi(tokens[1])
        reg1, _ := strconv.Atoi(tokens[2])
        reg2, _ := strconv.Atoi(tokens[3])
        op := instructions.OP_ADD << 12 | reg0 << 9 | 0 << 8 | reg1 << 5 | reg2
        output = append(output, uint16(op))
        break;
      case "DBG":
        op := instructions.OP_DBG << 12
        output = append(output, uint16(op))
        break;
      case "HALT":
        op := instructions.OP_HALT << 12
        output = append(output, uint16(op))
        break;
      default:
        fmt.Println("Unknown instruction: ", instr)
        break;
      }
  }

  output = append(output, 0x0000)
  copy(output[1:], output)
  output[0] = uint16(prog_start)

  return output
}
