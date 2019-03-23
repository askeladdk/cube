package il

import "github.com/askeladdk/cube/ast"

func ParseProgram(filename, source string) (ast.Node, error) {
	lexer := NewLexer(filename, source)
	parser := NewParser(lexer)
	if node, err := parser.Parse(nil); err != nil {
		return nil, err
	} else {
		return &ast.Program{node}, nil
	}
}
