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
		fmt.Println("expression statement")
		if node.Expression != nil {
			w.Walk(node.Expression)
		}

	case *ast.Program:
		fmt.Println("program")
		return w.walkProgram(node)

	case *ast.LetStatement:
		fmt.Println("let")
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
		fmt.Println("const")
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
				diagnostics.Fatal,
				"constants can only be of a single type, got: %v", typeString(node.Name.Type),
			)
		}

		w.env.SetVariable(node.Name.Value, environment.Object{Token: node.Token})

		w.expectType(node.Value, node.Token, node.Name.Type)

		return w.Walk(node.Value)

	case *ast.ReturnStatement:
		fmt.Println("return")
		return w.Walk(node.ReturnValue)

	case *ast.TypeStatement:
		fmt.Println("type")
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
		fmt.Printf("keyword: %v\n", typeString(node.Type))
		return node.Type, node

	case *ast.Null:
		fmt.Println("null")
		return node.Type, node

	case *ast.CommentStatement:
		fmt.Println("comment")
		return types, node

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			w.Walk(s)
		}

	case *ast.Identifier:
		fn, exists := w.env.GetFunction(node.Value, true)

		if exists {
			fmt.Printf("identifier: %v %v\n", node.Value, typeString(fn.Type))
			return fn.Type, node
		}

		v, exists := w.env.GetVariable(node.Value, true)

		if exists {
			fmt.Printf("identifier: %v %v\n", node.Value, typeString(v.Type))
			return v.Type, node
		}

		fmt.Printf("identifier: %v %v\n", node.Value, typeString(node.Type))
		return node.Type, node

	case *ast.Boolean:
		fmt.Printf("bool %v\n", node.Value)
		return node.Type, node

	case *ast.IntegerLiteral:
		fmt.Printf("int %v\n", node.Value)
		return node.Type, node

	case *ast.VectorLiteral:
		fmt.Println("vector")
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
		fmt.Println("]")
		return types, node

	case *ast.StringLiteral:
		fmt.Println("string literal")
		return node.Type, node

	case *ast.PrefixExpression:
		fmt.Println("prefix expression")
		return w.Walk(node.Right)

	case *ast.For:
		fmt.Println("for")
		w.Walk(node.Statement)
		w.env = w.env.Enclose()
		t, n := w.Walk(node.Value)
		w.env = w.env.Open()
		return t, n

	case *ast.While:
		fmt.Println("while")
		w.Walk(node.Statement)
		w.env = w.env.Enclose()
		t, n := w.Walk(node.Value)
		w.env = w.env.Open()
		return t, n

	case *ast.InfixExpression:
		fmt.Println("infix")
		tl, n := w.Walk(node.Left)
		if node.Right != nil {
			w.expectType(node.Right, node.Token, tl)
		}

		return tl, n

	case *ast.IfExpression:
		fmt.Println("if")
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
		fmt.Println("function")
		w.env = w.env.Enclose()

		params := []string{}
		for _, p := range node.Parameters {
			w.env.SetVariable(
				p.TokenLiteral(),
				environment.Object{
					Token: p.Token,
					Type:  p.Type,
				},
			)
			params = append(params, p.String())
		}

		w.expectType(node.Body, node.Token, node.Type)
		w.Walk(node.Body)
		w.env = w.env.Open()

		w.env.SetFunction(
			node.Name.Value,
			environment.Object{
				Token: node.Token,
				Type:  node.Type,
			},
		)

	case *ast.CallExpression:
		fmt.Println("call")
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

// Expect a type, where expectation is the left-side node
// node is the right-side node to traverse
func (w *Walker) expectType(node ast.Node, tok token.Item, expectation []*ast.Type) {
	fmt.Printf("+ check type START (%v), left: %v\n", tok.Value, typeString(expectation))

	actual, _ := w.Walk(node)
	ok, in, missing := w.typesIn(expectation, actual)
	fmt.Printf("+ check type DONE, right: %v - matches: %v\n", typeString(actual), ok)

	if len(in) > 0 && tok.Class != token.ItemLet && tok.Class != token.ItemConst {
		w.addErrorf(
			tok,
			diagnostics.Info,
			"token `%v` expects unnecessary types: %v",
			tok.Value,
			typeString(in),
		)
	}

	if ok {
		return
	}

	w.addErrorf(
		tok,
		diagnostics.Fatal,
		"token `%v` type mismatch, assigning (%v) to (%v), missing (%v)",
		tok.Value,
		typeString(actual),
		typeString(expectation),
		typeString(missing),
	)
}
