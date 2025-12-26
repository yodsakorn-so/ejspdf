package renderer

import "github.com/dop251/goja"

type Runtime struct {
	vm *goja.Runtime
}

func New() *Runtime {
	return &Runtime{
		vm: goja.New(),
	}
}
