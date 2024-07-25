package transpile

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/err"
)

type Transpiler struct {
	code  []string
	error err.Errors
}

func (t *Transpiler) Transpile(node ast.Node, env *Environment) ast.Node {
	switch node := node.(type) {

	// Statements
	case *ast.ExpressionStatement:
		return t.Transpile(node.Expression, env)

	case *ast.Program:
		return t.transpileProgram(node, env)

	case *ast.LetStatement:
		env.SetVariable(node.Name.Value, node)
		t.Transpile(node.Value, env)

	case *ast.ReturnStatement:
		t.Transpile(node.ReturnValue, env)
	}

	return node
}

func (t *Transpiler) transpileProgram(program *ast.Program, env *Environment) ast.Node {
	var node ast.Node

	for _, statement := range program.Statements {
		node := t.Transpile(statement, env)

		switch n := node.(type) {
		case *ast.ReturnStatement:
			return t.Transpile(n.ReturnValue, env)
		}
	}

	return node
}
