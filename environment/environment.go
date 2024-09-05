package environment

import (
	"fmt"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/r"
)

type Environment struct {
	variables  map[string]Variable
	types      map[string]Type
	functions  map[string]Function
	class      map[string]Class
	matrix     map[string]Matrix
	returnType ast.Types
	outer      *Environment
}

func Enclose(outer *Environment, t ast.Types) *Environment {
	env := New()
	env.returnType = t
	env.outer = outer
	return env
}

func (env *Environment) ReturnType() ast.Types {
	ret := env.returnType
	if ret == nil && env.outer != nil {
		return env.outer.ReturnType()
	}
	return ret
}

func Open(env *Environment) *Environment {
	return env.outer
}

// types
var baseTypes = []string{
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
}

// objects
var baseObjects = []string{
	"list",
	"object",
	"matrix",
	"dataframe",
}

func New() *Environment {
	v := make(map[string]Variable)
	t := make(map[string]Type)
	f := make(map[string]Function)
	c := make(map[string]Class)
	m := make(map[string]Matrix)

	env := &Environment{
		functions: f,
		variables: v,
		types:     t,
		class:     c,
		matrix:    m,
		outer:     nil,
	}

	for _, t := range baseTypes {
		env.SetType(t, Type{Used: true, Type: []*ast.Type{{Name: t, List: false}}})
	}

	for _, t := range baseObjects {
		env.SetType(t, Type{Used: true, Type: []*ast.Type{{Name: t, List: false}}})
	}

	fns, err := r.ListBaseFunctions()

	if err != nil {
		fmt.Printf("failed to fetch base R functions: %v", err.Error())
		return env
	}

	for _, pkg := range fns {
		for _, fn := range pkg.Functions {
			env.SetFunction(fn, Function{Value: &ast.FunctionLiteral{}, Package: pkg.Name})
		}
	}

	return env
}

func (e *Environment) SetTypeUsed(name string) (Type, bool) {
	obj, ok := e.types[name]

	if !ok && e.outer != nil {
		return e.outer.SetTypeUsed(name)
	}

	obj.Used = true
	e.types[name] = obj

	return obj, ok
}

func (e *Environment) GetVariable(name string, outer bool) (Variable, bool) {
	obj, ok := e.variables[name]
	if !ok && e.outer != nil && outer {
		obj, ok = e.outer.GetVariable(name, outer)
	}
	return obj, ok
}

func (e *Environment) GetVariableParent(name string) (Variable, bool) {
	if e.outer == nil {
		return Variable{}, false
	}
	obj, ok := e.outer.variables[name]
	return obj, ok
}

func (e *Environment) SetVariable(name string, val Variable) Variable {
	e.variables[name] = val
	return val
}

func (e *Environment) SetVariableUsed(name string) (Variable, bool) {
	obj, ok := e.variables[name]

	if !ok && e.outer != nil {
		return e.outer.SetVariableUsed(name)
	}

	obj.Used = true
	e.variables[name] = obj

	return obj, ok
}

func (e *Environment) SetVariableNotMissing(name string) {
	v, exists := e.GetVariable(name, false)

	if !exists {
		return
	}

	v.CanMiss = false
	e.SetVariable(name, v)
}

func (e *Environment) GetType(name string) (Type, bool) {
	obj, ok := e.types[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(name)
	}

	if ok {
		e.SetTypeUsed(name)
	}
	return obj, ok
}

func (e *Environment) SetType(name string, val Type) Type {
	e.types[name] = val
	return val
}

func (e *Environment) GetFunction(name string, outer bool) (Function, bool) {
	obj, ok := e.functions[name]
	if !ok && e.outer != nil && outer {
		obj, ok = e.outer.GetFunction(name, outer)
	}
	return obj, ok
}

func (e *Environment) SetFunction(name string, val Function) Function {
	e.functions[name] = val
	return val
}

func (e *Environment) GetClass(name string) (Class, bool) {
	obj, ok := e.class[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetClass(name)
	}
	return obj, ok
}

func (e *Environment) SetClass(name string, val Class) Class {
	e.class[name] = val
	return val
}

func (e *Environment) GetMatrix(name string) (Matrix, bool) {
	obj, ok := e.matrix[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetMatrix(name)
	}
	return obj, ok
}

func (e *Environment) SetMatrix(name string, val Matrix) Matrix {
	e.matrix[name] = val
	return val
}

func (e *Environment) Types() map[string]Type {
	return e.types
}

func (e *Environment) Variables() map[string]Variable {
	return e.variables
}
