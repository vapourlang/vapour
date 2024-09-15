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
	factor     map[string]Factor
	signature  map[string]Signature
	method     map[string]Methods
	returnType ast.Types
	outer      *Environment
}

var library string

func SetLibrary(path string) {
	library = path
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
	"factor",
	"matrix",
	"vector",
	"struct",
	"dataframe",
	"impliedList",
}

func NewGlobalEnvironment() *Environment {
	v := make(map[string]Variable)
	t := make(map[string]Type)
	f := make(map[string]Function)
	c := make(map[string]Class)
	m := make(map[string]Matrix)
	s := make(map[string]Signature)
	fct := make(map[string]Factor)
	meth := make(map[string]Methods)

	env := &Environment{
		functions: f,
		variables: v,
		types:     t,
		class:     c,
		matrix:    m,
		signature: s,
		factor:    fct,
		method:    meth,
		outer:     nil,
	}

	for _, t := range baseTypes {
		env.SetType(Type{Used: true, Name: t, Type: []*ast.Type{{Name: t, List: false}}})
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

func New() *Environment {
	v := make(map[string]Variable)
	t := make(map[string]Type)
	f := make(map[string]Function)
	c := make(map[string]Class)
	m := make(map[string]Matrix)
	s := make(map[string]Signature)
	fct := make(map[string]Factor)
	meth := make(map[string]Methods)

	return &Environment{
		functions: f,
		variables: v,
		types:     t,
		class:     c,
		matrix:    m,
		signature: s,
		factor:    fct,
		method:    meth,
		outer:     nil,
	}
}

func (e *Environment) SetSignatureUsed(name string) (Signature, bool) {
	obj, ok := e.signature[name]

	if !ok && e.outer != nil {
		return e.outer.SetSignatureUsed(name)
	}

	obj.Used = true
	e.signature[name] = obj

	return obj, ok
}

func (e *Environment) SetTypeUsed(pkg, name string) (Type, bool) {
	obj, ok := e.types[makeTypeKey(pkg, name)]

	if !ok && e.outer != nil {
		return e.outer.SetTypeUsed(pkg, name)
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

	obj, ok := e.outer.GetVariable(name, true)

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

func makeTypeKey(pkg, name string) string {
	if pkg == "" {
		return name
	}

	return pkg + "::" + name
}

func (e *Environment) GetType(pkg, name string) (Type, bool) {
	obj, ok := e.types[makeTypeKey(pkg, name)]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetType(pkg, name)
	}

	if ok {
		e.SetTypeUsed(pkg, name)
	}

	return obj, ok
}

func (e *Environment) SetType(val Type) Type {
	e.types[makeTypeKey(val.Package, val.Name)] = val
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

func (e *Environment) GetFactor(name string) (Factor, bool) {
	obj, ok := e.factor[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetFactor(name)
	}
	return obj, ok
}

func (e *Environment) SetFactor(name string, val Factor) Factor {
	e.factor[name] = val
	return val
}

func (e *Environment) GetSignature(name string) (Signature, bool) {
	obj, ok := e.signature[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetSignature(name)
	}
	return obj, ok
}

func (e *Environment) SetSignature(name string, val Signature) Signature {
	e.signature[name] = val
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

func (e *Environment) AddMethod(name string, val Method) Method {
	e.method[name] = append(e.method[name], val)
	return val
}

func (e *Environment) GetMethods(name string) (Methods, bool) {
	obj, ok := e.method[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetMethods(name)
	}
	return obj, ok
}

func (e *Environment) GetMethod(name string, t *ast.Type) (Method, bool) {
	obj, ok := e.method[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetMethods(name)
	}

	if !ok {
		return Method{}, false
	}

	for _, o := range obj {
		if o.Value.Method == t {
			return o, true
		}
	}

	return Method{}, false
}

func (e *Environment) HasMethods(name string, t *ast.Type) bool {
	obj, ok := e.GetMethods(name)

	if !ok {
		return false
	}

	for _, o := range obj {
		if o.Value.Method == t {
			return true
		}
	}

	return false
}

func (e *Environment) Types() map[string]Type {
	return e.types
}

func (e *Environment) Variables() map[string]Variable {
	return e.variables
}
