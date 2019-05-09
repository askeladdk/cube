package il

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

	I64
	IDENT
	FUNC
	SET
	RET
	GOTO
	IFZ
	ADD
	SUB
	MUL
)

type Token struct {
	Type   TokenType
	LineNo int
	Value  string
}
