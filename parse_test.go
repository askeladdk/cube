package cube

import "testing"

func TestParse_1(t *testing.T) {
	source := `
	func question() u64 {
		answer:
			reti 42
	}`

	lexer := NewLexer("test.cubeasm", source)
	ctx := parseContext{
		lexer:    lexer,
		program:  &program{},
		funcdefs: map[string]int{},
	}

	if err := ctx.parse(); err != nil {
		t.Fatal(err)
	} else if ctx.program.funcs[0].blocks[0].name != "answer" {
		t.Fatalf("wrong block name")
	} else if ctx.program.funcs[0].blocks[0].insrs[0].opcode != Opcode_RETI {
		t.Fatalf("wrong opcode")
	} else if ctx.program.funcs[0].blocks[0].insrs[0].operands[0] != 42 {
		t.Fatalf("wrong operand")
	}
}

func TestParse_2(t *testing.T) {
	source := `
	func a(b u64, c u64) u64 {
		var d u64
		var e u64
		entry:
			add d, b, c
			addi e, d, 0xffee
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

func TestParse_3(t *testing.T) {
	source := `
	func hang(a u64) u64 {
		loop1:
			jmp loop2
		loop2:
			jmp loop3
		loop3:
			jmp loop4
		loop4:
			jnz a, loop1, loop4
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
