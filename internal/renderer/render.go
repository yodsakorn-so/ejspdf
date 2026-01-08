package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

func (r *Runtime) RenderEJS(ejsJS []byte, tpl string, data any, filename string) (string, error) {
	// 1. Inject 'fs' and 'path' mocks for EJS include support
	r.setupNodePolyfills()

	// 2. Load EJS library
	// HACK: The bundled EJS has an empty mock for 'fs' at module ID 1.
	// We replace it to redirect calls to our 'native-fs'.
	jsCode := string(ejsJS)
	target := "1:[function(require,module,exports){"
	replacement := "1:[function(require,module,exports){module.exports=require('native-fs');"
	
	if !strings.Contains(jsCode, target) {
		// If exact match fails, try a more flexible search or report error
		return "", fmt.Errorf("failed to patch EJS: target signature not found. This might be due to a version mismatch or unexpected file encoding")
	}
	
	jsCode = strings.Replace(jsCode, target, replacement, 1)

	_, err := r.vm.RunString(jsCode)
	if err != nil {
		return "", err
	}
	
	// Override ejs.resolveInclude to use our native path module
	// First, expose 'path' and 'fs' (native-fs) to global scope
	r.vm.RunString(`
		var nativePath = require('path');
		var nativeFS = require('native-fs');

		ejs.resolveInclude = function(name, filename, isDir) {
			var dirname = nativePath.dirname(filename);
			var resolved = nativePath.resolve(dirname, name);
			if (!nativeFS.existsSync(resolved)) {
				// EJS standard behavior: try adding .ejs extension if missing
				// simplified here
			}
			return resolved;
		}
	`)

	render, ok := goja.AssertFunction(r.vm.Get("ejs").ToObject(r.vm).Get("render"))
	if !ok {
		return "", fmt.Errorf("ejs.render not found")
	}

	// 3. Prepare Options
	opts := r.vm.NewObject()
	if filename != "" {
		opts.Set("filename", filename)
	}

	// 4. Render
	// ejs.render(template, data, options)
	val, err := render(goja.Undefined(), r.vm.ToValue(tpl), r.vm.ToValue(data), opts)
	if err != nil {
		return "", err
	}

	return val.String(), nil
}

func (r *Runtime) setupNodePolyfills() {
	// Mock 'fs' module
	fsObj := r.vm.NewObject()
	
	// Use explicit variable for function to ensure it's not GC'd or lost (though not likely issue)
	readFileFunc := func(call goja.FunctionCall) goja.Value {
		pathVar := call.Argument(0).String()
		b, err := os.ReadFile(pathVar)
		if err != nil {
			panic(r.vm.ToValue(fmt.Sprintf("fs.readFileSync failed: %v", err)))
		}
		return r.vm.ToValue(string(b))
	}
	fsObj.Set("readFileSync", readFileFunc)

	existsFunc := func(call goja.FunctionCall) goja.Value {
		pathVar := call.Argument(0).String()
		_, err := os.Stat(pathVar)
		exists := err == nil || !os.IsNotExist(err)
		return r.vm.ToValue(exists)
	}
	fsObj.Set("existsSync", existsFunc)
	
	// Mock 'path' module (EJS uses path.resolve/join sometimes)
	pathObj := r.vm.NewObject()
	pathObj.Set("resolve", func(call goja.FunctionCall) goja.Value {
		// Simple implementation: join all args
		var paths []string
		for _, arg := range call.Arguments {
			paths = append(paths, arg.String())
		}
		return r.vm.ToValue(filepath.Join(paths...))
	})
	pathObj.Set("join", func(call goja.FunctionCall) goja.Value {
		var paths []string
		for _, arg := range call.Arguments {
			paths = append(paths, arg.String())
		}
		return r.vm.ToValue(filepath.Join(paths...))
	})
	pathObj.Set("dirname", func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).String()
		return r.vm.ToValue(filepath.Dir(input))
	})
	pathObj.Set("extname", func(call goja.FunctionCall) goja.Value {
		return r.vm.ToValue(filepath.Ext(call.Argument(0).String()))
	})

	// Inject 'require' to return our mocks
	r.vm.Set("require", func(call goja.FunctionCall) goja.Value {
		moduleName := call.Argument(0).String()
		switch moduleName {
		case "fs", "native-fs":
			return fsObj
		case "path":
			return pathObj
		default:
			// For other modules, return undefined or empty object
			return r.vm.NewObject()
		}
	})
}
