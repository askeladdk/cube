package cube

type opcode int

const (
	opcode_ADD opcode = iota
	opcode_JMP
	opcode_JNZ
	opcode_MUL
	opcode_RET
	opcode_SET
	opcode_SUB
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
