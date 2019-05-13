package il

import "testing"

func TestScanNumbers(t *testing.T) {
	lexer := NewLexer("<source>", `
		0xa0f9
		0b10011
		42
		-42
	`)
	numbers := []string{
		"0xa0f9",
		"0b10011",
		"42",
		"-42",
	}

	for _, expected := range numbers {
		if token := lexer.Scan(); token.Type != INTEGER {
			t.Fatal(token)
		} else if token.Value != expected {
			t.Fatal(token)
		}
	}
}

func TestScanKeywords(t *testing.T) {
	lexer := NewLexer("<source>", `
		goto
		func
		int64
		ifz
		set
		sub
		add
		mul
		var
	`)

	tokens := []TokenType{
		GOTO,
		FUNC,
		INT64,
		IFZ,
		SET,
		SUB,
		ADD,
		MUL,
		VAR,
	}

	for _, expected := range tokens {
		if token := lexer.Scan(); token.Type != expected {
			t.Fatal(token)
		}
	}
}

func TestScanIdentifiers(t *testing.T) {
	lexer := NewLexer("<source>", `
		a
		sett
		hëlló
		你好
	`)

	identifiers := []string{
		"a",
		"sett",
		"hëlló",
		"你好",
	}

	for _, expected := range identifiers {
		if token := lexer.Scan(); token.Type != IDENT {
			t.Fatal(token)
		} else if token.Value != expected {
			t.Fatal(token)
		}
	}
}

func TestScanIllegal(t *testing.T) {
	lexer := NewLexer("<source>", "%")
	token := lexer.Scan()
	if token.Type != ILLEGAL {
		t.Fatalf("should be illegal")
	}
}

func TestScanMisc(t *testing.T) {
	lexer := NewLexer("<source>", `
	; comment ⌘
	() {} :, - := _ ?
	`)

	tokens := []TokenType{
		PAREN_L,
		PAREN_R,
		CURLY_L,
		CURLY_R,
		COLON,
		COMMA,
		MINUS,
		ASSIGN,
		IDENT,
		ILLEGAL,
	}

	for _, expected := range tokens {
		if token := lexer.Scan(); token.Type != expected {
			t.Fatal(token)
		}
	}
}
