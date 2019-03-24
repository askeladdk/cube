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

func (this *Lexer) matchKeyword(offset int, rest string, tokenType TokenType) TokenType {
	if this.position-this.initial == offset+len(rest) {
		s := this.source[this.initial+offset : this.position]
		if s == rest {
			return tokenType
		}
	}
	return IDENT
}

func (this *Lexer) peekat(offset int) byte {
	return this.source[this.initial+offset]
}

func (this *Lexer) identifierType() TokenType {
	switch this.peekat(0) {
	case 'a':
		return this.matchKeyword(1, "dd", ADD)
	case 'i':
		switch this.peekat(1) {
		case '3':
			if this.peekat(2) == '2' {
				return I32
			}
		case '6':
			if this.peekat(2) == '4' {
				return I64
			}
		case 'f':
			if this.peekat(2) == 'z' {
				return IFZ
			}
		}
	case 'f':
		return this.matchKeyword(1, "unc", FUNC)
	case 'g':
		return this.matchKeyword(1, "oto", GOTO)
	case 'm':
		return this.matchKeyword(1, "ul", MUL)
	case 'r':
		return this.matchKeyword(1, "et", RET)
	case 's':
		switch this.peekat(1) {
		case 'u':
			if this.peekat(2) == 'b' {
				return SUB
			}
		case 'e':
			if this.peekat(2) == 't' {
				return SET
			}
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
