package renderer

import (
	_ "embed"

	"github.com/dop251/goja"
)

//go:embed assets/ejs.js
var ejsSource string

type Runtime struct {
	vm *goja.Runtime
}

func New() (*Runtime, error) {
	vm := goja.New()

	if _, err := vm.RunString(ejsSource); err != nil {
		return nil, err
	}

	return &Runtime{
		vm: vm,
	}, nil
}
