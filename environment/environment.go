package environment

import (
	"fmt"
	"reflect"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/r"
)

type Environment struct {
	variables map[string]Object
	types     map[string]Object
	functions map[string]Object
	class     map[string]Object
	Fn        Object // Function (if environment is a function)
	outer     *Environment
}

func (e *Environment) Enclose(fn Object) *Environment {
	env := New(fn)
	env.outer = e
	return env
}

func (e *Environment) Open() *Environment {
	return e.outer
}

var baseTypes = []string{
	// types
	"factor",
	"int",
	"any",
	"num",
	"char",
	"bool",
	"null",
	"na",
	"na_char",
	"na_int",
	"na_real",
	"na_complex",
	"nan",
	// objects
	"list",
	"object",
	"matrix",
	"dataframe",
}

func New(fn Object) *Environment {
	v := make(map[string]Object)
	t := make(map[string]Object)
	f := make(map[string]Object)
	c := make(map[string]Object)

	env := &Environment{
		functions: f,
		variables: v,
		types:     t,
		class:     c,
		outer:     nil,
		Fn:        fn,
	}

	for _, t := range baseTypes {
		env.SetType(t, Object{Type: []*ast.Type{{Name: t, List: false}}})
	}

	fns, err := r.ListBaseFunctions()

	if err != nil {
		fmt.Printf("failed to fetch base R functions: %v", err.Error())
		return env
	}

	for _, pkg := range fns {
		for _, fn := range pkg.Functions {
			env.SetFunction(fn, Object{Name: fn, Package: pkg.Name})
		}
	}

	return env
}

func (e *Environment) variablesNotUsed() []Object {
	var unused []Object
	for _, v := range e.variables {
		if !v.Used && v.Name != "..." {
			unused = append(unused, v)
		}
	}

	return unused
}

func (e *Environment) AllVariablesUsed() ([]Object, bool) {
	v := e.variablesNotUsed()
	return v, len(v) == 0
}

func (e *Environment) GetVariable(name string, outer bool) (Object, bool) {
	obj, ok := e.variables[name]
	if !ok && e.outer != nil && outer {
		obj, ok = e.outer.GetVariable(name, outer)
	}
	return obj, ok
}

func (e *Environment) SetVariable(name string, val Object) Object {
	e.variables[name] = val
	return val
}

func (e *Environment) SetVariableUsed(name string) {
	v, exists := e.GetVariable(name, false)

	if !exists {
		return
	}

	v.Used = true
	e.SetVariable(name, v)
}

func (e *Environment) SetVariableNotMissing(name string) {
	v, exists := e.GetVariable(name, false)

	if !exists {
		return
	}

	v.CanMiss = false
	e.SetVariable(name, v)
}

func (e *Environment) GetType(name string) (Object, bool) {
	obj, ok := e.types[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name)
	}
	return obj, ok
}

func (e *Environment) SetType(name string, val Object) Object {
	e.types[name] = val
	return val
}

func (e *Environment) GetFunction(name string, outer bool) (Object, bool) {
	obj, ok := e.functions[name]
	if !ok && e.outer != nil && outer {
		obj, ok = e.outer.GetFunction(name, outer)
	}
	return obj, ok
}

func (e *Environment) SetFunction(name string, val Object) Object {
	e.functions[name] = val
	return val
}

func (e *Environment) GetFunctionEnvironment() (bool, Object) {
	var exists bool

	if e.Fn.Token.Value != "" {
		exists = true
	}

	return exists, e.Fn
}

func (e *Environment) GetClass(name string) (Object, bool) {
	obj, ok := e.class[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetClass(name)
	}
	return obj, ok
}

func (e *Environment) SetClass(name string, val Object) Object {
	e.class[name] = val
	return val
}

func (e *Environment) GetTypeFromSignature(fn *ast.FunctionLiteral) (string, bool) {
	for name, types := range e.types {
		for _, t := range types.Type {
			if t.Name != "func" {
				continue
			}

			var allIdentical []bool
			for i, t := range types.Attributes {
				if len(fn.Parameters)-1 < i {
					continue
				}
				identical := reflect.DeepEqual(t.Type, fn.Parameters[i].Type)
				allIdentical = append(allIdentical, identical)
			}

			for _, i := range allIdentical {
				if !i {
					return "", false
				}
			}

			if len(fn.Type) == 0 {
				return "", false
			}

			if len(types.Object) < 2 {
				return "", false
			}

			// check return type
			if fn.Type[0].Name == types.Object[1].Name {
				return name, true
			}

			return "", false
		}
	}

	return "", false
}
