package il

import (
	"errors"
	"fmt"
	"strings"

	"github.com/askeladdk/cube"
)

type funcSymbolStack []*funcSymbol

func (this funcSymbolStack) push(s *funcSymbol) funcSymbolStack {
	return append(this, s)
}

func (this funcSymbolStack) pop() (funcSymbolStack, *funcSymbol) {
	n := len(this)
	if n > 0 {
		return this[:n-1], this[n-1]
	} else {
		return nil, nil
	}
}

func (this funcSymbolStack) peek() (*funcSymbol, bool) {
	n := len(this)
	if n > 0 {
		return this[n-1], true
	} else {
		return nil, false
	}
}

func getDataType(n Node) (*cube.Type, bool) {
	switch m := n.(type) {
	case *TypeName:
		return m.Type, true
	default:
		return nil, false
	}
}

// Symbol resolution

type unresolvedLabels map[string][]*LabelUse

type symbolResolver struct {
	symbolTable      symbolTable
	funcStack        funcSymbolStack
	locals           map[string]*localSymbol
	blocks           map[string]*blockSymbol
	unresolvedLabels map[string][]*LabelUse
}

func newSymbolResolver(symbolTable symbolTable) *symbolResolver {
	return &symbolResolver{
		symbolTable: symbolTable,
	}
}

func (this *symbolResolver) DoPass(ast Node) (Node, error) {
	return Traverse(this, ast)
}

func (this *symbolResolver) Visit(n Node) (Node, error) {
	switch m := n.(type) {
	case *Program:
		progInfo := &funcSymbol{
			level: 0,
			name:  "",
		}
		this.symbolTable[progInfo.name] = progInfo
		this.funcStack = this.funcStack.push(progInfo)
	case *Function:
		if _, exists := this.symbolTable[m.Name]; exists {
			return nil, errors.New(fmt.Sprintf("func '%s' exists", m.Name))
		} else {
			funcinfo := &funcSymbol{
				level: len(this.funcStack),
				name:  m.Name,
			}
			this.symbolTable[m.Name] = funcinfo
			this.funcStack = this.funcStack.push(funcinfo)
			this.unresolvedLabels = unresolvedLabels{}
			this.locals = map[string]*localSymbol{}
			this.blocks = map[string]*blockSymbol{}
			m.symbol = funcinfo
		}
	case *Signature:
		curfunc, _ := this.funcStack.peek()
		curfunc.rtype, _ = getDataType(m.Returns)
	case *Parameter:
		if _, exists := this.locals[m.Name]; !exists {
			dtype, _ := getDataType(m.TypeName)
			m.symbol = &localSymbol{
				parent: len(this.locals),
				index:  len(this.locals),
				dtype:  dtype,
				param:  true,
			}
			this.locals[m.Name] = m.symbol
			curfunc, _ := this.funcStack.peek()
			curfunc.locals = append(curfunc.locals, m.symbol)
		} else {
			return nil, errors.New(fmt.Sprintf("parameter '%s' exists", m.Name))
		}
	case *Local:
		if _, exists := this.locals[m.Name]; !exists {
			dtype, _ := getDataType(m.TypeName)
			m.symbol = &localSymbol{
				parent: len(this.locals),
				index:  len(this.locals),
				dtype:  dtype,
				param:  false,
			}
			this.locals[m.Name] = m.symbol
			curfunc, _ := this.funcStack.peek()
			curfunc.locals = append(curfunc.locals, m.symbol)
		} else {
			return nil, errors.New(fmt.Sprintf("local '%s' exists", m.Name))
		}
	case *Use:
		if localinfo, exists := this.locals[m.Name]; exists {
			m.symbol = localinfo
		} else {
			return nil, errors.New(fmt.Sprintf("local '%s' does not exist", m.Name))
		}
	case *Block:
		if _, exists := this.blocks[m.Name]; exists {
			return nil, errors.New(fmt.Sprintf("block '%s' exists", m.Name))
		} else {
			m.symbol = &blockSymbol{
				name:  m.Name,
				index: len(this.blocks),
			}

			this.blocks[m.Name] = m.symbol

			curfunc, _ := this.funcStack.peek()
			curfunc.blocks = append(curfunc.blocks, m.symbol)

			if unresolved, ok := this.unresolvedLabels[m.Name]; ok {
				for _, node := range unresolved {
					node.symbol = m.symbol
				}
				delete(this.unresolvedLabels, m.Name)
			}
		}
	case *LabelUse:
		if symbol, exists := this.blocks[m.Name]; exists {
			m.symbol = symbol
		} else {
			unresolved := this.unresolvedLabels[m.Name]
			this.unresolvedLabels[m.Name] = append(unresolved, m)
		}
	}
	return n, nil
}

func (this *symbolResolver) PostVisit(n Node) (Node, error) {
	switch n.(type) {
	case *Program:
		this.funcStack, _ = this.funcStack.pop()
	case *Function:
		if len(this.unresolvedLabels) > 0 {
			curfunc, _ := this.funcStack.peek()
			var labels []string
			for k, _ := range this.unresolvedLabels {
				labels = append(labels, k)
			}
			joined := strings.Join(labels, ", ")
			return nil, errors.New(fmt.Sprintf("func %s: unresolved labels: %s", curfunc.name, joined))
		} else {
			this.funcStack, _ = this.funcStack.pop()
		}
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
	case *Instruction:
		if use, ok := m.Dst.(*Use); !ok {
			return nil, errors.New("invalid ast!")
		} else if srcAType, err := getTypeOfNode(m.SrcA); err != nil {
			return nil, err
		} else if srcBType, err := getTypeOfNode(m.SrcB); err != nil {
			return nil, err
		} else if srcAType != srcBType {
			return nil, errors.New("source operand types do not match")
		} else if srcAType != use.symbol.dtype {
			return nil, errors.New("destination operand type does not match")
		} else {
			m.Opcode, _ = tokenTypeToOpcodeType_i64[m.OpcodeToken]
		}
	}

	return n, nil
}

func (this *typeChecker) PostVisit(n Node) (Node, error) {
	return n, nil
}
