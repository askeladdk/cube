package cube

type Instruction struct {
	opcode   opcode
	operands [3]operand
}

type BasicBlock struct {
	name         string
	instructions []Instruction

	sccomponent  int
	ssaparams    []int
	jmpcode      opcode
	jmpretarg    int
	jmpargs      [2][]int
	successors   [2]*BasicBlock
	predecessors []*BasicBlock
}

type Local struct {
	name        string
	dataType    *Type
	generations int
	lastssareg  int
	isParameter bool
}

type SSAReg struct {
	local      *Local
	generation int
}

type Procedure struct {
	name       string
	returnType *Type
	constants  []uint64
	locals     []Local
	ssaregs    []SSAReg
	blocks     []*BasicBlock
	entryPoint *BasicBlock
}
