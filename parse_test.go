package cube

import "testing"

func TestParse_1(t *testing.T) {
	source := `
	func plusone(a u64) u64 {
		var b u64
		entry:
			add b, a, 1
			ret b
	}`

	err := Compile(&Config{
		Filename: "test.cubeasm",
		Source:   source,
		Procedure: func(proc *Procedure) error {
			if proc.blocks[0].name != "entry" {
				t.Fatalf("wrong block name")
			} else if proc.blocks[0].jmpretarg != 1 {
				t.Fatalf("wrong operand")
			}
			return nil
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParse_2(t *testing.T) {
	source := `
	func cfg(z u64) u64 {
		x: jnz z, b, c
		b: jmp d
		d: jmp g
		g: jmp d
		c: jmp e
		e: jmp m
		m: jmp c
	}`

	err := Compile(&Config{
		Filename: "test.cubeasm",
		Source:   source,
		Procedure: func(proc *Procedure) error {
			if len(proc.blocks) != 7 {
				t.Fatalf("wrong nr of blocks")
			}
			return nil
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}
