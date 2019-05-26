package cube

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type unresolvedLabel struct {
	block   *BasicBlock
	succidx int
}

type parseContext struct {
	config *Config
	lexer  *Lexer
	peek   Token

	localdefs        map[string]int
	blockdefs        map[string]*BasicBlock
	curproc          Procedure
	curblock         *BasicBlock
	unresolvedLabels map[string][]unresolvedLabel
}

func (this *parseContext) registerLocal(name string, dtype *Type, param bool) error {
	if _, exists := this.localdefs[name]; exists {
		return this.error(fmt.Sprintf("local %s is redefined here", name))
	} else {
		index := len(this.localdefs)
		newlocal := Local{
			name:        name,
			dataType:    dtype,
			isParameter: true,
		}

		this.curproc.locals = append(this.curproc.locals, newlocal)
		this.localdefs[name] = index
		return nil
	}
}

func (this *parseContext) lookupLocal(name string) (int, error) {
	if local, ok := this.localdefs[name]; !ok {
		return 0, this.error(fmt.Sprintf("undefined local '%s' referenced here", name))
	} else {
		return local, nil
	}
}

func (this *parseContext) error(errmsg string) error {
	return errors.New(fmt.Sprintf("%s:%d: %s", this.config.Filename, this.peek.LineNo, errmsg))
}

func (this *parseContext) unexpected() error {
	if this.peek.Type == EOF {
		return this.error("unexpected end of file")
	} else {
		return this.error(fmt.Sprintf("unexpected symbol '%s'", this.peek.Value))
	}
}

func (this *parseContext) advance() error {
	this.peek = this.lexer.Scan()

	if this.peek.Type == ILLEGAL {
		return this.error(fmt.Sprintf("illegal character '%s'", this.peek.Value))
	}

	return nil
}

func (this *parseContext) match(tt TokenType) (bool, error) {
	if this.peek.Type == tt {
		return true, this.advance()
	} else {
		return false, nil
	}
}

func (this *parseContext) expect(tt TokenType) (Token, error) {
	current := this.peek
	if this.peek.Type == tt {
		return current, this.advance()
	} else {
		return current, this.unexpected()
	}
}

func (this *parseContext) ident() (string, error) {
	if token, err := this.expect(IDENT); err != nil {
		return "", err
	} else {
		return token.Value, nil
	}
}

func (this *parseContext) typename() (*Type, error) {
	switch this.peek.Type {
	case U64:
		return TypeUntyped64, this.advance()
	default:
		return nil, this.unexpected()
	}
}

func (this *parseContext) parameters() error {
	if matched, err := this.match(PAREN_R); err != nil {
		return err
	} else if matched {
		return nil
	} else {
		for {
			if name, err := this.ident(); err != nil {
				return err
			} else if dtype, err := this.typename(); err != nil {
				return err
			} else if err := this.registerLocal(name, dtype, true); err != nil {
				return err
			} else {
				switch this.peek.Type {
				case COMMA:
					if err := this.advance(); err != nil {
						return err
					} else {
						continue
					}
				case PAREN_R:
					return this.advance()
				default:
					return this.unexpected()
				}
			}
		}
	}
}

func (this *parseContext) vars() error {
	for {
		if matched, err := this.match(VAR); err != nil {
			return err
		} else if !matched {
			return nil
		} else {
			if name, err := this.ident(); err != nil {
				return err
			} else if dtype, err := this.typename(); err != nil {
				return err
			} else if err := this.registerLocal(name, dtype, false); err != nil {
				return err
			}
		}
	}
}

func parseInt(val string) (uint64, error) {
	base := 0
	if len(val) > 1 && val[0] == '0' && val[1] == 'b' {
		base = 2
	}

	if num, err := strconv.ParseInt(val, base, 64); err != nil {
		return 0, err
	} else {
		return uint64(num), nil
	}
}

func (this *parseContext) constant(num uint64) int {
	for i, c := range this.curproc.constants {
		if c == num {
			return i
		}
	}

	i := len(this.curproc.constants)
	this.curproc.constants = append(this.curproc.constants, num)
	return i
}

func (this *parseContext) local() (int, error) {
	if ident, err := this.ident(); err != nil {
		return 0, err
	} else if local, err := this.lookupLocal(ident); err != nil {
		return 0, err
	} else {
		return local, nil
	}
}

func (this *parseContext) atom() (operand, error) {
	switch this.peek.Type {
	case INTEGER:
		if num, err := parseInt(this.peek.Value); err != nil {
			return operandNil, this.error(err.Error())
		} else {
			return operandCon(this.constant(num)), this.advance()
		}
	case IDENT:
		if local, err := this.lookupLocal(this.peek.Value); err != nil {
			return operandNil, err
		} else {
			return operandLoc(local), this.advance()
		}
	default:
		return operandNil, this.unexpected()
	}
}

func (this *parseContext) label(succidx int) (*BasicBlock, error) {
	if name, err := this.ident(); err != nil {
		return nil, err
	} else if block, ok := this.blockdefs[name]; !ok {
		unresolved, _ := this.unresolvedLabels[name]
		this.unresolvedLabels[name] = append(unresolved, unresolvedLabel{
			block:   this.curblock,
			succidx: succidx,
		})
		return nil, nil
	} else {
		return block, nil
	}
}

