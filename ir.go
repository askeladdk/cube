package cube

type Instruction struct {
	opcode   Opcode
	operands [3]int
}

type BasicBlock struct {
	name         string
	instructions []Instruction

	sccomponent  int
	jmpcode      Opcode
	jmparg       int
	successors   [2]*BasicBlock
	predecessors []*BasicBlock
}

type Local struct {
	name        string
	dataType    *Type
	index       int
	parent      int
	isParameter bool
}

type Procedure struct {
	name       string
	returnType *Type
	constants  []uint64
	locals     []Local
	blocks     []*BasicBlock
	entryPoint *BasicBlock
}
