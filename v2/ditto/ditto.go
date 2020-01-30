package ditto

import (
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Ditto struct {
	name string
	vm *jsonnet.VM
	ast ast.Node
}

var cache = make(map[string]*Ditto)

func New(name string, net string) *Ditto {
	if cache[name] != nil {
		return cache[name]
	}

	vm := jsonnet.MakeVM()
	a, err := jsonnet.SnippetToAST(name, net)
	if nil != err {
		panic(err)
	}

	d := &Ditto{
		name: name,
		vm:   vm,
		ast: a,
	}

	cache[name] = d
	return d
}

func (d *Ditto) JSON(net string) (string, error) {
	return d.vm.Evaluate(d.ast)
}