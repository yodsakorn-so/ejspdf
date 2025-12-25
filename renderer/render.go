package renderer

import (
	"fmt"

	"github.com/dop251/goja"
)

func (r *Runtime) Render(tpl string, data any) (string, error) {
	ejs := r.vm.Get("ejs")
	if ejs == nil {
		return "", fmt.Errorf("ejs is not loaded")
	}

	renderFn, ok := goja.AssertFunction(
		ejs.ToObject(r.vm).Get("render"),
	)
	if !ok {
		return "", fmt.Errorf("ejs.render is not a function")
	}

	result, err := renderFn(
		goja.Undefined(),
		r.vm.ToValue(tpl),
		r.vm.ToValue(data),
	)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
