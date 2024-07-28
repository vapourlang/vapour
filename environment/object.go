package environment

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/token"
)

type Object struct {
	Token  token.Item
	Type   []*ast.Type
	Object []*ast.Type
	Value  string
	Name   string
	List   bool
}
