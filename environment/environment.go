package environment

import (
	"fmt"

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

	// types list
	env.SetType("factor", Object{Type: []*ast.Type{{Name: "factor", List: true}}})
	env.SetType("int", Object{Type: []*ast.Type{{Name: "int", List: true}}})
	env.SetType("any", Object{Type: []*ast.Type{{Name: "any", List: true}}})
	env.SetType("num", Object{Type: []*ast.Type{{Name: "num", List: true}}})
	env.SetType("char", Object{Type: []*ast.Type{{Name: "char", List: true}}})
	env.SetType("bool", Object{Type: []*ast.Type{{Name: "bool", List: true}}})
	env.SetType("null", Object{Type: []*ast.Type{{Name: "null", List: true}}})
	env.SetType("na", Object{Type: []*ast.Type{{Name: "na", List: true}}})
	env.SetType("na_char", Object{Type: []*ast.Type{{Name: "na_char", List: true}}})
	env.SetType("na_int", Object{Type: []*ast.Type{{Name: "na_int", List: true}}})
	env.SetType("na_real", Object{Type: []*ast.Type{{Name: "na_real", List: true}}})
	env.SetType("na_complex", Object{Type: []*ast.Type{{Name: "na_complex", List: true}}})
	env.SetType("nan", Object{Type: []*ast.Type{{Name: "nan", List: true}}})

	// objects
	env.SetType("list", Object{Type: []*ast.Type{{Name: "list", List: false}}})
	env.SetType("object", Object{Type: []*ast.Type{{Name: "object", List: false}}})
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

func (e *Environment) GetType(name string, list bool) (Object, bool) {
	n := name
	if list == true {
		n += "_"
	}
	obj, ok := e.types[n]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name, list)
	}
	return obj, ok
}

func (e *Environment) SetType(name string, val Object) Object {
	if val.List == true {
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
