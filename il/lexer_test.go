package il

import "testing"

func TestScanHexNumber(t *testing.T) {
	lexer := NewLexer("<source>", "0xa0f9 ")
	token := lexer.Scan()
	if token.Type != INTEGER || token.Value != "0xa0f9" {
		t.Fatalf("wrong token %d %s", token.Type, token.Value)
	}
}

func TestScanBinNumber(t *testing.T) {
	lexer := NewLexer("<source>", "0b10011 ")
	token := lexer.Scan()
	if token.Type != INTEGER || token.Value != "0b10011" {
		t.Fatalf("wrong token %d %s", token.Type, token.Value)
	}
}

func TestScanInteger(t *testing.T) {
	lexer := NewLexer("<source>", "42 ")
	token := lexer.Scan()
	if token.Type != INTEGER || token.Value != "42" {
		t.Fatalf("wrong token %d %s", token.Type, token.Value)
	}
}

func TestScanNegativeInteger(t *testing.T) {
	lexer := NewLexer("<source>", "-42 ")
	token := lexer.Scan()
	if token.Type != INTEGER || token.Value != "-42" {
		t.Fatalf("wrong token %d %s", token.Type, token.Value)
	}
}

func TestScanIdentifier(t *testing.T) {
	lexer := NewLexer("<source>", "hëllõ ")
	token := lexer.Scan()
	if token.Type != IDENT || token.Value != "hëllõ" {
		t.Fatalf("wrong token %d %s", token.Type, token.Value)
	}
}

func TestScanKeyword(t *testing.T) {
	lexer := NewLexer("<source>", "func ")
	token := lexer.Scan()
	if token.Type != FUNC || token.Value != "func" {
		t.Fatalf("wrong token %d %s", token.Type, token.Value)
	}
}

func TestScanKeyword2(t *testing.T) {
	lexer := NewLexer("<source>", "function ")
	token := lexer.Scan()
	if token.Type != IDENT || token.Value != "function" {
		t.Fatalf("wrong token %d %s", token.Type, token.Value)
	}
}

func TestScanComment(t *testing.T) {
	lexer := NewLexer("<source>", ";42\n;1337")
	token := lexer.Scan()
	if token.Type != EOF {
		t.Fatalf("should be eof")
	}
}

func TestScanWhitespace(t *testing.T) {
	lexer := NewLexer("<source>", "\n; this is a test\n\t1337\n\n")
	token := lexer.Scan()
	if token.Type != INTEGER || token.LineNo != 3 || token.Value != "1337" {
		t.Fatalf("wrong token :%d %d %s", token.LineNo, token.Type, token.Value)
	}
}

func TestScanIllegal(t *testing.T) {
	lexer := NewLexer("<source>", "%")
	token := lexer.Scan()
	if token.Type != ILLEGAL {
		t.Fatalf("should be illegal")
	}
}

func TestScanEverything(t *testing.T) {
	source := `
    ; comment ⌘
    0 --0x1337 0b110011
    func _你好(i32) {
		var
        hëlló(wÖrld i32):
            日本語 := add wÖrld, i64 42
            ret 日本語
    }`

	lexer := NewLexer("<source>", source)
	for token := lexer.Scan(); token.Type != EOF; token = lexer.Scan() {
		switch token.Type {
		case ILLEGAL:
			t.Fatalf("illegal character: %s", token.Value)
			break
		}
	}
}
