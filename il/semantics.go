package il

import (
	"errors"
	"fmt"

	"github.com/askeladdk/cube"
)

type registerSymbol struct {
	name  string
	index int
	dtype *cube.Type
	node  Node
	param bool
}

type blockSymbol struct {
	name  string
	index int
	node  Node
}

type funcSymbol struct {
	name   string
	level  int
	node   Node
	locals map[string]*registerSymbol
	blocks map[string]*blockSymbol
	params []*registerSymbol
	rtype  *cube.Type
}

type funcInfoStack []*funcSymbol

func (this funcInfoStack) push(s *funcSymbol) funcInfoStack {
	return append(this, s)
}

func (this funcInfoStack) pop() (funcInfoStack, *funcSymbol) {
	n := len(this)
	if n > 0 {
		return this[:n-1], this[n-1]
	} else {
		return nil, nil
	}
}

func (this funcInfoStack) peek() (*funcSymbol, bool) {
	n := len(this)
	if n > 0 {
		return this[n-1], true
	} else {
		return nil, false
	}
}

type symbolTable map[string]*funcSymbol

func getDataType(n Node) (*cube.Type, bool) {
	switch m := n.(type) {
	case *TypeName:
		return m.Type, true
	default:
		return nil, false
	}
}

// First pass symbol resolution

type symbolResolver struct {
	symbolTable symbolTable
	funcStack   funcInfoStack
}

func (this *symbolResolver) Visit(n Node) (Node, error) {
	switch m := n.(type) {
	case *Program:
		progInfo := &funcSymbol{
			level: 0,
			name:  "",
			node:  n,
		}
		this.symbolTable[progInfo.name] = progInfo
		this.funcStack = this.funcStack.push(progInfo)
	case *Function:
		if _, exists := this.symbolTable[m.Name]; exists {
			return nil, errors.New(fmt.Sprintf("func '%s' exists", m.Name))
		} else {
			funcinfo := &funcSymbol{
				level:  len(this.funcStack),
				name:   m.Name,
				locals: map[string]*registerSymbol{},
				blocks: map[string]*blockSymbol{},
				node:   n,
			}
			this.symbolTable[m.Name] = funcinfo
			this.funcStack = this.funcStack.push(funcinfo)
			m.symbol = funcinfo
		}
	case *Signature:
		curfunc, _ := this.funcStack.peek()
		curfunc.rtype, _ = getDataType(m.Returns)
	case *Parameter:
		curfunc, _ := this.funcStack.peek()
		if _, exists := curfunc.locals[m.Name]; exists {
			return nil, errors.New(fmt.Sprintf("parameter '%s' exists", m.Name))
		} else {
			dtype, _ := getDataType(m.TypeName)
			m.symbol = &registerSymbol{
				name:  m.Name,
				dtype: dtype,
				index: len(curfunc.locals),
				node:  m,
				param: true,
			}
			curfunc.locals[m.Name] = m.symbol
			curfunc.params = append(curfunc.params, m.symbol)
		}
	case *Def:
		curfunc, _ := this.funcStack.peek()
		if locinfo, exists := curfunc.locals[m.Name]; exists {
			m.symbol = locinfo
		} else {
			dtype, _ := getDataType(m.TypeName)
			m.symbol = &registerSymbol{
				name:  m.Name,
				dtype: dtype,
				index: len(curfunc.locals),
				node:  m,
				param: false,
			}
			curfunc.locals[m.Name] = m.symbol
		}
	case *Use:
		curfunc, _ := this.funcStack.peek()
		if localinfo, exists := curfunc.locals[m.Name]; exists {
			m.symbol = localinfo
		} else {
			return nil, errors.New(fmt.Sprintf("register '%s' does not exist", m.Name))
		}
	case *Block:
		curfunc, _ := this.funcStack.peek()
		if _, exists := curfunc.blocks[m.Name]; exists {
			return nil, errors.New(fmt.Sprintf("block '%s' exists", m.Name))
		} else {
			curfunc.blocks[m.Name] = &blockSymbol{
				name:  m.Name,
				index: len(curfunc.blocks),
				node:  m,
			}
		}
	}
	return n, nil
}

func (this *symbolResolver) PostVisit(n Node) (Node, error) {
	switch n.(type) {
	case *Program:
		this.funcStack, _ = this.funcStack.pop()
	case *Function:
		this.funcStack, _ = this.funcStack.pop()
	}
	return n, nil
}

// Second pass symbol resolution

type symbolBacklinkResolver symbolResolver

func (this *symbolBacklinkResolver) Visit(n Node) (Node, error) {
	switch m := n.(type) {
	case *Function:
		this.funcStack = this.funcStack.push(m.symbol)
	case *LabelUse:
		curfunc, _ := this.funcStack.peek()
		if blockInfo, exists := curfunc.blocks[m.Name]; exists {
			m.symbol = blockInfo
		} else {
			return nil, errors.New(fmt.Sprintf("block '%s' does not exist", m.Name))
		}
	}
	return n, nil
}

func (this *symbolBacklinkResolver) PostVisit(n Node) (Node, error) {
	switch n.(type) {
	case *Function:
		this.funcStack, _ = this.funcStack.pop()
	}
	return n, nil
}

// Type checker

type typeChecker struct {
	symbolTable symbolTable
}

func getTypeOfNode(n Node) (*cube.Type, error) {
	switch m := n.(type) {
	case *Integer:
		return cube.TypeInt64, nil
	case *Use:
		return m.symbol.dtype, nil
	default:
		return nil, errors.New("not a typed node")
	}
}

var tokenTypeToOpcodeType_i64 = map[TokenType]*cube.OpcodeType{
	SET: cube.Opcode_SET_I64,
	ADD: cube.Opcode_ADD_I64,
	SUB: cube.Opcode_SUB_I64,
	MUL: cube.Opcode_MUL_I64,
	RET: cube.Opcode_RET_I64,
}

func (this *typeChecker) Visit(n Node) (Node, error) {
	switch m := n.(type) {
	case *Set:
		if srcType, err := getTypeOfNode(m.Src); err != nil {
			return nil, err
		} else if def, ok := m.Dst.(*Def); !ok {
			return nil, errors.New("invalid ast!")
		} else if def.symbol.dtype != cube.TypeAuto {
			return nil, errors.New("type error!")
		} else {
			def.symbol.dtype = srcType
		}
	case *Instruction:
		if use, ok := m.Dst.(*Use); !ok {
			return nil, errors.New("invalid ast!")
		} else if srcAType, err := getTypeOfNode(m.SrcA); err != nil {
			return nil, err
		} else if srcBType, err := getTypeOfNode(m.SrcB); err != nil {
			return nil, err
		} else if srcAType != srcBType || use.symbol.dtype != srcAType {
			return nil, errors.New("source operand types do not match")
		} else {
			m.Opcode, _ = tokenTypeToOpcodeType_i64[m.OpcodeToken]
		}
	}

	return n, nil
}

func (this *typeChecker) PostVisit(n Node) (Node, error) {
	return n, nil
}
