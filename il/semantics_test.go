package il

import (
	"testing"

	"github.com/askeladdk/cube"
)

func TestSemanticAnalysis_NameResolve1(t *testing.T) {
	source := `
	func identity(a int64) int64 {
		entry:
			goto done
		done:
			ret a
	}`
	symbolTable := symbolTable{}
	resolver := newSymbolResolver(symbolTable)
	if ast, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if _, err := resolver.DoPass(ast); err != nil {
		t.Fatal(err)
	} else {
		if finfo, ok := symbolTable["identity"]; !ok {
			t.Fatalf("function not found")
		} else if finfo.level != 1 {
			t.Fatalf("level != 1")
		} else if len(finfo.locals) == 0 {
			t.Fatalf("no locals")
		} else if finfo.locals[0].dtype != cube.TypeInt64 {
			t.Fatalf("a type != int64")
		}
	}
}

func TestSemanticAnalysis_TypeCheck(t *testing.T) {
	source := `
	func inc(a int64) int64 {
		var b int64
		entry:
			add b, a, 1
			ret b
	}`
	symbolTable := symbolTable{}
	resolver := newSymbolResolver(symbolTable)
	typeChecker := &typeChecker{
		symbolTable: symbolTable,
	}

	if ast, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if ast, err := resolver.DoPass(ast); err != nil {
		t.Fatal(err)
	} else if _, err := Traverse(typeChecker, ast); err != nil {
		t.Fatal(err)
	} else {
		// t.Fatal(ast)
	}
}
