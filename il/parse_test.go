package il

import "testing"

func TestParse2(t *testing.T) {
	source := `
	func a(b int64, c int64) int64 {
		var d int64
		var e int64
	}
	`
	lexer := NewLexer("<source>", source)
	ctx := parseContext{
		lexer:    lexer,
		program:  &program{},
		funcdefs: map[string]int{},
	}

	if err := ctx.parse(); err != nil {
		t.Fatal(err)
	}
	t.Fail()
}
