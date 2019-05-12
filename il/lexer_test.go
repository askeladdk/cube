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
		i64
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
		I64,
		IFZ,
		SET,
		SUB,
		ADD,
		MUL,
		VAR,
		EOF,
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
