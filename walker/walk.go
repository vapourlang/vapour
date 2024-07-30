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
		env: environment.New(environment.Object{}),
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

		if len(node.Name.Type) > 0 {
			w.addErrorf(
				node.Token,
				diagnostics.Warn,
				"constants can only be of a single type, got: %v", typeString(node.Name.Type),
			)
		}

		w.env.SetVariable(node.Name.Value, environment.Object{Token: node.Token})

		w.expectType(node.Value, node.Token, node.Name.Type)

		return w.Walk(node.Value)

	case *ast.ReturnStatement:
		if node.ReturnValue == nil {
			return []*ast.Type{{Name: "null", List: false}}, node
		}

		t, n := w.Walk(node.ReturnValue)

		inFn, fn := w.env.GetFunctionEnvironment()

		if !inFn {
			return t, n
		}

		ok, _ := w.typesIn(fn.Type, t)

		if !ok {
			w.addErrorf(
				node.Token,
				diagnostics.Fatal,
				"function expects %v, return %v",
				typeString(fn.Type),
				typeString(t),
			)
		}

		return t, n

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
		fn, exists := w.env.GetFunction(node.Value, true)

		if exists {
			fmt.Printf("function: %v %v\n", node.Value, typeString(fn.Type))
			return fn.Type, node
		}

		v, exists := w.env.GetVariable(node.Value, true)

		if exists {
			fmt.Printf("identifier: %v %v\n", node.Value, typeString(v.Type))
			return v.Type, node
		}

		fmt.Printf("unknown identifier: %v %v\n", node.Value, typeString(node.Type))
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
		w.env = w.env.Enclose(w.env.Fn)
		t, n := w.Walk(node.Value)
		w.env = w.env.Open()
		return t, n

	case *ast.While:
		w.Walk(node.Statement)
		w.env = w.env.Enclose(w.env.Fn)
		t, n := w.Walk(node.Value)
		w.env = w.env.Open()
		return t, n

	case *ast.InfixExpression:
		tl, n := w.Walk(node.Left)
		if node.Right != nil {
			w.expectType(node.Right, node.Token, tl)
		}

		return tl, n

	case *ast.IfExpression:
		w.Walk(node.Condition)

		w.env = w.env.Enclose(w.env.Fn)
		w.Walk(node.Consequence)
		w.env = w.env.Open()

		if node.Alternative != nil {
			w.env = w.env.Enclose(w.env.Fn)
			w.Walk(node.Alternative)
			w.env = w.env.Open()
		}

	case *ast.FunctionLiteral:
		fmt.Println("FUNCTION")
		w.env = w.env.Enclose(
			environment.Object{
				Token: node.Token,
				Name:  node.Name.Value,
				Type:  node.Type,
			},
		)

		params := []string{}
		for _, p := range node.Parameters {
			w.env.SetVariable(
				p.Token.Value,
				environment.Object{
					Token: p.Token,
					Type:  p.Type,
				},
			)
			params = append(params, p.String())
		}

		for _, s := range node.Body.Statements {
			w.Walk(s)
		}

		w.env = w.env.Open()

		w.env.SetFunction(
			node.Name.Value,
			environment.Object{
				Token: node.Token,
				Type:  node.Type,
			},
		)

		return node.Type, node

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
