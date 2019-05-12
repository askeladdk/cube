package il

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/askeladdk/cube"
)

type Parser struct {
	lexer *Lexer
	peek  Token
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func (this *Parser) error(errmsg string) error {
	return errors.New(fmt.Sprintf("<%s>:%d: %s", this.lexer.Filename(), this.peek.LineNo, errmsg))
}

func (this *Parser) unexpected() error {
	if this.peek.Type == EOF {
		return this.error("unexpected end of file")
	} else {
		return this.error(fmt.Sprintf("unexpected symbol '%s'", this.peek.Value))
	}
}

func (this *Parser) advance() error {
	this.peek = this.lexer.Scan()

	if this.peek.Type == ILLEGAL {
		return this.error(fmt.Sprintf("illegal character '%s'", this.peek.Value))
	}

	return nil
}

func (this *Parser) check(tt TokenType) bool {
	return this.peek.Type == tt
}

func (this *Parser) accept(tt TokenType) (bool, error) {
	if this.check(tt) {
		return true, this.advance()
	} else {
		return false, nil
	}
}

func (this *Parser) expect(tt TokenType) (Token, error) {
	current := this.peek
	if this.check(tt) {
		return current, this.advance()
	} else {
		return current, this.unexpected()
	}
}

func (this *Parser) integer() (*Integer, error) {
	s := this.peek.Value

	base := 0
	if len(s) > 1 && s[0] == '0' && s[1] == 'b' {
		base = 2
	}

	if num, err := strconv.ParseInt(s, base, 64); err != nil {
		return nil, this.error(err.Error())
	} else {
		return &Integer{Value: num}, this.advance()
	}
}

func (this *Parser) use() (*Use, error) {
	name := this.peek.Value
	return &Use{Name: name}, this.advance()
}

func (this *Parser) labeluse() (*LabelUse, error) {
	name := this.peek.Value
	return &LabelUse{Name: name}, this.advance()
}

func (this *Parser) atom() (Node, error) {
	switch this.peek.Type {
	case IDENT:
		return this.use()
	case INTEGER:
		return this.integer()
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) ret() (*Return, error) {
	if src, err := this.atom(); err != nil {
		return nil, err
	} else {
		return &Return{
			Src: src,
		}, nil
	}
}

func (this *Parser) set() (*Set, error) {
	if dst, err := this.use(); err != nil {
		return nil, err
	} else if _, err := this.expect(COMMA); err != nil {
		return nil, err
	} else if src, err := this.atom(); err != nil {
		return nil, err
	} else if next, err := this.instructions(); err != nil {
		return nil, err
	} else {
		return &Set{
			Dst:  dst,
			Src:  src,
			Next: next,
		}, nil
	}
}

func (this *Parser) gotoinsr() (*Branch, error) {
	if label, err := this.labeluse(); err != nil {
		return nil, err
	} else {
		return &Branch{
			Label: label,
		}, nil
	}
}

func (this *Parser) conditional(opcodeToken TokenType) (*ConditionalBranch, error) {
	if opa, err := this.atom(); err != nil {
		return nil, err
	} else if _, err := this.expect(COMMA); err != nil {
		return nil, err
	} else if labela, err := this.labeluse(); err != nil {
		return nil, err
	} else if _, err := this.expect(COMMA); err != nil {
		return nil, this.unexpected()
	} else if labelb, err := this.labeluse(); err != nil {
		return nil, err
	} else {
		return &ConditionalBranch{
			OpcodeToken: opcodeToken,
			Cond:        opa,
			LabelA:      labela,
			LabelB:      labelb,
		}, nil
	}
}

func (this *Parser) instruction(opcodeToken TokenType) (*Instruction, error) {
	if dst, err := this.use(); err != nil {
		return nil, err
	} else if _, err := this.expect(COMMA); err != nil {
		return nil, err
	} else if srca, err := this.atom(); err != nil {
		return nil, err
	} else if _, err := this.expect(COMMA); err != nil {
		return nil, this.unexpected()
	} else if srcb, err := this.atom(); err != nil {
		return nil, err
	} else if next, err := this.instructions(); err != nil {
		return nil, err
	} else {
		return &Instruction{
			OpcodeToken: opcodeToken,
			Dst:         dst,
			SrcA:        srca,
			SrcB:        srcb,
			Next:        next,
		}, nil
	}
}

func (this *Parser) instructions() (Node, error) {
	tokenType := this.peek.Type
	possibleErr := this.unexpected()

	if err := this.advance(); err != nil {
		return nil, err
	} else {
		switch tokenType {
		case ADD:
			fallthrough
		case SUB:
			fallthrough
		case MUL:
			return this.instruction(tokenType)
		case RET:
			return this.ret()
		case SET:
			return this.set()
		case GOTO:
			return this.gotoinsr()
		case IFZ:
			return this.conditional(tokenType)
		default:
			return nil, possibleErr
		}
	}
}

func (this *Parser) typename() (*TypeName, error) {
	switch this.peek.Type {
	case I64:
		return &TypeName{cube.TypeInt64}, this.advance()
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) parameter() (*Parameter, error) {
	if name, err := this.expect(IDENT); err != nil {
		return nil, err
	} else if typename, err := this.typename(); err != nil {
		return nil, err
	} else if next, err := this.parameterNext(); err != nil {
		return nil, err
	} else {
		return &Parameter{
			Name:     name.Value,
			TypeName: typename,
			Next:     next,
		}, nil
	}
}

func (this *Parser) parameterNext() (*Parameter, error) {
	switch this.peek.Type {
	case COMMA:
		if err := this.advance(); err != nil {
			return nil, err
		} else {
			return this.parameter()
		}
	case PAREN_R:
		return nil, this.advance()
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) parameters() (*Parameter, error) {
	switch this.peek.Type {
	case IDENT:
		return this.parameter()
	case PAREN_R:
		return nil, this.advance()
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) signature() (*Signature, error) {
	if _, err := this.expect(PAREN_L); err != nil {
		return nil, err
	} else if parameters, err := this.parameters(); err != nil {
		return nil, err
	} else if rtype, err := this.typename(); err != nil {
		return nil, err
	} else {
		return &Signature{
			Parameters: parameters,
			Returns:    rtype,
		}, nil
	}
}

func (this *Parser) locals() (*Local, error) {
	if ok, err := this.accept(VAR); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	} else if name, err := this.expect(IDENT); err != nil {
		return nil, err
	} else if typename, err := this.typename(); err != nil {
		return nil, err
	} else if next, err := this.locals(); err != nil {
		return nil, err
	} else {
		return &Local{
			Name:     name.Value,
			TypeName: typename,
			Next:     next,
		}, nil
	}
}

func (this *Parser) block(name string) (*Block, error) {
	if _, err := this.expect(COLON); err != nil {
		return nil, err
	} else if instructions, err := this.instructions(); err != nil {
		return nil, err
	} else if next, err := this.blocks(); err != nil {
		return nil, err
	} else {
		return &Block{
			Name:         name,
			Instructions: instructions,
			Next:         next,
		}, nil
	}
}

func (this *Parser) blocks() (*Block, error) {
	token := this.peek
	switch token.Type {
	case IDENT:
		if err := this.advance(); err != nil {
			return nil, err
		} else {
			return this.block(token.Value)
		}
	case CURLY_R:
		return nil, this.advance()
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) function() (*Function, error) {
	if name, err := this.expect(IDENT); err != nil {
		return nil, err
	} else if signature, err := this.signature(); err != nil {
		return nil, err
	} else if _, err := this.expect(CURLY_L); err != nil {
		return nil, err
	} else if locals, err := this.locals(); err != nil {
		return nil, err
	} else if blocks, err := this.blocks(); err != nil {
		return nil, err
	} else if next, err := this.definitions(); err != nil {
		return nil, err
	} else {
		return &Function{
			Name:      name.Value,
			Signature: signature,
			Locals:    locals,
			Blocks:    blocks,
			Next:      next,
		}, nil
	}
}

func (this *Parser) definitions() (Node, error) {
	switch this.peek.Type {
	case FUNC:
		if err := this.advance(); err != nil {
			return nil, err
		} else {
			return this.function()
		}
	case EOF:
		return nil, nil
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) unit(next *Unit) (*Unit, error) {
	if definitions, err := this.definitions(); err != nil {
		return nil, err
	} else {
		return &Unit{
			Filename:    this.lexer.Filename(),
			Definitions: definitions,
			Next:        next,
		}, nil
	}
}

func (this *Parser) Parse(next *Unit) (*Unit, error) {
	if err := this.advance(); err != nil {
		return nil, err
	} else if unit, err := this.unit(next); err != nil {
		return nil, err
	} else {
		return unit, nil
	}
}
