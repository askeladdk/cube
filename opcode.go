package cube

type opcode struct {
	name string
}

func (this *opcode) String() string {
	return this.name
}

var (
	opcode_ADD = &opcode{"add"}
	opcode_JMP = &opcode{"jmp"}
	opcode_JNZ = &opcode{"jnz"}
	opcode_MOV = &opcode{"mov"}
	opcode_MUL = &opcode{"mul"}
	opcode_RET = &opcode{"ret"}
	opcode_SUB = &opcode{"sub"}
)

type operandType int

const (
	operandType_NIL operandType = iota
	operandType_LOC
	operandType_REG
	operandType_CON
	operandType_BLK
)

type operand struct {
	otype operandType
	value int
}

func newOperand(otype operandType, value int) operand {
	return operand{otype, value}
}

func (this operand) unpack() (operandType, int) {
	return this.otype, this.value
}

var operandNil = newOperand(operandType_NIL, 0)

func operandLoc(value int) operand {
	return newOperand(operandType_LOC, value)
}

func operandReg(value int) operand {
	return newOperand(operandType_REG, value)
}

func operandCon(value int) operand {
	return newOperand(operandType_CON, value)
}
