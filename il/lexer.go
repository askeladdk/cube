package il

import (
	"strings"
	"unicode"
)

func isbindigit(ch rune) bool {
	return '0' <= ch && ch <= '1'
}

func isdecdigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func ishexdigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('A' <= ch && ch <= 'F') || ('a' <= ch && ch <= 'f')
}

const eof rune = 0

type Lexer struct {
	reader   *strings.Reader
	filename string
	source   string
	peek     rune
	peekw    int
	lineno   int
	initial  int
	position int
}

func NewLexer(filename, source string) *Lexer {
	reader := strings.NewReader(source)
	peek, peekw, _ := reader.ReadRune()
	return &Lexer{
		reader:   reader,
		filename: filename,
		source:   source,
		peek:     peek,
		peekw:    peekw,
		lineno:   1,
	}
}

func (this *Lexer) Filename() string {
	return this.filename
}

func (this *Lexer) token(tokenType TokenType) Token {
	return Token{
		Type:   tokenType,
		LineNo: this.lineno,
		Value:  this.source[this.initial:this.position],
	}
}

func (this *Lexer) advance() {
	this.position += this.peekw
	this.peek, this.peekw, _ = this.reader.ReadRune()
}

func (this *Lexer) match(expected rune) bool {
	if this.peek != expected {
		return false
	} else {
		this.advance()
		return true
	}
}

func (this *Lexer) binnumber() Token {
	for isbindigit(this.peek) {
		this.advance()
	}
	return this.token(INTEGER)
}

func (this *Lexer) decnumber() Token {
	for isdecdigit(this.peek) {
		this.advance()
	}
	return this.token(INTEGER)
}

func (this *Lexer) hexnumber() Token {
	for ishexdigit(this.peek) {
		this.advance()
	}
	return this.token(INTEGER)
}

func (this *Lexer) comment() {
	for this.peek != eof && this.peek != '\n' {
		this.advance()
	}
}

func (this *Lexer) whitespace() {
	for {
		switch this.peek {
		case '\n':
			this.lineno += 1
			fallthrough
		case '\v':
			fallthrough
		case '\f':
			fallthrough
		case '\r':
			fallthrough
		case '\t':
			fallthrough
		case ' ':
			fallthrough
		case 0x85: // U+0085 (NEL)
			fallthrough
		case 0xA0: // U+00A0 (NBSP)
			this.advance()
			continue
		case ';':
			this.comment()
		default:
			return
		}
	}
}

var keywords = []struct {
	ident     string
	tokenType TokenType
}{
	// must be in alphabetical order
	{"add", ADD},
	{"func", FUNC},
	{"goto", GOTO},
	{"ifz", IFZ},
	{"int64", INT64},
	{"mul", MUL},
	{"ret", RET},
	{"set", SET},
	{"sub", SUB},
	{"var", VAR},
}

func (this *Lexer) identifierType() TokenType {
	test := this.source[this.initial:this.position]
	j := 0
outer:
	for _, keyword := range keywords {
		if len(test) == len(keyword.ident) {
			for ; j < len(test); j++ {
				if test[j] != keyword.ident[j] {
					continue outer
				}
			}
			return keyword.tokenType
		}
	}
	return IDENT
}

func (this *Lexer) identifier() Token {
	for this.peek == '_' || isdecdigit(this.peek) || unicode.IsLetter(this.peek) {
		this.advance()
	}
	return this.token(this.identifierType())
}

func (this *Lexer) Scan() Token {
	this.whitespace()

	this.initial = this.position

	ch := this.peek
	this.advance()

	switch ch {
	case eof:
		return this.token(EOF)
	case '(':
		return this.token(PAREN_L)
	case ')':
		return this.token(PAREN_R)
	case '{':
		return this.token(CURLY_L)
	case '}':
		return this.token(CURLY_R)
	case ',':
		return this.token(COMMA)
	case ':':
		if this.match('=') {
			return this.token(ASSIGN)
		} else {
			return this.token(COLON)
		}
	case '-':
		if !isdecdigit(this.peek) {
			return this.token(MINUS)
		}
		this.advance()
		fallthrough
	case '0':
		if this.match('b') {
			return this.binnumber()
		} else if this.match('x') {
			return this.hexnumber()
		}
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		fallthrough
	case '4':
		fallthrough
	case '5':
		fallthrough
	case '6':
		fallthrough
	case '7':
		fallthrough
	case '8':
		fallthrough
	case '9':
		return this.decnumber()
	case '_':
		return this.identifier()
	}

	if unicode.IsLetter(ch) {
		return this.identifier()
	} else {
		return this.token(ILLEGAL)
	}
}
