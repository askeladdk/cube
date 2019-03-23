package il

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/askeladdk/cube"
	"github.com/askeladdk/cube/ast"
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

func (this *Parser) integer() (*ast.Integer, error) {
	s := this.peek.Value

	base := 0
	if len(s) > 1 && s[0] == '0' && s[1] == 'b' {
		base = 2
	}

	if num, err := strconv.ParseInt(s, base, 64); err != nil {
		return nil, this.error(err.Error())
	} else {
		return &ast.Integer{Value: num}, this.advance()
	}
}

func (this *Parser) identifier() (*ast.Identifier, error) {
	name := this.peek.Value
	return &ast.Identifier{Name: name}, this.advance()
}

func (this *Parser) atom() (ast.Node, error) {
	switch this.peek.Type {
	case IDENT:
		return this.identifier()
	case INTEGER:
		return this.integer()
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) ret() (*ast.Return, error) {
	if node, err := this.atom(); err != nil {
		return nil, err
	} else {
		return &ast.Return{Src: node}, nil
	}
}

func (this *Parser) threeaddress(opcode cube.OpcodeType) (*ast.ThreeAddressInstruction, error) {
	if dst, err := this.identifier(); err != nil {
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
		return &ast.ThreeAddressInstruction{opcode, dst, srca, srcb, next}, nil
	}
}

func (this *Parser) instructions() (ast.Node, error) {
	tokenType := this.peek.Type
	switch tokenType {
	case ADD:
		fallthrough
	case MUL:
		if err := this.advance(); err != nil {
			return nil, err
		} else {
			opcode, _ := tokenToOpcodeType[tokenType]
			return this.threeaddress(opcode)
		}
	case RET:
		if err := this.advance(); err != nil {
			return nil, err
		} else {
			return this.ret()
		}
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) typename() (*ast.TypeName, error) {
	switch this.peek.Type {
	case I32:
		return &ast.TypeName{&cube.TypeInt32}, this.advance()
	case I64:
		return &ast.TypeName{&cube.TypeInt64}, this.advance()
	default:
		return nil, this.unexpected()
	}
}

func (this *Parser) parameter() (*ast.Parameter, error) {
	if name, err := this.expect(IDENT); err != nil {
		return nil, err
	} else if typename, err := this.typename(); err != nil {
		return nil, err
	} else if ok, err := this.accept(COMMA); err != nil {
		return nil, err
	} else if ok {
		if next, err := this.parameter(); err != nil {
			return nil, err
		} else {
			return &ast.Parameter{Name: name.Value, TypeName: typename, Next: next}, nil
		}
	} else {
		return &ast.Parameter{Name: name.Value, TypeName: typename}, nil
	}
}

func (this *Parser) parameters() (*ast.Parameter, error) {
	if _, err := this.expect(PAREN_L); err != nil {
		return nil, err
	} else if ok, err := this.accept(PAREN_R); err != nil {
		return nil, err
	} else if ok {
		return nil, nil
	} else if parameter, err := this.parameter(); err != nil {
		return nil, err
	} else if _, err := this.expect(PAREN_R); err != nil {
		return nil, err
	} else {
		return parameter, nil
	}
}

func (this *Parser) block(name string) (*ast.Block, error) {
	if _, err := this.expect(COLON); err != nil {
		return nil, err
	} else if instructions, err := this.instructions(); err != nil {
		return nil, err
	} else if next, err := this.blocks(); err != nil {
		return nil, err
	} else {
		return &ast.Block{name, instructions, next}, nil
	}
}

func (this *Parser) blocks() (*ast.Block, error) {
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

func (this *Parser) funcbody() (*ast.Block, error) {
	if _, err := this.expect(CURLY_L); err != nil {
		return nil, err
	} else if block, err := this.blocks(); err != nil {
		return nil, err
	} else {
		return block, nil
	}
}

func (this *Parser) function() (*ast.Function, error) {
	if name, err := this.expect(IDENT); err != nil {
		return nil, err
	} else if params, err := this.parameters(); err != nil {
		return nil, err
	} else if returns, err := this.typename(); err != nil {
		return nil, err
	} else if blocks, err := this.funcbody(); err != nil {
		return nil, err
	} else if next, err := this.definitions(); err != nil {
		return nil, err
	} else {
		return &ast.Function{
			Name:       name.Value,
			Parameters: params,
			Returns:    returns,
			Blocks:     blocks,
			Next:       next,
		}, nil
	}
}

func (this *Parser) definitions() (ast.Node, error) {
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

func (this *Parser) unit(next *ast.Unit) (*ast.Unit, error) {
	if definitions, err := this.definitions(); err != nil {
		return nil, err
	} else {
		return &ast.Unit{
			Filename:    this.lexer.Filename(),
			Definitions: definitions,
			Next:        next,
		}, nil
	}
}

func (this *Parser) Parse(next *ast.Unit) (*ast.Unit, error) {
	if err := this.advance(); err != nil {
		return nil, err
	} else if unit, err := this.unit(next); err != nil {
		return nil, err
	} else {
		return unit, nil
	}
}
