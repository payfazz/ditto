package ditto

import (
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Ditto struct {
	name string
	ast ast.Node
}

var cache = make(map[string]*Ditto)

func New(name string, net string) *Ditto {
	if cache[name] != nil {
		return cache[name]
	}

	a, err := jsonnet.SnippetToAST(name, net)
	if nil != err {
		panic(err)
	}

	d := &Ditto{
		name: name,
		ast: a,
	}

	cache[name] = d
	return d
}

func (d *Ditto) JSON() (string, error) {
	vm := jsonnet.MakeVM()
	return vm.Evaluate(d.ast)
}