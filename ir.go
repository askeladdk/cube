package cube

type instruction struct {
	Opcode   Opcode
	Operands [3]int
}

type basicBlock struct {
	Name         string
	Instructions []instruction

	sccomponent  int
	jmpcode      Opcode
	jmparg       int
	successors   [2]*basicBlock
	predecessors []*basicBlock
}

type local struct {
	Name        string
	Type        *Type
	Index       int
	Parent      int
	IsParameter bool
}

type Procedure struct {
	name       string
	returnType *Type
	constants  []uint64
	locals     []local
	blocks     []*basicBlock
	entryPoint *basicBlock
}
