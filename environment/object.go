package environment

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/token"
)

// this should be an interface but I haven't got the time right now
type Function struct {
	Token   token.Item
	Package string
	Value   *ast.FunctionLiteral
}

type Variable struct {
	Token    token.Item
	Value    []*ast.Type
	HasValue bool
	CanMiss  bool
	IsConst  bool
	Used     bool
}

type Type struct {
	Token      token.Item
	Type       []*ast.Type
	Used       bool
	Object     string
	Attributes []*ast.TypeAttributesStatement
}

type Class struct {
	Token token.Item
	Value *ast.DecoratorClass
}
