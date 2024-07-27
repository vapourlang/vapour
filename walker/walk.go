package walker

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/environment"
	"github.com/devOpifex/vapour/err"
	"github.com/devOpifex/vapour/token"
)

type Walker struct {
	code   string
	errors err.Errors
	env    *environment.Environment
}

func New() *Walker {
	return &Walker{
		env: environment.New(),
	}
}

func (w *Walker) Walk(node ast.Node) ([]*ast.Type, ast.Node) {
	var types []*ast.Type

	switch node := node.(type) {

	case *ast.ExpressionStatement:
		if node.Expression != nil {
			return w.Walk(node.Expression)
		}

	case *ast.Program:
		return w.walkProgram(node)

	case *ast.LetStatement:
		_, exists := w.env.GetVariable(node.Name.Value, true)

		if exists {
			w.addError(node.Token, node.Name.Value+" is already declared")
		}

		w.expectType(node.Value, node.Name.Type)

	case *ast.ConstStatement:
		_, exists := w.env.GetVariable(node.Name.Value, true)

		if exists {
			w.addError(node.Token, node.Name.Value+" is already declared")
			return w.Walk(node.Value)
		}

		w.env.SetVariable(node.Name.Value, environment.Object{Token: node.Token})
		w.Walk(node.Value)

	case *ast.ReturnStatement:
		w.Walk(node.ReturnValue)

	case *ast.TypeStatement:
		_, exists := w.env.GetType(node.Name.Value)

		if exists {
			w.addError(node.Token, "type "+node.Name.Value+" already defined")
		}

		w.env.SetType(
			node.Name.Value,
			environment.Object{
				Token:  node.Token,
				Type:   node.Type,
				Object: node.Object,
				Name:   node.Name.Value,
				List:   node.List,
			},
		)

	case *ast.Keyword:
		return types, node

	case *ast.Null:
		return node.Type, node

	case *ast.CommentStatement:
		return types, node

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			w.Walk(s)
		}

	case *ast.Identifier:
		return node.Type, node

	case *ast.Boolean:
		return []*ast.Type{{Name: "bool", List: false}}, node

	case *ast.IntegerLiteral:
		return []*ast.Type{{Name: "int", List: false}}, node

	case *ast.VectorLiteral:
		for _, s := range node.Value {
			w.Walk(s)
		}

	case *ast.SquareRightLiteral:
		return types, node

	case *ast.StringLiteral:
		return types, node

	case *ast.PrefixExpression:
		return w.Walk(node.Right)

	case *ast.InfixExpression:
		w.Walk(node.Left)
		if node.Right != nil {
			return w.Walk(node.Right)
		}

	case *ast.IfExpression:
		w.Walk(node.Condition)

		if node.Alternative != nil {
			w.Walk(node.Alternative)
		}

	case *ast.FunctionLiteral:
		w.env = environment.NewEnclosed(w.env)

		params := []string{}
		for _, p := range node.Parameters {
			w.env.SetVariable(
				p.TokenLiteral(),
				environment.Object{
					Token: node.Token,
					Name:  node.Name.Value,
				},
			)
			params = append(params, p.String())
		}
		w.Walk(node.Body)

	case *ast.CallExpression:
		w.Walk(node.Function)
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

func (w *Walker) addError(tok token.Item, m string) {
	w.errors = append(w.errors, err.New(tok, m))
}

func (w *Walker) expectType(node ast.Node, actual []*ast.Type) {
	expected, _ := w.Walk(node)
	ok, expected, _ := ast.AllTypesMatch(actual, expected)

	if ok {
		return
	}

	w.addError(token.Item{}, "wrong types")
}
