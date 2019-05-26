package cube

import (
	"fmt"
	"io"
)

func printproc(w io.Writer, proc *Procedure) {
	fmt.Fprintf(w, "func %s(", proc.name)

	funargidx := 0
	for funargidx < len(proc.locals) && proc.locals[funargidx].isParameter {
		local := &proc.locals[funargidx]
		fmt.Fprintf(w, "%s, ", local)
		funargidx += 1
	}

	fmt.Fprintf(w, ") %s {\n", proc.returnType)

	for i := funargidx; i < len(proc.locals); i++ {
		local := &proc.locals[i]
		fmt.Fprintf(w, " var %s\n", local)
	}

	op2str := func(op operand) string {
		if otype, val := op.unpack(); otype == operandType_CON {
			return fmt.Sprintf("0x%x", proc.constants[val])
		} else if otype == operandType_LOC {
			return proc.locals[val].name
		} else if otype == operandType_REG {
			return fmt.Sprintf("%s", &proc.ssaregs[val])
		} else {
			return ""
		}
	}

	for _, blk := range proc.blocks {
		fmt.Fprintf(w, " %s(", blk)
		for _, val := range blk.ssaparams {
			fmt.Fprintf(w, "%s, ", &proc.ssaregs[val])
		}
		fmt.Fprintf(w, "):\n")

		for _, insr := range blk.instructions {
			fmt.Fprintf(w, "  %s ", insr.opcode)
			for _, op := range insr.operands {
				if op.otype != operandType_NIL {
					fmt.Fprintf(w, "%s, ", op2str(op))
				}
			}
			fmt.Fprintf(w, "\n")
		}

		switch blk.jmpcode {
		case opcode_RET:
			fmt.Fprintf(w, "  ret %s\n", op2str(blk.jmpretval))
		case opcode_JMP:
			fmt.Fprintf(w, "  jmp %s(", blk.successors[0])
			for _, a := range blk.jmpargs[0] {
				fmt.Fprintf(w, "%s, ", &proc.ssaregs[a])
			}
			fmt.Fprintf(w, ")\n")
		default:
			fmt.Fprintf(w, "  %s, %s(", op2str(blk.jmpretval), blk.successors[0])
			for _, a := range blk.jmpargs[0] {
				fmt.Fprintf(w, "%s, ", &proc.ssaregs[a])
			}
			fmt.Fprintf(w, "), %s(", blk.successors[1])
			for _, a := range blk.jmpargs[1] {
				fmt.Fprintf(w, "%s, ", &proc.ssaregs[a])
			}
			fmt.Fprintf(w, ")\n")
		}
	}

	fmt.Fprintf(w, "}\n")
}
