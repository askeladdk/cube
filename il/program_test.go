package il

import "testing"

func TestParseProgram(t *testing.T) {
	source := `
	func pow(b i64, e i64) i64 {
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
		"Block<entry>",
		"Instruction<SET>",
		"Identifier<r>",
		"Integer<1>",
		"Instruction<GOTO>",
		"Identifier<loop>",
		"Block<loop>",
		"Instruction<IFZ>",
		"Identifier<e>",
		"Identifier<done>",
		"Identifier<body>",
		"Block<body>",
		"Instruction<MUL>",
		"Identifier<r>",
		"Identifier<r>",
		"Identifier<b>",
		"Instruction<SUB>",
		"Identifier<e>",
		"Identifier<e>",
		"Integer<1>",
		"Instruction<GOTO>",
		"Identifier<loop>",
		"Block<done>",
		"Instruction<RET>",
		"Identifier<r>",
	}

	if node, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if err := validateTrace(node, test); err != nil {
		t.Fatal(err)
	}
}
