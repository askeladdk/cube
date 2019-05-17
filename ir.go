package cube

type operand uint64

type instruction struct {
	opcode   *OpcodeType
	operands [3]operand
}

type block struct {
	name  string
	index int
	insrs []instruction
}

type local struct {
	name   string
	dtype  *Type
	index  int
	parent int
	param  bool
}

type function struct {
	name   string
	rtype  *Type
	locals []local
	blocks []block
	index  int
}

type program struct {
	funcs []function
}
