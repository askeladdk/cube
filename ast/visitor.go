package ast

import "reflect"

type Visitor interface {
	Visit(Node) (Node, error)
	PostVisit(Node) (Node, error)
}

type visitor struct {
	vi Visitor
}

func (this *visitor) Visit(n Node) (Node, error) {
	if n == nil || reflect.ValueOf(n).IsNil() {
		return nil, nil
	} else if m, err := this.vi.Visit(n); err != nil {
		return nil, err
	} else if k, err := m.Traverse(this); err != nil {
		return nil, err
	} else {
		return this.PostVisit(k)
	}
}

func (this *visitor) PostVisit(n Node) (Node, error) {
	return this.vi.PostVisit(n)
}

func Traverse(vi Visitor, n Node) (Node, error) {
	vi2 := visitor{vi}
	return vi2.Visit(n)
}
