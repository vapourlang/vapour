package environment

import (
	"github.com/vapourlang/vapour/ast"
	"github.com/vapourlang/vapour/token"
)

// this should be an interface but I haven't got the time right now
type Function struct {
	Token   token.Item
	Package string
	Value   *ast.FunctionLiteral
	Name    string
}

type Methods []Method

type Method struct {
	Token   token.Item
	Package string
	Value   *ast.FunctionLiteral
	Name    string
}

type Variable struct {
	Token    token.Item
	Value    ast.Types
	HasValue bool
	CanMiss  bool
	IsConst  bool
	Used     bool
	Name     string
}

type Type struct {
	Token      token.Item
	Type       ast.Types
	Package    string
	Used       bool
	Object     string
	Name       string
	Attributes []*ast.TypeAttributesStatement
}

type Class struct {
	Token token.Item
	Value *ast.DecoratorClass
}

type Matrix struct {
	Token token.Item
	Value *ast.DecoratorMatrix
}

type Factor struct {
	Token token.Item
	Value *ast.DecoratorFactor
}

type Env struct {
	Token token.Item
	Value *ast.DecoratorEnvironment
}

type Signature struct {
	Token token.Item
	Value *ast.TypeFunction
	Used  bool
}
