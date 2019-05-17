package cube

type OperandType int

const (
	OPNIL = iota
	OPLAB
	OPLOC
	OPIMM
)

type OpcodeType struct {
	Name         string
	OperandTypes [3]OperandType
}

func (this *OpcodeType) String() string {
	return this.Name
}

var (
	Opcode_ADD  = &OpcodeType{"add", [3]OperandType{OPLOC, OPLOC, OPLOC}}
	Opcode_ADDI = &OpcodeType{"addi", [3]OperandType{OPLOC, OPLOC, OPIMM}}
	Opcode_JMP  = &OpcodeType{"jmp", [3]OperandType{OPLOC, OPLAB, OPNIL}}
	Opcode_JNZ  = &OpcodeType{"jnz", [3]OperandType{OPLOC, OPLAB, OPLAB}}
	Opcode_MUL  = &OpcodeType{"mul", [3]OperandType{OPLOC, OPLOC, OPLOC}}
	Opcode_NOP  = &OpcodeType{"nop", [3]OperandType{OPNIL, OPNIL, OPNIL}}
	Opcode_RET  = &OpcodeType{"ret", [3]OperandType{OPLOC, OPNIL, OPNIL}}
	Opcode_RETI = &OpcodeType{"reti", [3]OperandType{OPIMM, OPNIL, OPNIL}}
	Opcode_SET  = &OpcodeType{"set", [3]OperandType{OPLOC, OPLOC, OPLOC}}
	Opcode_SUB  = &OpcodeType{"sub", [3]OperandType{OPLOC, OPLOC, OPLOC}}
)
