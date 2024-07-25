package environment

import "github.com/devOpifex/vapour/ast"

type Environment struct {
	variables map[string]Object
	types     map[string]Object
	functions map[string]Object
	outer     *Environment
}

func NewEnclosed(outer *Environment) *Environment {
	env := New()
	env.outer = outer
	return env
}

func New() *Environment {
	v := make(map[string]Object)
	t := make(map[string]Object)
	f := make(map[string]Object)
	return &Environment{functions: f, variables: v, types: t, outer: nil}
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

func (e *Environment) GetType(name string) (Object, bool) {
	obj, ok := e.types[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name)
	}
	return obj, ok
}

func (e *Environment) HasAllTypes(types []*ast.Type) ([]string, bool) {
	var has []bool
	var missing []string
	for _, v := range types {
		_, ok := e.types[v.Name]
		if !ok && e.outer != nil {
			_, ok = e.outer.GetType(v.Name)
		}

		has = append(has, ok)

		if !ok {
			missing = append(missing, v.Name)
		}
	}

	exists := true

	for _, v := range has {
		if !v {
			exists = false
		}
	}

	return missing, exists
}

func (e *Environment) SetType(name string, val Object) Object {
	e.types[name] = val
	return val
}

func (e *Environment) GetFunction(name string) (Object, bool) {
	obj, ok := e.functions[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetFunction(name)
	}
	return obj, ok
}

func (e *Environment) SetFunctions(name string, val Object) Object {
	e.functions[name] = val
	return val
}
