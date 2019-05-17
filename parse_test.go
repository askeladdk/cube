package cube

import "testing"

func TestParse2(t *testing.T) {
	source := `
	func a(b u64, c u64) u64 {
		var d u64
		var e u64
		entry:
			add d, b, c
			adi e, d, 0xffee
			ret e
	}`

	lexer := NewLexer("test.cubeasm", source)
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
