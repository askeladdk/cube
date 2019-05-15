package il

import (
	"errors"
	"fmt"

	"github.com/askeladdk/cube"
)

type opcode struct {
	name     string
	src0Type int
	src1Type int
}

type instruction struct {
	opcode int
	dst    int
	src0   int
	src1   int
}

type block struct {
	name  string
	index int
	insrs []instruction
}

type local struct {
	name   string
	dtype  *cube.Type
	index  int
	parent int
	param  bool
}

type function struct {
	name   string
	rtype  *cube.Type
	locals []local
	blocks []block
	index  int
}

type program struct {
	funcs []function
}

type parseContext struct {
	lexer   *Lexer
	program *program

	peek Token

	localdefs  map[string]int
	funcdefs   map[string]int
	activeFunc function
}

func (this *parseContext) error(errmsg string) error {
	return errors.New(fmt.Sprintf("<%s>:%d: %s", this.lexer.Filename(), this.peek.LineNo, errmsg))
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

func (this *parseContext) typename() (*cube.Type, error) {
	switch this.peek.Type {
	case INT64:
		return cube.TypeInt64, this.advance()
	default:
		return nil, this.unexpected()
	}
}

func (this *parseContext) parameter() (string, *cube.Type, error) {
	if name, err := this.expect(IDENT); err != nil {
		return "", nil, err
	} else if dtype, err := this.typename(); err != nil {
		return "", nil, err
	} else {
		return name.Value, dtype, this.advance()
	}
}

func (this *parseContext) registerLocal(name string, dtype *cube.Type, param bool) error {
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

func (this *parseContext) blocks() error {
	for {
		if this.peek.Type == CURLY_R {
			return nil
		} else if _, err := this.expect(IDENT); err != nil {
			return err
		} else if _, err := this.expect(COLON); err != nil {
			return err
		} else {
			return nil
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
