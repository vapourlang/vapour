package environment

import "fmt"

type Environment struct {
	variables map[string]Object
	types     map[string]Object
	functions map[string]Object
	outer     *Environment
}

func (e *Environment) Enclose() *Environment {
	env := New()
	env.outer = e
	return env
}

func (e *Environment) Open() *Environment {
	return e.outer
}

func New() *Environment {
	v := make(map[string]Object)
	t := make(map[string]Object)
	f := make(map[string]Object)

	env := &Environment{functions: f, variables: v, types: t, outer: nil}

	// types
	env.SetType("int", Object{})
	env.SetType("num", Object{})
	env.SetType("null", Object{})
	env.SetType("na", Object{})
	env.SetType("na_char", Object{})
	env.SetType("na_int", Object{})
	env.SetType("na_real", Object{})
	env.SetType("na_complex", Object{})
	env.SetType("nan", Object{})
	env.SetType("char", Object{})
	env.SetType("bool", Object{})

	// objects
	env.SetType("list", Object{})
	env.SetType("matrix", Object{})
	env.SetType("dataframe", Object{})

	return env
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

func (e *Environment) Print() {
	fmt.Println("++++++++++++++++++++++++++++ Environment")
	fmt.Println("--- Inner")
	fmt.Println("- Variables")
	for k := range e.variables {
		fmt.Printf("%v\n", k)
	}
	fmt.Println("- Functions")
	for k := range e.functions {
		fmt.Printf("%v\n", k)
	}

	fmt.Println("- Types")
	for k := range e.types {
		fmt.Printf("%v\n", k)
	}

	fmt.Println("--- Outer")
	fmt.Println("- Variables")
	if e.outer != nil {
		for k := range e.outer.variables {
			fmt.Printf("%v\n", k)
		}
	}
	fmt.Println("- Functions")
	if e.outer != nil {
		for k := range e.outer.functions {
			fmt.Printf("%v\n", k)
		}
	}
	fmt.Println("- Types")
	if e.outer != nil {
		for k := range e.outer.types {
			fmt.Printf("%v\n", k)
		}
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++")
}
