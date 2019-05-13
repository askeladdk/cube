package il

import "testing"

func TestParseProgram(t *testing.T) {
	source := `
	func pow(b int64, e int64) int64 {
		var r int64
        entry:
			set r, 1
			goto loop
        loop:
            ifz e, done, body
        body:
            mul r, r, b
            sub e, e, 1
            goto loop
        done:
            ret r
    }`

	test := []string{
		"Program<>",
		"Unit<test.cubeasm>",
		"Function<pow>",
		"Signature<>",
		"Parameter<b>",
		"TypeName<int64>",
		"Parameter<e>",
		"TypeName<int64>",
		"TypeName<int64>",
		"Local<r>",
		"TypeName<int64>",
		"Block<entry>",
		"Set<>",
		"Use<r>",
		"Integer<1>",
		"Branch<>",
		"LabelUse<loop>",
		"Block<loop>",
		"ConditionalBranch<18>",
		"Use<e>",
		"LabelUse<done>",
		"LabelUse<body>",
		"Block<body>",
		"Instruction<21>",
		"Use<r>",
		"Use<r>",
		"Use<b>",
		"Instruction<20>",
		"Use<e>",
		"Use<e>",
		"Integer<1>",
		"Branch<>",
		"LabelUse<loop>",
		"Block<done>",
		"Return<>",
		"Use<r>",
	}

	if node, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if err := validateTrace(node, test); err != nil {
		t.Fatal(err)
	}
}
