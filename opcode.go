package cube

type Opcode int

const (
	Opcode_ADD Opcode = iota
	Opcode_ADDI
	Opcode_JMP
	Opcode_JNZ
	Opcode_MUL
	Opcode_RET
	Opcode_SET
	Opcode_SUB
)
