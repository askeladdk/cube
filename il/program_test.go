package il

import "testing"

func TestParseProgram(t *testing.T) {
	source := `
    func fma(a i64, b i64, c i64) i64 {
        entry:
            mul d, a, b
            add e, d, c
            ret e
    }`

	test := []string{
		"Program<>",
		"Unit<test.cubeasm>",
		"Function<fma>",
		"Parameter<a>",
		"TypeName<int64>",
		"Parameter<b>",
		"TypeName<int64>",
		"Parameter<c>",
		"TypeName<int64>",
		"TypeName<int64>",
		"Block<entry>",
		"ThreeAddressInstruction<MUL>",
		"Identifier<d>",
		"Identifier<a>",
		"Identifier<b>",
		"ThreeAddressInstruction<ADD>",
		"Identifier<e>",
		"Identifier<d>",
		"Identifier<c>",
		"Return<>",
		"Identifier<e>",
	}

	if node, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if err := validateTrace(node, test); err != nil {
		t.Fatal(err)
	}
}
