package il

import "github.com/askeladdk/cube"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	NEWLINE

	INTEGER
	PAREN_L
	PAREN_R
	CURLY_L
	CURLY_R
	COMMA
	COLON
	MINUS
	ASSIGN

	I32
	I64
	IDENT
	FUNC
	RET
	ADD
	MUL
)

type Token struct {
	Type   TokenType
	LineNo int
	Value  string
}

var tokenToOpcodeType = map[TokenType]cube.OpcodeType{
	ADD: cube.ADD,
	MUL: cube.MUL,
	RET: cube.RET,
}
