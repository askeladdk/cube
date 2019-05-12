package il

import (
	"testing"

	"github.com/askeladdk/cube"
)

func TestSemanticAnalysis_NameResolve1(t *testing.T) {
	source := `
	func identity(a i64) i64 {
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
		} else if aloc, ok := finfo.locals["a"]; !ok {
			t.Fatalf("local a not found")
		} else if aloc.dtype != cube.TypeInt64 {
			t.Fatalf("a type != int64")
		}
	}
}

func TestSemanticAnalysis_TypeCheck(t *testing.T) {
	source := `
	func inc(a i64) i64 {
		entry:
			set b, 0
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
