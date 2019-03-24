package il

func ParseProgram(filename, source string) (Node, error) {
	lexer := NewLexer(filename, source)
	parser := NewParser(lexer)
	if node, err := parser.Parse(nil); err != nil {
		return nil, err
	} else {
		return &Program{node}, nil
	}
}
