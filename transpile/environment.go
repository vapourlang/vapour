package transpile

import "github.com/devOpifex/vapour/ast"

type Environment struct {
	variables map[string]ast.Node
	types     map[string]ast.Node
	functions map[string]ast.Node
	outer     *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	v := make(map[string]ast.Node)
	t := make(map[string]ast.Node)
	f := make(map[string]ast.Node)
	return &Environment{functions: f, variables: v, types: t, outer: nil}
}

func (env *Environment) addEnclosedEnvironment() {
	env = NewEnclosedEnvironment(env)
}

func (e *Environment) GetVariable(name string) (ast.Node, bool) {
	obj, ok := e.variables[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetVariable(name)
	}
	return obj, ok
}

func (e *Environment) SetVariable(name string, val ast.Node) ast.Node {
	e.variables[name] = val
	return val
}

func (e *Environment) GetType(name string) (ast.Node, bool) {
	obj, ok := e.types[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name)
	}
	return obj, ok
}

func (e *Environment) SetType(name string, val ast.Node) ast.Node {
	e.types[name] = val
	return val
}

func (e *Environment) GetFunction(name string) (ast.Node, bool) {
	obj, ok := e.functions[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetFunction(name)
	}
	return obj, ok
}

func (e *Environment) SetFunctions(name string, val ast.Node) ast.Node {
	e.functions[name] = val
	return val
}
