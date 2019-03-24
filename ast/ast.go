package ast

import (
	"fmt"

	"github.com/askeladdk/cube"
)

type Node interface {
	fmt.Stringer
	Traverse(Visitor) (Node, error)
}

type TypeName struct {
	Type *cube.Type
}

func (this *TypeName) Traverse(vi Visitor) (Node, error) {
	return this, nil
}

func (this *TypeName) String() string {
	return fmt.Sprintf("TypeName<%s>", this.Type)
}

type Integer struct {
	Value int64
}

func (this *Integer) Traverse(vi Visitor) (Node, error) {
	return this, nil
}

func (this *Integer) String() string {
	return fmt.Sprintf("Integer<%v>", this.Value)
}

type Identifier struct {
	Name string
}

func (this *Identifier) Traverse(vi Visitor) (Node, error) {
	return this, nil
}

func (this *Identifier) String() string {
	return fmt.Sprintf("Identifier<%s>", this.Name)
}

type Parameter struct {
	Name     string
	TypeName Node
	Next     Node
}

func (this *Parameter) Traverse(vi Visitor) (Node, error) {
	if typename, err := vi.Visit(this.TypeName); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.Next = next
		this.TypeName = typename
		return this, nil
	}
}

func (this *Parameter) String() string {
	return fmt.Sprintf("Parameter<%s>", this.Name)
}

type Block struct {
	Name string
	// Parameters   Node
	Instructions Node
	Next         Node
}

func (this *Block) Traverse(vi Visitor) (Node, error) {
	// if parameter, err := vi.Visit(this.Parameters); err != nil {
	//     return nil, err
	if instruction, err := vi.Visit(this.Instructions); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		// this.Parameters = parameter
		this.Instructions = instruction
		this.Next = next
		return this, nil
	}
}

func (this *Block) String() string {
	return fmt.Sprintf("Block<%s>", this.Name)
}

type Function struct {
	Name       string
	Parameters Node
	Returns    Node
	Blocks     Node
	Next       Node
}

func (this *Function) Traverse(vi Visitor) (Node, error) {
	if parameter, err := vi.Visit(this.Parameters); err != nil {
		return nil, err
	} else if returns, err := vi.Visit(this.Returns); err != nil {
		return nil, err
	} else if blocks, err := vi.Visit(this.Blocks); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.Parameters = parameter
		this.Returns = returns
		this.Blocks = blocks
		this.Next = next
		return this, nil
	}
}

func (this *Function) String() string {
	return fmt.Sprintf("Function<%s>", this.Name)
}

type Instruction struct {
	OpcodeType cube.OpcodeType
	OpA        Node
	OpB        Node
	OpC        Node
	Next       Node
}

func (this *Instruction) Traverse(vi Visitor) (Node, error) {
	if opa, err := vi.Visit(this.OpA); err != nil {
		return nil, err
	} else if opb, err := vi.Visit(this.OpB); err != nil {
		return nil, err
	} else if opc, err := vi.Visit(this.OpC); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.OpA = opa
		this.OpB = opb
		this.OpC = opc
		this.Next = next
		return this, nil
	}
}

func (this *Instruction) String() string {
	return fmt.Sprintf("Instruction<%s>", &this.OpcodeType)
}

type Error struct {
	Message string
	Node    Node
}

func (this *Error) Traverse(vi Visitor) (Node, error) {
	return this, nil
}

func (this *Error) String() string {
	return fmt.Sprintf("Error<%s>", this.Message)
}

type Unit struct {
	Filename    string
	Definitions Node
	Next        Node
}

func (this *Unit) Traverse(vi Visitor) (Node, error) {
	if definitions, err := vi.Visit(this.Definitions); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.Definitions = definitions
		this.Next = next
		return this, nil
	}
}

func (this *Unit) String() string {
	return fmt.Sprintf("Unit<%s>", this.Filename)
}

type Program struct {
	Units Node
}

func (this *Program) Traverse(vi Visitor) (Node, error) {
	if units, err := vi.Visit(this.Units); err != nil {
		return nil, err
	} else {
		this.Units = units
		return this, nil
	}
}

func (this *Program) String() string {
	return "Program<>"
}
