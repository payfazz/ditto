package ditto

import (
	"fmt"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Ditto struct {
	name string
	ast  ast.Node
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
		ast:  a,
	}

	cache[name] = d
	return d
}

func PrepareVM(ext map[string]interface{}, tla map[string]interface{}) *jsonnet.VM {
	vm := jsonnet.MakeVM()

	for k, v := range ext {
		if val, ok := v.(string); ok {
			vm.ExtVar(k, val)
		} else {
			str := fmt.Sprintf("%v", v)
			vm.ExtCode(k, str)
		}
	}

	for k, v := range tla {
		if val, ok := v.(string); ok {
			vm.TLAVar(k, val)
		} else {
			str := fmt.Sprintf("%v", v)
			vm.TLACode(k, str)
		}
	}
	return vm
}

func (d *Ditto) JSON(vm *jsonnet.VM) (string, error) {
	if vm == nil {
		vm = jsonnet.MakeVM()
	}
	return vm.Evaluate(d.ast)
}
