package cube

type OperandType int

const (
	OPNONE = iota
	OPLABEL
	OPLOCAL
	OPIMMEDIATE
)

type OpcodeType struct {
	Name         string
	OperandTypes [3]OperandType
}

func (this *OpcodeType) String() string {
	return this.Name
}

var (
	Opcode_ADD = &OpcodeType{"ADD", [3]OperandType{OPLOCAL, OPLOCAL, OPLOCAL}}
	Opcode_ADI = &OpcodeType{"ADI", [3]OperandType{OPLOCAL, OPLOCAL, OPIMMEDIATE}}
	Opcode_JMP = &OpcodeType{"JMP", [3]OperandType{OPLOCAL, OPLABEL, OPNONE}}
	Opcode_JNZ = &OpcodeType{"JNZ", [3]OperandType{OPLOCAL, OPLABEL, OPLABEL}}
	Opcode_MUL = &OpcodeType{"MUL", [3]OperandType{OPLOCAL, OPLOCAL, OPLOCAL}}
	Opcode_NOP = &OpcodeType{"NOP", [3]OperandType{OPNONE, OPNONE, OPNONE}}
	Opcode_REI = &OpcodeType{"REI", [3]OperandType{OPIMMEDIATE, OPNONE, OPNONE}}
	Opcode_RET = &OpcodeType{"RET", [3]OperandType{OPLOCAL, OPNONE, OPNONE}}
	Opcode_SET = &OpcodeType{"SET", [3]OperandType{OPLOCAL, OPLOCAL, OPLOCAL}}
	Opcode_SUB = &OpcodeType{"SUB", [3]OperandType{OPLOCAL, OPLOCAL, OPLOCAL}}
)
