package cube

import "testing"

func TestScanNumbers(t *testing.T) {
	lexer := NewLexer(`
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
	lexer := NewLexer(`
		jmp
		func
		u64
		jnz
		set
		sub
		add
		mul
		var
		ret
		reti
		addi
	`)

	tokens := []TokenType{
		JMP,
		FUNC,
		U64,
		JNZ,
		SET,
		SUB,
		ADD,
		MUL,
		VAR,
		RET,
		RETI,
		ADDI,
	}

	for _, expected := range tokens {
		if token := lexer.Scan(); token.Type != expected {
			t.Fatal(token)
		}
	}
}

func TestScanIdentifiers(t *testing.T) {
	lexer := NewLexer(`
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
	lexer := NewLexer("%")
	token := lexer.Scan()
	if token.Type != ILLEGAL {
		t.Fatalf("should be illegal")
	}
}

func TestScanMisc(t *testing.T) {
	lexer := NewLexer(`
	; comment ⌘
	() {} :, _ ?
	`)

	tokens := []TokenType{
		PAREN_L,
		PAREN_R,
		CURLY_L,
		CURLY_R,
		COLON,
		COMMA,
		IDENT,
		ILLEGAL,
	}

	for _, expected := range tokens {
		if token := lexer.Scan(); token.Type != expected {
			t.Fatal(token)
		}
	}
}
