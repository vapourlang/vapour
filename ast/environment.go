package ast

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Node)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]Node
	outer *Environment
}

func (e *Environment) Get(name string) (Node, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Node) Node {
	e.store[name] = val
	return val
}
