package environment

import (
	"fmt"

	"github.com/devOpifex/vapour/ast"
)

type Environment struct {
	variables map[string]Object
	types     map[string]Object
	functions map[string]Object
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

	env := &Environment{
		functions: f,
		variables: v,
		types:     t,
		outer:     nil,
		Fn:        fn,
	}

	// types
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
	env.SetType("matrix", Object{Type: []*ast.Type{{Name: "matrix", List: false}}})
	env.SetType("dataframe", Object{Type: []*ast.Type{{Name: "dataframe", List: false}}})

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

func (e *Environment) GetFunctionEnvironment() (bool, Object) {
	var exists bool

	if e.Fn.Token.Value != "" {
		exists = true
	}

	return exists, e.Fn
}

func (e *Environment) Print() {
	fmt.Println("++++++++++++++++++++++++++++ Environment")
	fmt.Println("------ Inner")
	fmt.Println("--- Variables")
	for k := range e.variables {
		fmt.Printf("%v\n", k)
	}
	fmt.Println("--- Functions")
	for k := range e.functions {
		fmt.Printf("%v\n", k)
	}

	fmt.Println("--- Types")
	for k := range e.types {
		fmt.Printf("%v\n", k)
	}

	fmt.Println("------ Outer")
	fmt.Println("--- Variables")
	if e.outer != nil {
		for k := range e.outer.variables {
			fmt.Printf("%v\n", k)
		}
	}
	fmt.Println("--- Functions")
	if e.outer != nil {
		for k := range e.outer.functions {
			fmt.Printf("%v\n", k)
		}
	}
	fmt.Println("--- Types")
	if e.outer != nil {
		for k := range e.outer.types {
			fmt.Printf("%v\n", k)
		}
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++")
}
