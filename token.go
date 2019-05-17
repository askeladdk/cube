package cube

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

	U64
	IDENT
	FUNC
	SET
	RET
	REI
	JMP
	JNZ
	ADD
	ADI
	SUB
	MUL
	VAR
)

type Token struct {
	Type   TokenType
	LineNo int
	Value  string
}
