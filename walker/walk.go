package walker

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/environment"
)

type Walker struct {
	errors diagnostics.Diagnostics
	env    *environment.Environment
}

func New() *Walker {
	return &Walker{
		env: environment.New(),
	}
}

func (w *Walker) Run(node ast.Node) {
	w.Walk(node)
}

func (w *Walker) Walk(node ast.Node) (ast.Types, ast.Node) {
	var types []*ast.Type

	switch node := node.(type) {

	case *ast.Program:
		return w.walkProgram(node)

	case *ast.ExpressionStatement:
		if node.Expression != nil {
			return w.Walk(node.Expression)
		}

	case *ast.Square:
		return w.walkSquare(node)

	case *ast.LetStatement:
		return w.walkLetStatement(node)

	case *ast.ConstStatement:
		return w.walkConstStatement(node)

	case *ast.ReturnStatement:
		return w.walkReturnStatement(node)

	case *ast.TypeStatement:
		w.walkTypeStatement(node)

	case *ast.DecoratorClass:
		w.walkDecoratorClass(node)

	case *ast.DecoratorGeneric:
		w.walkDecoratorGeneric(node)

	case *ast.DecoratorDefault:
		w.walkDecoratorDefault(node)

	case *ast.Keyword:
		return ast.Types{node.Type}, node

	case *ast.Null:
		return ast.Types{node.Type}, node

	case *ast.CommentStatement:
		return types, node

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			w.Walk(s)
		}

	case *ast.Attrbute:
		return types, node

	case *ast.Identifier:
		w.walkIdentifier(node)

	case *ast.Boolean:
		return ast.Types{node.Type}, node

	case *ast.IntegerLiteral:
		return ast.Types{node.Type}, node

	case *ast.FloatLiteral:
		return ast.Types{node.Type}, node

	case *ast.VectorLiteral:
		return w.walkVectorLiteral(node)

	case *ast.SquareRightLiteral:
		return types, node

	case *ast.StringLiteral:
		return ast.Types{node.Type}, node

	case *ast.PrefixExpression:
		return w.Walk(node.Right)

	case *ast.For:
		w.env = w.env.Enclose()
		w.Walk(node.Name)
		w.Walk(node.Vector)
		t, n := w.Walk(node.Value)
		w.env = w.env.Open()
		return t, n

	case *ast.While:
		w.Walk(node.Statement)
		w.env = w.env.Enclose()
		t, n := w.Walk(node.Value)
		w.env = w.env.Open()
		return t, n

	case *ast.InfixExpression:
		return w.walkInfixExpression(node)

	case *ast.IfExpression:
		w.Walk(node.Condition)

		w.env = w.env.Enclose()
		w.Walk(node.Consequence)
		w.env = w.env.Open()

		if node.Alternative != nil {
			w.env = w.env.Enclose()
			w.Walk(node.Alternative)
			w.env = w.env.Open()
		}

	case *ast.FunctionLiteral:
		w.walkFunctionLiteral(node)

	case *ast.CallExpression:
		return w.walkCallExpression(node)
	}

	return types, node
}

func (w *Walker) walkProgram(program *ast.Program) ([]*ast.Type, ast.Node) {
	var node ast.Node
	var types []*ast.Type

	for _, statement := range program.Statements {
		types, node = w.Walk(statement)

		switch n := node.(type) {
		case *ast.ReturnStatement:
			if n.ReturnValue != nil {
				w.Walk(n.ReturnValue)
			}
		}
	}

	return types, node
}

func (w *Walker) walkCallExpression(node *ast.CallExpression) ([]*ast.Type, ast.Node) {
	for _, v := range node.Arguments {
		w.Walk(v.Value)
	}

	return w.Walk(node.Function)
}

func (w *Walker) walkInfixExpression(node *ast.InfixExpression) ([]*ast.Type, ast.Node) {
	lt, ln := w.Walk(node.Left)

	if node.Operator == "=" {
		w.Walk(node.Right)
	}

	return lt, ln
}

func (w *Walker) walkLetStatement(node *ast.LetStatement) (ast.Types, ast.Node) {
	return w.Walk(node.Value)
}

func (w *Walker) walkConstStatement(node *ast.ConstStatement) (ast.Types, ast.Node) {
	return w.Walk(node.Value)
}

func (w *Walker) walkReturnStatement(node *ast.ReturnStatement) (ast.Types, ast.Node) {
	return w.Walk(node.ReturnValue)
}

func (w *Walker) walkDecoratorDefault(node *ast.DecoratorDefault) {
	w.Walk(node.Func)
}

func (w *Walker) walkDecoratorGeneric(node *ast.DecoratorGeneric) {
	w.Walk(node.Func)
}

func (w *Walker) walkDecoratorClass(node *ast.DecoratorClass) {
}

func (w *Walker) walkTypeStatement(node *ast.TypeStatement) {
}

func (w *Walker) walkIdentifier(node *ast.Identifier) {
}

func (w *Walker) walkVectorLiteral(node *ast.VectorLiteral) ([]*ast.Type, ast.Node) {
	var ts ast.Types
	for _, s := range node.Value {
		t, _ := w.Walk(s)
		ts = append(ts, t...)
	}

	ok := allTypesIdentical(ts)

	if !ok {
		w.addFatalf(
			node.Token,
			"vectors of different types (%v)",
			ts,
		)
	}

	return ts, node
}

func (w *Walker) walkFunctionLiteral(node *ast.FunctionLiteral) {
	w.env = w.env.Enclose()

	if node.Body != nil {
		for _, s := range node.Body.Statements {
			w.Walk(s)
		}
	}

	w.env = w.env.Open()
}

func (w *Walker) walkSquare(node *ast.Square) ([]*ast.Type, ast.Node) {
	var types []*ast.Type
	var n ast.Node
	for _, s := range node.Statements {
		types, n = w.Walk(s)
	}

	return types, n
}
