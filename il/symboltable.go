package il

import "github.com/askeladdk/cube"

type localSymbol struct {
	parent int
	index  int
	dtype  *cube.Type
	param  bool
}

type blockSymbol struct {
	name   string
	index  int
	params []int
}

type funcSymbol struct {
	name   string
	level  int
	locals []*localSymbol
	blocks []*blockSymbol
	rtype  *cube.Type
}

type symbolTable map[string]*funcSymbol
