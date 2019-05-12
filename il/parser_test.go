package il

import (
	"errors"
	"fmt"
	"testing"
)

type tracer struct {
	trace []string
}

func (this *tracer) Visit(n Node) (Node, error) {
	this.trace = append(this.trace, n.String())
	return n, nil
}

func (this *tracer) PostVisit(n Node) (Node, error) {
	return n, nil
}

func validateTrace(n Node, test []string) error {
	tracer := tracer{}
	Traverse(&tracer, n)
	trace := tracer.trace

	if len(trace) != len(test) {
		return errors.New("len(trace) != len(test)")
	}

	for i, s := range trace {
		if s != test[i] {
			return errors.New(fmt.Sprintf("%d: %s != %s", i, s, test[i]))
		}
	}

	return nil
}

func TestParseFuncZero(t *testing.T) {
	source := `
    func zero() i64 {
        entry:
            ret 0
    }`

	test := []string{
		"Function<zero>",
		"Signature<>",
		"TypeName<int64>",
		"Block<entry>",
		"Return<>",
		"Integer<0>",
	}

	lexer := NewLexer("<source>", source)
	parser := NewParser(lexer)
	if err := parser.advance(); err != nil {
		t.Fatal(err)
	} else if node, err := parser.definitions(); err != nil {
		t.Fatal(err)
	} else if err := validateTrace(node, test); err != nil {
		t.Fatal(err)
	}
}

func TestParseBlockDoubleFunc(t *testing.T) {
	source := `
    func double(a i64) i64 {
		var b i64
        entry:
            add b, a, a
            ret b
    }`

	test := []string{
		"Function<double>",
		"Signature<>",
		"Parameter<a>",
		"TypeName<int64>",
		"TypeName<int64>",
		"Local<b>",
		"TypeName<int64>",
		"Block<entry>",
		"Instruction<19>",
		"Use<b>",
		"Use<a>",
		"Use<a>",
		"Return<>",
		"Use<b>",
	}

	lexer := NewLexer("<source>", source)
	parser := NewParser(lexer)
	if err := parser.advance(); err != nil {
		t.Fatal(err)
	} else if node, err := parser.definitions(); err != nil {
		t.Fatal(err)
	} else if err := validateTrace(node, test); err != nil {
		t.Fatal(err)
	}
}
