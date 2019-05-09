package il

import (
	"testing"

	"github.com/askeladdk/cube"
)

func TestSemanticAnalysis_NameResolve1(t *testing.T) {
	source := `
	func identity(a i64) i64 {
		entry:
			ret a
	}`
	semantics := &semanticAnalysis{
		funcs: map[string]*funcInfo{},
	}
	if ast, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if ast, err := Traverse((*nameResolver)(semantics), ast); err != nil {
		t.Fatal(err)
	} else if _, err := Traverse((*labelUseResolver)(semantics), ast); err != nil {
		t.Fatal(err)
	} else {
		if finfo, ok := semantics.funcs["identity"]; !ok {
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
	semantics := &semanticAnalysis{
		funcs: map[string]*funcInfo{},
	}
	if ast, err := ParseProgram("test.cubeasm", source); err != nil {
		t.Fatal(err)
	} else if ast, err := Traverse((*nameResolver)(semantics), ast); err != nil {
		t.Fatal(err)
	} else if ast, err := Traverse((*labelUseResolver)(semantics), ast); err != nil {
		t.Fatal(err)
	} else if _, err := Traverse((*typeChecker)(semantics), ast); err != nil {
		t.Fatal(err)
	} else {
		// t.Fatal(ast)
	}
}
