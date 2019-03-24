package cube

type OpcodeType int

const (
	NOP OpcodeType = iota
	SET
	ADD
	SUB
	MUL
	RET
	GOTO
	IFZ
)

type Opcode struct {
	Name string
}

var Opcodes = []Opcode{
	Opcode{"NOP"},
	Opcode{"SET"},
	Opcode{"ADD"},
	Opcode{"SUB"},
	Opcode{"MUL"},
	Opcode{"RET"},
	Opcode{"GOTO"},
	Opcode{"IFZ"},
}

func (this *OpcodeType) String() string {
	return Opcodes[*this].Name
}
