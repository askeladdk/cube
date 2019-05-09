package il

import (
	"errors"
	"fmt"

	"github.com/askeladdk/cube"
)

type localInfo struct {
	name  string
	index int
	dtype *cube.Type
	node  Node
	param bool
}

type blockInfo struct {
	name  string
	index int
	node  Node
}

type funcInfo struct {
	name   string
	level  int
	node   Node
	locals map[string]*localInfo
	blocks map[string]*blockInfo
	params []*localInfo
	rtype  *cube.Type
}

type funcInfoStack []*funcInfo

func (this funcInfoStack) push(s *funcInfo) funcInfoStack {
	return append(this, s)
}

func (this funcInfoStack) pop() (funcInfoStack, *funcInfo) {
	n := len(this)
	if n > 0 {
		return this[:n-1], this[n-1]
	} else {
		return nil, nil
	}
}

func (this funcInfoStack) peek() (*funcInfo, bool) {
	n := len(this)
	if n > 0 {
		return this[n-1], true
	} else {
		return nil, false
	}
}

type semanticAnalysis struct {
	funcs     map[string]*funcInfo
	funcStack funcInfoStack
}

func getDataType(n Node) (*cube.Type, bool) {
	switch m := n.(type) {
	case *TypeName:
		return m.Type, true
	default:
		return nil, false
	}
}

// First pass symbol resolution

type nameResolver semanticAnalysis

func (this *nameResolver) Visit(n Node) (Node, error) {
	switch m := n.(type) {
	case *Program:
		progInfo := &funcInfo{
			level: 0,
			name:  "",
			node:  n,
		}
		this.funcs[progInfo.name] = progInfo
		this.funcStack = this.funcStack.push(progInfo)
	case *Function:
		if _, exists := this.funcs[m.Name]; exists {
			return nil, errors.New(fmt.Sprintf("func '%s' exists", m.Name))
		} else {
			funcinfo := &funcInfo{
				level:  len(this.funcStack),
				name:   m.Name,
				locals: map[string]*localInfo{},
				blocks: map[string]*blockInfo{},
				node:   n,
			}
			this.funcs[m.Name] = funcinfo
			this.funcStack = this.funcStack.push(funcinfo)
			m.funcInfo = funcinfo
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
			m.localInfo = &localInfo{
				name:  m.Name,
				dtype: dtype,
				index: len(curfunc.locals),
				node:  m,
				param: true,
			}
			curfunc.locals[m.Name] = m.localInfo
			curfunc.params = append(curfunc.params, m.localInfo)
		}
	case *Def:
		curfunc, _ := this.funcStack.peek()
		if locinfo, exists := curfunc.locals[m.Name]; exists {
			m.localInfo = locinfo
		} else {
			dtype, _ := getDataType(m.TypeName)
			m.localInfo = &localInfo{
				name:  m.Name,
				dtype: dtype,
				index: len(curfunc.locals),
				node:  m,
				param: false,
			}
			curfunc.locals[m.Name] = m.localInfo
		}
	case *Use:
		curfunc, _ := this.funcStack.peek()
		if localinfo, exists := curfunc.locals[m.Name]; exists {
			m.localInfo = localinfo
		} else {
			return nil, errors.New(fmt.Sprintf("register '%s' does not exist", m.Name))
		}
	case *Block:
		curfunc, _ := this.funcStack.peek()
		if _, exists := curfunc.blocks[m.Name]; exists {
			return nil, errors.New(fmt.Sprintf("block '%s' exists", m.Name))
		} else {
			curfunc.blocks[m.Name] = &blockInfo{
				name:  m.Name,
				index: len(curfunc.blocks),
				node:  m,
			}
		}
	}
	return n, nil
}

func (this *nameResolver) PostVisit(n Node) (Node, error) {
	switch n.(type) {
	case *Program:
		this.funcStack, _ = this.funcStack.pop()
	case *Function:
		this.funcStack, _ = this.funcStack.pop()
	}
	return n, nil
}

// Second pass symbol resolution

type labelUseResolver semanticAnalysis

func (this *labelUseResolver) Visit(n Node) (Node, error) {
	switch m := n.(type) {
	case *Function:
		this.funcStack = this.funcStack.push(m.funcInfo)
	case *LabelUse:
		curfunc, _ := this.funcStack.peek()
		if blockInfo, exists := curfunc.blocks[m.Name]; exists {
			m.blockInfo = blockInfo
		} else {
			return nil, errors.New(fmt.Sprintf("block '%s' does not exist", m.Name))
		}
	}
	return n, nil
}

func (this *labelUseResolver) PostVisit(n Node) (Node, error) {
	switch n.(type) {
	case *Function:
		this.funcStack, _ = this.funcStack.pop()
	}
	return n, nil
}

// Type checker

type typeChecker semanticAnalysis

func getTypeOfNode(n Node) (*cube.Type, error) {
	switch m := n.(type) {
	case *Integer:
		return cube.TypeInt64, nil
	case *Use:
		return m.localInfo.dtype, nil
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
		} else if def.localInfo.dtype != cube.TypeAuto {
			return nil, errors.New("type error!")
		} else {
			def.localInfo.dtype = srcType
		}
	case *Instruction:
		if use, ok := m.Dst.(*Use); !ok {
			return nil, errors.New("invalid ast!")
		} else if srcAType, err := getTypeOfNode(m.SrcA); err != nil {
			return nil, err
		} else if srcBType, err := getTypeOfNode(m.SrcB); err != nil {
			return nil, err
		} else if srcAType != srcBType || use.localInfo.dtype != srcAType {
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
