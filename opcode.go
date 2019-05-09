package cube

type OpcodeType struct {
	Name string
}

var (
	Opcode_NOP     = &OpcodeType{"NOP"}
	Opcode_SET_I64 = &OpcodeType{"SET_I64"}
	Opcode_ADD_I64 = &OpcodeType{"ADD_I64"}
	Opcode_SUB_I64 = &OpcodeType{"SUB_I64"}
	Opcode_MUL_I64 = &OpcodeType{"MUL_I64"}
	Opcode_RET_I64 = &OpcodeType{"RET_I64"}
	Opcode_GOTO    = &OpcodeType{"GOTO"}
	Opcode_IFZ     = &OpcodeType{"IFZ"}
)

func (this *OpcodeType) String() string {
	return this.Name
}
