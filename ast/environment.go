package ast

type Environment struct {
	variables map[string]Node
	types     map[string]Node
	functions map[string]Node
	outer     *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	v := make(map[string]Node)
	t := make(map[string]Node)
	f := make(map[string]Node)
	return &Environment{functions: f, variables: v, types: t, outer: nil}
}

func (env *Environment) addEnclosedEnvironment() {
	env = NewEnclosedEnvironment(env)
}

func (e *Environment) GetVariable(name string) (Node, bool) {
	obj, ok := e.variables[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetVariable(name)
	}
	return obj, ok
}

func (e *Environment) SetVariable(name string, val Node) Node {
	e.variables[name] = val
	return val
}

func (e *Environment) GetType(name string) (Node, bool) {
	obj, ok := e.types[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name)
	}
	return obj, ok
}

func (e *Environment) SetType(name string, val Node) Node {
	e.types[name] = val
	return val
}

func (e *Environment) GetFunction(name string) (Node, bool) {
	obj, ok := e.functions[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetFunction(name)
	}
	return obj, ok
}

func (e *Environment) SetFunctions(name string, val Node) Node {
	e.functions[name] = val
	return val
}
