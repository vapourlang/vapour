package environment

import (
	"github.com/devOpifex/vapour/ast"
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

	// types
	env.SetType("factor", Object{Type: []*ast.Type{{Name: "factor", List: false}}})
	env.SetType("int", Object{Type: []*ast.Type{{Name: "int", List: false}}})
	env.SetType("any", Object{Type: []*ast.Type{{Name: "any", List: false}}})
	env.SetType("num", Object{Type: []*ast.Type{{Name: "num", List: false}}})
	env.SetType("char", Object{Type: []*ast.Type{{Name: "char", List: false}}})
	env.SetType("bool", Object{Type: []*ast.Type{{Name: "bool", List: false}}})
	env.SetType("null", Object{Type: []*ast.Type{{Name: "null", List: false}}})
	env.SetType("na", Object{Type: []*ast.Type{{Name: "na", List: false}}})
	env.SetType("na_char", Object{Type: []*ast.Type{{Name: "na_char", List: false}}})
	env.SetType("na_int", Object{Type: []*ast.Type{{Name: "na_int", List: false}}})
	env.SetType("na_real", Object{Type: []*ast.Type{{Name: "na_real", List: false}}})
	env.SetType("na_complex", Object{Type: []*ast.Type{{Name: "na_complex", List: false}}})
	env.SetType("nan", Object{Type: []*ast.Type{{Name: "nan", List: false}}})

	// objects
	env.SetType("list", Object{Type: []*ast.Type{{Name: "list", List: false}}})
	env.SetType("object", Object{Type: []*ast.Type{{Name: "object", List: false}}})
	env.SetType("matrix", Object{Type: []*ast.Type{{Name: "matrix", List: false}}})
	env.SetType("dataframe", Object{Type: []*ast.Type{{Name: "dataframe", List: false}}})

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

func (e *Environment) GetType(name string, list bool) (Object, bool) {
	n := name
	if list {
		n += "_"
	}
	obj, ok := e.types[n]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name, list)
	}
	return obj, ok
}

func (e *Environment) SetType(name string, val Object) Object {
	if val.List {
		name += "_"
	}
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
