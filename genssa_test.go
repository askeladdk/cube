package cube

import "testing"

func TestSSA_1(t *testing.T) {
	source := `
	func foo(a u64, b u64) u64 {
	var c u64
	entry:
		add a, a, a
		jmp label2
	label2:
		add b, b, b
		jmp label3
	label3:
		add c, a, b
		ret c
	}`

	err := Compile(&Config{
		Filename: "test.cubeasm",
		Source:   source,
		Procedure: func(proc *Procedure) error {
			proc = Pass_BuildCFG(proc)
			proc, _ = reallycrudessa(proc)
			return nil
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}
