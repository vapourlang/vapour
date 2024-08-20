package environment

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/token"
)

// this should be an interface but I haven't got the time right now
type Function struct {
	Token   token.Item
	Package string
	Value   ast.FunctionLiteral
}

type Variable struct {
	Token    token.Item
	Value    []*ast.Type
	HasValue bool
	CanMiss  bool
	IsConst  bool
	List     bool
	Used     bool
	Name     string
}

type Type struct {
	Token      token.Item
	Name       string
	Type       []*ast.Type
	Used       bool
	Attributes []*ast.TypeAttributesStatement
}

type Class struct {
	Token token.Item
	Value *ast.DecoratorClass
}
