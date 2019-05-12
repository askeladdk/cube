package il

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

type Def struct {
	Name     string
	TypeName Node
	symbol   *registerSymbol
}

func (this *Def) Traverse(vi Visitor) (Node, error) {
	if typename, err := vi.Visit(this.TypeName); err != nil {
		return nil, err
	} else {
		this.TypeName = typename
		return this, nil
	}
}

func (this *Def) String() string {
	return fmt.Sprintf("Def<%s>", this.Name)
}

type Use struct {
	Name   string
	symbol *registerSymbol
}

func (this *Use) Traverse(vi Visitor) (Node, error) {
	return this, nil
}

func (this *Use) String() string {
	return fmt.Sprintf("Use<%s>", this.Name)
}

type LabelUse struct {
	Name   string
	symbol *blockSymbol
}

func (this *LabelUse) Traverse(vi Visitor) (Node, error) {
	return this, nil
}

func (this *LabelUse) String() string {
	return fmt.Sprintf("LabelUse<%s>", this.Name)
}

type Parameter struct {
	Name     string
	TypeName Node
	Next     Node
	symbol   *registerSymbol
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
	Name         string
	Instructions Node
	Next         Node
	symbol       *blockSymbol
}

func (this *Block) Traverse(vi Visitor) (Node, error) {
	if instruction, err := vi.Visit(this.Instructions); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.Instructions = instruction
		this.Next = next
		return this, nil
	}
}

func (this *Block) String() string {
	return fmt.Sprintf("Block<%s>", this.Name)
}

type Signature struct {
	Parameters Node
	Returns    Node
}

func (this *Signature) Traverse(vi Visitor) (Node, error) {
	if parameters, err := vi.Visit(this.Parameters); err != nil {
		return nil, err
	} else if returns, err := vi.Visit(this.Returns); err != nil {
		return nil, err
	} else {
		this.Parameters = parameters
		this.Returns = returns
		return this, nil
	}
}

func (this *Signature) String() string {
	return "Signature<>"
}

type Function struct {
	Name      string
	Signature Node
	Blocks    Node
	Next      Node
	symbol    *funcSymbol
}

func (this *Function) Traverse(vi Visitor) (Node, error) {
	if signature, err := vi.Visit(this.Signature); err != nil {
		return nil, err
	} else if blocks, err := vi.Visit(this.Blocks); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.Signature = signature
		this.Blocks = blocks
		this.Next = next
		return this, nil
	}
}

func (this *Function) String() string {
	return fmt.Sprintf("Function<%s>", this.Name)
}

type Set struct {
	Dst  Node
	Src  Node
	Next Node
}

func (this *Set) Traverse(vi Visitor) (Node, error) {
	if dst, err := vi.Visit(this.Dst); err != nil {
		return nil, err
	} else if src, err := vi.Visit(this.Src); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.Dst = dst
		this.Src = src
		this.Next = next
		return this, nil
	}
}

func (this *Set) String() string {
	return "Set<>"
}

type Branch struct {
	Label Node
}

func (this *Branch) Traverse(vi Visitor) (Node, error) {
	if label, err := vi.Visit(this.Label); err != nil {
		return nil, err
	} else {
		this.Label = label
		return this, nil
	}
}

func (this *Branch) String() string {
	return "Branch<>"
}

type ConditionalBranch struct {
	Cond        Node
	LabelA      Node
	LabelB      Node
	OpcodeToken TokenType
}

func (this *ConditionalBranch) Traverse(vi Visitor) (Node, error) {
	if cond, err := vi.Visit(this.Cond); err != nil {
		return nil, err
	} else if labela, err := vi.Visit(this.LabelA); err != nil {
		return nil, err
	} else if labelb, err := vi.Visit(this.LabelB); err != nil {
		return nil, err
	} else {
		this.Cond = cond
		this.LabelA = labela
		this.LabelB = labelb
		return this, nil
	}
}

func (this *ConditionalBranch) String() string {
	return fmt.Sprintf("ConditionalBranch<%d>", this.OpcodeToken)
}

type Return struct {
	Src Node
}

func (this *Return) Traverse(vi Visitor) (Node, error) {
	if src, err := vi.Visit(this.Src); err != nil {
		return nil, err
	} else {
		this.Src = src
		return this, nil
	}
}

func (this *Return) String() string {
	return "Return<>"
}

type Instruction struct {
	Dst         Node
	SrcA        Node
	SrcB        Node
	Next        Node
	OpcodeToken TokenType
	Opcode      *cube.OpcodeType
}

func (this *Instruction) Traverse(vi Visitor) (Node, error) {
	if dst, err := vi.Visit(this.Dst); err != nil {
		return nil, err
	} else if srca, err := vi.Visit(this.SrcA); err != nil {
		return nil, err
	} else if srcb, err := vi.Visit(this.SrcB); err != nil {
		return nil, err
	} else if next, err := vi.Visit(this.Next); err != nil {
		return nil, err
	} else {
		this.Dst = dst
		this.SrcA = srca
		this.SrcB = srcb
		this.Next = next
		return this, nil
	}
}

func (this *Instruction) String() string {
	return fmt.Sprintf("Instruction<%d>", this.OpcodeToken)
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
