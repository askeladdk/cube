package cube

type OpcodeType int

const (
	NOP OpcodeType = iota
	ADD
	MUL
	RET
)

type Opcode struct {
	Name string
}

var Opcodes = []Opcode{
	Opcode{"NOP"},
	Opcode{"ADD"},
	Opcode{"MUL"},
	Opcode{"RET"},
}

func (this *OpcodeType) String() string {
	return Opcodes[*this].Name
}
