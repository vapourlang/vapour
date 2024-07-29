package walker

import (
	"fmt"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/environment"
	"github.com/devOpifex/vapour/token"
)

type Walker struct {
	code   string
	errors diagnostics.Diagnostics
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
		// check that variables is not yet declared
		_, exists := w.env.GetVariable(node.Name.Value, false)

		if exists {
			w.addErrorf(
				node.Token,
				diagnostics.Fatal,
				"%v variable is already declared",
				node.Name.Value,
			)
		}

		ok := w.typesExists(node.Name.Type)

		if !ok {
			w.addErrorf(
				node.Token,
				diagnostics.Fatal,
				"type %v is not defined", typeString(node.Name.Type),
			)
		}

		w.env.SetVariable(
			node.Name.Value,
			environment.Object{Token: node.Token, Type: node.Name.Type},
		)
		w.expectType(node.Value, node.Token, node.Name.Type)

		return w.Walk(node.Value)

	case *ast.ConstStatement:
		_, exists := w.env.GetVariable(node.Name.Value, false)

		if exists {
			w.addErrorf(
				node.Token,
				diagnostics.Fatal,
				"%v constant is already declared",
				node.Name.Value,
			)
			return w.Walk(node.Value)
		}

		ok := w.typesExists(node.Name.Type)

		if !ok {
			w.addErrorf(
				node.Token,
				diagnostics.Fatal,
				"type %v is not defined", typeString(node.Name.Type),
			)
		}

		w.env.SetVariable(node.Name.Value, environment.Object{Token: node.Token})

		return w.Walk(node.Value)

	case *ast.ReturnStatement:
		return w.Walk(node.ReturnValue)

	case *ast.TypeStatement:
		_, exists := w.env.GetType(node.Name.Value)

		if exists {
			w.addErrorf(node.Token, diagnostics.Fatal, "type %v already defined", node.Name.Value)
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
		return node.Type, node

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
		return node.Type, node

	case *ast.IntegerLiteral:
		return node.Type, node

	case *ast.VectorLiteral:
		var ts []*ast.Type
		for _, s := range node.Value {
			t, _ := w.Walk(s)
			ts = append(ts, t...)
		}

		// check that all types are equal
		ok := w.allSameTypes(ts)

		if !ok {
			w.addErrorf(
				node.Token,
				diagnostics.Warn,
				"vector must contain all same types, got %v",
				typeString(ts),
			)
		}

		return ts, node

	case *ast.SquareRightLiteral:
		return types, node

	case *ast.StringLiteral:
		return node.Type, node

	case *ast.PrefixExpression:
		return w.Walk(node.Right)

	case *ast.For:
		w.Walk(node.Statement)
		w.env = w.env.Enclose()
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
		t, n := w.Walk(node.Left)
		if node.Right != nil {
			w.expectType(node.Right, token.Item{}, t)
			w.Walk(node.Right)
		}

		return t, n

	case *ast.IfExpression:
		w.Walk(node.Condition)

		w.env = w.env.Enclose()

		w.Walk(node.Consequence)

		if node.Alternative != nil {
			w.env = w.env.Enclose()
			w.Walk(node.Alternative)
			w.env = w.env.Open()
		}
		w.env = w.env.Open()

	case *ast.FunctionLiteral:
		w.env = w.env.Enclose()

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
		w.env = w.env.Open()

	case *ast.CallExpression:
		return w.Walk(node.Function)
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

func (w *Walker) addError(tok token.Item, s diagnostics.Severity, m string) {
	w.errors = append(w.errors, diagnostics.New(tok, m, s))
}

func (w *Walker) addErrorf(tok token.Item, s diagnostics.Severity, fm string, a ...interface{}) {
	str := fmt.Sprintf(fm, a...)
	w.errors = append(w.errors, diagnostics.New(tok, str, s))
}

func (w *Walker) expectType(node ast.Node, tok token.Item, actual []*ast.Type) {
	expected, _ := w.Walk(node)
	ok, expected, _ := w.allTypesMatch(actual, expected)

	if ok {
		return
	}

	w.addErrorf(
		tok,
		diagnostics.Fatal,
		"wrong types, got (%v), may also be (%v)",
		typeString(actual),
		typeString(expected),
	)
}
