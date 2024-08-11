package environment

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/token"
)

// this should be an interface but I haven't got the time right now
type Object struct {
	Token      token.Item
	Type       []*ast.Type
	Object     []*ast.Type
	Value      string
	Name       string
	List       bool
	Const      bool
	Parameters []Object
	Class      []string
	Attributes []*ast.TypeAttributesStatement
	Used       bool
	CanMiss    bool
}
