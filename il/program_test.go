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
		"Def<r>",
		"Integer<1>",
		"Instruction<GOTO>",
		"LabelUse<loop>",
		"Block<loop>",
		"Instruction<IFZ>",
		"Use<e>",
		"LabelUse<done>",
		"LabelUse<body>",
		"Block<body>",
		"Instruction<MUL>",
		"Def<r>",
		"Use<r>",
		"Use<b>",
		"Instruction<SUB>",
		"Def<e>",
		"Use<e>",
		"Integer<1>",
		"Instruction<GOTO>",
		"LabelUse<loop>",
		"Block<done>",
		"Instruction<RET>",
		"Use<r>",
	}

	if node, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if err := validateTrace(node, test); err != nil {
		t.Fatal(err)
	}
}
