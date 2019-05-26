package cube

import "fmt"

type Instruction struct {
	opcode   *opcode
	operands [3]operand
}

type BasicBlock struct {
	name         string
	instructions []Instruction

	sccomponent  int
	ssaparams    []int
	jmpcode      *opcode
	jmpretval    operand
	jmpargs      [2][]int
	successors   [2]*BasicBlock
	predecessors []*BasicBlock
}

func (this *BasicBlock) String() string {
	return this.name
}

type Local struct {
	name        string
	dataType    *Type
	generations int
	lastssareg  int
	isParameter bool
	isDefined   bool
}

func (this *Local) String() string {
	return fmt.Sprintf("%s %s", this.name, this.dataType)
}

type SSAReg struct {
	local      *Local
	generation int
}

func (this *SSAReg) String() string {
	return fmt.Sprintf("%s%d", this.local.name, this.generation)
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
