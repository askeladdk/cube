package cube

import (
	"errors"
	"fmt"
	"strconv"
)

type parseContext struct {
	lexer   *Lexer
	program *program

	peek Token

	localdefs   map[string]int
	blockdefs   map[string]int
	funcdefs    map[string]int
	activeFunc  function
	activeBlock block
}

func (this *parseContext) registerLocal(name string, dtype *Type, param bool) error {
	if _, exists := this.localdefs[name]; exists {
		return this.error(fmt.Sprintf("local %s is redefined here", name))
	} else {
		index := len(this.localdefs)
		newlocal := local{
			name:   name,
			index:  index,
			parent: index,
			dtype:  dtype,
			param:  true,
		}

		this.activeFunc.locals = append(this.activeFunc.locals, newlocal)
		this.localdefs[name] = index
		return nil
	}
}

func (this *parseContext) lookupLocal(name string) (operand, error) {
	if local, ok := this.localdefs[name]; !ok {
		return 0, this.error(fmt.Sprintf("undefined local '%s' referenced here", name))
	} else {
		return operand(local), nil
	}
}

func (this *parseContext) error(errmsg string) error {
	return errors.New(fmt.Sprintf("%s:%d: %s", this.lexer.Filename(), this.peek.LineNo, errmsg))
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

func (this *parseContext) immediate() (operand, error) {
	switch this.peek.Type {
	case INTEGER:
		val := this.peek.Value

		base := 0
		if len(val) > 1 && val[0] == '0' && val[1] == 'b' {
			base = 2
		}

		if num, err := strconv.ParseInt(val, base, 64); err != nil {
			return 0, this.error(err.Error())
		} else {
			return operand(num), this.advance()
		}
	default:
		return 0, this.unexpected()
	}
}

func (this *parseContext) instruction_r(opcode *OpcodeType) error {
	if r0, err := this.ident(); err != nil {
		return err
	} else if r0idx, err := this.lookupLocal(r0); err != nil {
		return err
	} else {
		insr := instruction{
			opcode:   opcode,
			operands: [3]operand{r0idx},
		}
		this.activeBlock.insrs = append(this.activeBlock.insrs, insr)
		return nil
	}
}

func (this *parseContext) instruction_rrr(opcode *OpcodeType) error {
	if ident0, err := this.ident(); err != nil {
		return err
	} else if op0, err := this.lookupLocal(ident0); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if ident1, err := this.ident(); err != nil {
		return err
	} else if op1, err := this.lookupLocal(ident1); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if ident2, err := this.ident(); err != nil {
		return err
	} else if op2, err := this.lookupLocal(ident2); err != nil {
		return err
	} else {
		insr := instruction{
			opcode:   opcode,
			operands: [3]operand{op0, op1, op2},
		}
		this.activeBlock.insrs = append(this.activeBlock.insrs, insr)
		return nil
	}
}

func (this *parseContext) instruction_rri(opcode *OpcodeType) error {
	if ident0, err := this.ident(); err != nil {
		return err
	} else if op0, err := this.lookupLocal(ident0); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if ident1, err := this.ident(); err != nil {
		return err
	} else if op1, err := this.lookupLocal(ident1); err != nil {
		return err
	} else if _, err := this.expect(COMMA); err != nil {
		return err
	} else if op2, err := this.immediate(); err != nil {
		return err
	} else {
		insr := instruction{
			opcode:   opcode,
			operands: [3]operand{op0, op1, op2},
		}
		this.activeBlock.insrs = append(this.activeBlock.insrs, insr)
		return nil
	}
}

func (this *parseContext) instructions() error {
	for {
		tokenType := this.peek.Type
		if err := this.advance(); err != nil {
			return err
		} else {
			switch tokenType {
			case ADD:
				if err := this.instruction_rrr(Opcode_ADD); err != nil {
					return err
				}
			case ADI:
				if err := this.instruction_rri(Opcode_ADI); err != nil {
					return err
				}
			case RET:
				if err := this.instruction_r(Opcode_RET); err != nil {
					return err
				} else {
					return nil
				}
			default:
				return this.unexpected()
			}
		}
	}
}

func (this *parseContext) blocks() error {
	this.blockdefs = map[string]int{}

	for {
		if this.peek.Type == CURLY_R {
			return nil
		} else if name, err := this.ident(); err != nil {
			return err
		} else if _, err := this.expect(COLON); err != nil {
			return err
		} else if _, exists := this.blockdefs[name]; exists {
			return this.error(fmt.Sprintf("block %s redefined here", name))
		} else {
			this.activeBlock = block{
				name:  name,
				index: len(this.activeFunc.blocks),
			}
			this.blockdefs[name] = this.activeBlock.index
			if err := this.instructions(); err != nil {
				return err
			} else {
				this.activeFunc.blocks = append(this.activeFunc.blocks, this.activeBlock)
				return nil
			}
		}
	}
}

func (this *parseContext) function() error {
	this.activeFunc = function{}
	this.localdefs = map[string]int{}

	if name, err := this.ident(); err != nil {
		return err
	} else if _, exists := this.funcdefs[name]; exists {
		return this.error(fmt.Sprintf("function %s redefined here", name))
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
	} else {
		index := len(this.program.funcs)
		this.activeFunc.index = index
		this.activeFunc.name = name
		this.activeFunc.rtype = rtype
		this.funcdefs[name] = index
		this.program.funcs = append(this.program.funcs, this.activeFunc)
		return nil
	}
}

func (this *parseContext) definitions() error {
	for {
		switch this.peek.Type {
		case FUNC:
			if err := this.advance(); err != nil {
				return err
			} else if err := this.function(); err != nil {
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

func Parse2(lexer *Lexer, program *program) error {
	return (&parseContext{
		lexer:    lexer,
		program:  program,
		funcdefs: map[string]int{},
	}).parse()
}
