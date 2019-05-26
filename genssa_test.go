package cube

import (
	"strings"
	"testing"
)

func TestSSA_1(t *testing.T) {
	source := `
	func pow(b u64, e u64) u64 {
	var r u64
	entry:
		mov r, 1
		jmp loop
	loop:
		jnz e, done, body
	body:
		mul r, r, b
		sub e, e, 1
		jmp loop
	done:
		ret r
	}`

	err := Compile(&Config{
		Filename: "test.cubeasm",
		Source:   source,
		Procedure: func(proc *Procedure) error {
			proc = Pass_BuildCFG(proc)
			proc, _ = reallycrudessa(proc)
			var sb strings.Builder
			printproc(&sb, proc)
			s := sb.String()
			_ = s
			return nil
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}
