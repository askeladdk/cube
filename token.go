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

	U64
	IDENT
	FUNC
	SET
	RET
	RETI
	JMP
	JNZ
	ADD
	ADDI
	SUB
	MUL
	VAR
)

type Token struct {
	Type   TokenType
	LineNo int
	Value  string
}