func (this *parseContext) resolveLabel(name string, block *BasicBlock) {
	unresolved, _ := this.unresolvedLabels[name]
	for _, u := range unresolved {
		u.block.successors[u.succidx] = block
	}
	delete(this.unresolvedLabels, name)
}

func (this *parseContext) emit(opc opcode, op0, op1, op2 operand) error {
	this.curblock.instructions = append(this.curblock.instructions, Instruction{
		opcode:   opc,
		operands: [3]operand{op0, op1, op2},
	})
	return nil
}

func (this *parseContext) ret() error {
	if op1, err := this.atom(); err != nil {
		return err
	} else {
		this.curblock.jmpcode = opcode_RET
		this.curblock.jmpretarg = op1.value
		return this.emit(opcode_RET, operandNil, op1, operandNil)
	}
}

func (this *parseContext) jmp() error {
	if op0, err := this.label(0); err != nil {
		return err
	} else {
		this.curblock.jmpcode = opcode_JMP
		this.curblock.successors[0] = op0
		return this.emit(opcode_JMP, operandNil, operandNil, operandNil)
	}
}

func (this *parseContext) jnz() error {
	if op0, err := this.local(); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if op1, err := this.label(0); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if op2, err := this.label(1); err != nil {
		return err
	} else if op1 == op2 {
		this.curblock.jmpcode = opcode_JMP
		this.curblock.successors[0] = op1
		return this.emit(opcode_JMP, operandNil, operandNil, operandNil)
	} else {
		this.curblock.jmpcode = opcode_JNZ
		this.curblock.jmpretarg = op0
		this.curblock.successors[0] = op1
		this.curblock.successors[1] = op2
		return this.emit(opcode_JNZ, operandNil, operandLoc(op0), operandNil)
	}
}

func (this *parseContext) instruction_raa(opc opcode) error {
	if dstloc, err := this.local(); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if op1, err := this.atom(); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if op2, err := this.atom(); err != nil {
		return err
	} else {
		return this.emit(opc, operandLoc(dstloc), op1, op2)
	}
}

func (this *parseContext) instruction_ra(opc opcode) error {
	if dstloc, err := this.local(); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if op1, err := this.atom(); err != nil {
		return err
	} else {
		return this.emit(opc, operandLoc(dstloc), op1, operandNil)
	}
}

func (this *parseContext) instructions() error {
	for {
		tokenType := this.peek.Type
		if err := this.advance(); err != nil {
			return err
		} else {
			var err error
			switch tokenType {
			case ADD:
				err = this.instruction_raa(opcode_ADD)
			case SUB:
				err = this.instruction_raa(opcode_SUB)
			case MOV:
				err = this.instruction_ra(opcode_MOV)
			case RET:
				return this.ret()
			case JMP:
				return this.jmp()
			case JNZ:
				return this.jnz()
			default:
				return this.unexpected()
			}

			if err != nil {
				return err
			}
		}
	}
}

func (this *parseContext) blocks() error {
	this.blockdefs = map[string]*BasicBlock{}
	this.unresolvedLabels = map[string][]unresolvedLabel{}

	for do := true; do; do = this.peek.Type != CURLY_R {
		if name, err := this.ident(); err != nil {
			return err
		} else if _, err := this.expect(COLON); err != nil {
			return err
		} else if _, exists := this.blockdefs[name]; exists {
			return this.error(fmt.Sprintf("block %s redefined here", name))
		} else {
			this.curblock = &BasicBlock{
				name: name,
			}
			this.blockdefs[name] = this.curblock

			this.resolveLabel(name, this.curblock)

			if err := this.instructions(); err != nil {
				return err
			} else {
				this.curproc.blocks = append(this.curproc.blocks, this.curblock)
			}
		}
	}

	return nil
}

func (this *parseContext) procedure() error {
	this.curproc = Procedure{}
	this.localdefs = map[string]int{}

	if name, err := this.ident(); err != nil {
		return err
	} else if _, err := this.expect(PAREN_L); err != nil {
		return err
	} else if err := this.parameters(); err != nil {
		return err
	} else if rtype, err := this.typename(); err != nil {
		return err
	} else if _, err := this.expect(CURLY_L); err != nil {
		return err
	} else if err := this.vars(); err != nil {
		return err
	} else if err := this.blocks(); err != nil {
		return err
	} else if _, err := this.expect(CURLY_R); err != nil {
		return err
	} else if len(this.unresolvedLabels) > 0 {
		var labels []string
		for k, _ := range this.unresolvedLabels {
			labels = append(labels, k)
		}
		if len(labels) > 1 {
			joinedLabels := strings.Join(labels, ", ")
			return this.error(fmt.Sprintf("unresolved references to labels %s", joinedLabels))
		} else {
			return this.error(fmt.Sprintf("unresolved reference to label %s", labels[0]))
		}
	} else {
		this.curproc.name = name
		this.curproc.returnType = rtype
		this.curproc.entryPoint = this.curproc.blocks[0]
		return this.config.Procedure(&this.curproc)
	}
}

func (this *parseContext) definitions() error {
	for {
		switch this.peek.Type {
		case FUNC:
			if err := this.advance(); err != nil {
				return err
			} else if err := this.procedure(); err != nil {
				return err
			}
		case EOF:
			return nil
		default:
			return this.unexpected()
		}
	}
}

func (this *parseContext) parse() error {
	if err := this.advance(); err != nil {
		return err
	} else {
		return this.definitions()
	}
}
