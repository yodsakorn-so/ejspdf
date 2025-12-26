package renderer

import (
	"fmt"

	"github.com/dop251/goja"
)

func (r *Runtime) RenderEJS(ejsJS []byte, tpl string, data any) (string, error) {
	_, err := r.vm.RunString(string(ejsJS))
	if err != nil {
		return "", err
	}

	render, ok := goja.AssertFunction(r.vm.Get("ejs").ToObject(r.vm).Get("render"))
	if !ok {
		return "", fmt.Errorf("ejs.render not found")
	}

	val, err := render(goja.Undefined(), r.vm.ToValue(tpl), r.vm.ToValue(data))
	if err != nil {
		return "", err
	}

	return val.String(), nil
}
