package walker

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/environment"
)

type Walker struct {
	errors diagnostics.Diagnostics
	env    *environment.Environment
	state  *state
}

type state struct {
	inconst bool
}

func New() *Walker {
	return &Walker{
		env:   environment.New(environment.Object{}),
		state: &state{},
	}
}

func (w *Walker) Walk(node ast.Node) ([]*ast.Type, ast.Node) {
	var types []*ast.Type

	switch node := node.(type) {

	case *ast.Program:
		return w.walkProgram(node)

	case *ast.ExpressionStatement:
		if node.Expression != nil {
			return w.Walk(node.Expression)
		}

	case *ast.LetStatement:
    return w.walkLetStatement(types, node)

	case *ast.ConstStatement:
    return w.walkConstStatement(types, node)

	case *ast.ReturnStatement:
		if node.ReturnValue == nil {
			return []*ast.Type{{Name: "null", List: false}}, node
		}

		t, n := w.Walk(node.ReturnValue)

		inFn, fn := w.env.GetFunctionEnvironment()

		if !inFn {
			return t, n
		}

		ok, _ := w.validTypes(fn.Type, t)

		if !ok {
			w.addFatalf(
				node.Token,
				"function expects %v, return %v",
				typeString(fn.Type),
				typeString(t),
			)
		}

		return t, n

	case *ast.TypeStatement:
		_, exists := w.env.GetType(node.Name.Value, node.List)

		if exists {
			w.addFatalf(node.Token, "type %v already defined", node.Name.Value)
		}

		w.env.SetType(
			node.Name.Value,
			environment.Object{
				Token:      node.Token,
				Type:       node.Type,
				Name:       node.Name.Value,
				Attributes: node.Attributes,
				List:       node.List,
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
			return fn.Type, node
		}

		v, exists := w.env.GetVariable(node.Value, true)

		if exists {
			return v.Type, node
		}

		return node.Type, node

	case *ast.Boolean:
		return node.Type, node

	case *ast.IntegerLiteral:
		return node.Type, node

	case *ast.FloatLiteral:
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
			w.addWarnf(
				node.Token,
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
		w.env = w.env.Enclose(w.env.Fn)
		w.Walk(node.Statement)
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
		return w.walkInfixExpression(node)

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
		w.env = w.env.Enclose(
			environment.Object{
				Token: node.Token,
				Name:  node.Name.Value,
				Type:  node.Type,
			},
		)

		var params []environment.Object
		paramsMap := make(map[string]bool)
		for _, p := range node.Parameters {
			paramsObject := environment.Object{
				Token: p.Token,
				Type:  p.Type,
				Name:  p.Token.Value,
			}

			params = append(params, paramsObject)

			w.env.SetVariable(
				p.Token.Value,
				paramsObject,
			)

			_, exists := paramsMap[p.Token.Value]

			if exists {
				w.addFatalf(p.Token, "duplicated function parameter `%v`", p.Token.Value)
			}
			paramsMap[p.Token.Value] = true
		}

		for _, s := range node.Body.Statements {
			w.Walk(s)
		}

		w.env = w.env.Open()

		if node.Name.Value != "" {
			w.env.SetFunction(
				node.Name.Value,
				environment.Object{
					Token:      node.Token,
					Type:       node.Type,
					Parameters: params,
				},
			)
		}

		return node.Type, node

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
	token := node.Function.Item()

	fn, fnExists := w.env.GetFunction(token.Value, true)
	ty, tyExists := w.env.GetType(token.Value, false)

	if !tyExists {
		ty, tyExists = w.env.GetType(token.Value, true)
	}

	if !fnExists && !tyExists {
		for _, v := range node.Arguments {
			w.Walk(v.Value)
		}
		return w.Walk(node.Function)
	}

	if tyExists && ty.Type[0].Name == "list" {
		for _, v := range node.Arguments {
			w.Walk(v.Value)
		}
		return w.Walk(node.Function)
	}

	// handle if it's a type too
	if tyExists {
		for argIndex, arg := range node.Arguments {
			argType, _ := w.Walk(arg.Value)

			if arg.Name == "" && ty.Type[0].Name == "object" {
				w.addFatalf(
					arg.Token,
					"type `%v` is an object, all arguments cannot be named",
					token.Value,
				)
				continue
			}

			if arg.Name != "" && ty.Type[0].Name == "struct" && argIndex == 0 {
				w.addFatalf(
					arg.Token,
					"type `%v` is a struct, first argument cannot be named",
					token.Value,
				)
				continue
			}

			if arg.Name == "" && ty.Type[0].Name == "struct" && argIndex > 0 {
				w.addFatalf(
					arg.Token,
					"type `%v` is a struct, attributes must be named",
					token.Value,
				)
				continue
			}

			found := false
			for _, a := range ty.Attributes {
				if arg.Name == "" {
					continue
				}

				if arg.Name != a.Name.Value {
					continue
				}

				found = true
				ok, _ := w.validTypes(a.Type, argType)

				if !ok {
					w.addFatalf(
						arg.Token,
						"call type `%v` argument `%v` of type `%v`, expected `%v`",
						token.Value,
						arg.Name,
						typeString(argType),
						typeString(a.Type),
					)
				}
			}

			if !found && arg.Name != "" {
				w.addWarnf(
					arg.Token,
					"call type `%v` argument `%v` is not an attribute",
					token.Value,
					arg.Name,
				)
			}
		}
		return w.Walk(node.Function)
	}

	hasElipsis := hasElipsis(fn.Parameters)
	elipsisType := getElipsisType(fn.Parameters)

	for argIndex, arg := range node.Arguments {
		argType, _ := w.Walk(arg.Value)

		// the function accepts ...
		// we just check the type
		if hasElipsis && elipsisType != nil {
			ok, _ := w.validTypes(elipsisType, argType)

			if !ok {
				w.addWarnf(
					arg.Token,
					"call of `%v` with argument #%v of type `%v`, expected parameter of type `%v`",
					token.Value,
					argIndex+1,
					typeString(argType),
					typeString(elipsisType),
				)
			}
			continue
		}

		found := false
		for pIndex, p := range fn.Parameters {
			if arg.Name != p.Name && arg.Name != "" {
				continue
			}

			if arg.Name == "" && argIndex != pIndex {
				continue
			}

			if arg.Name == "" && argIndex == pIndex {
				found = true
				ok, _ := w.validTypes(p.Type, argType)

				if !ok {
					w.addWarnf(
						arg.Token,
						"`%v()` argument #%v got type `%v`, expected `%v`",
						token.Value,
						argIndex+1,
						typeString(argType),
						typeString(p.Type),
					)
				}
				continue
			}

			found = true
			ok, _ := w.validTypes(p.Type, argType)

			if !ok {
				w.addWarnf(
					arg.Token,
					"`%v()` argument `%v` got type `%v`, expected `%v`",
					token.Value,
					arg.Name,
					typeString(argType),
					typeString(p.Type),
				)
			}
		}

		if !found && arg.Name == "" {
			w.addWarnf(
				arg.Token,
				"call `%v()` argument #%v is not a parameter",
				token.Value,
				argIndex+1,
			)
		}

		if !found && arg.Name != "" {
			w.addWarnf(
				arg.Token,
				"call `%v()` argument `%v` is not a parameter",
				token.Value,
				arg.Name,
			)
		}
	}

	return w.Walk(node.Function)
}

func (w *Walker) walkInfixExpression(node *ast.InfixExpression) ([]*ast.Type, ast.Node) {
	lt, ln := w.Walk(node.Left)

	if !w.state.inconst && node.Operator == "=" {
		switch n := ln.(type) {
		case *ast.Identifier:
			v, exists := w.env.GetVariable(n.Value, false)

			if exists && v.Const {
				w.addFatalf(n.Token, "`%v` is a constant", n.Value)
			}
		}
	}

	if node.Right != nil {
		if len(lt) > 0 && node.Operator != "$" && node.Operator != "[[" && node.Operator != "::" && node.Operator != "[" {
			w.expectType(node.Right, node.Token, lt)
			return w.Walk(node.Right)
		}

		if node.Operator == "$" {
			if len(lt) == 0 {
				return lt, ln
			}

			_, rn := w.Walk(node.Right)

			ts, exists := w.env.GetType(lt[0].Name, lt[0].List)

			if !exists {
				return lt, ln
			}

			found := false
			for _, a := range ts.Attributes {
				if rn.Item().Value == a.Name.Value {
					found = true
				}
			}

			if !found {
				w.addFatalf(
					rn.Item(),
					"attribute `%v` does not exist on type `%v`",
					rn.Item().Value,
					lt[0].Name,
				)
			}

			return lt, ln
		}

		rt, rn := w.Walk(node.Right)

		if node.Operator == "=" {
			return rt, rn
		}
	}

	return lt, ln
}

func hasElipsis(params []environment.Object) bool {
	for _, p := range params {
		if p.Name == "..." {
			return true
		}
	}

	return false
}

func getElipsisType(params []environment.Object) []*ast.Type {
	for _, p := range params {
		if p.Name == "..." {
			return p.Type
		}
	}

	return nil
}

func (w *Walker) walkLetStatement(types []*ast.Type, node *ast.Node) ([]*ast.Type, ast.Node) {
		// check that variables is not yet declared
		_, exists := w.env.GetVariable(node.Name.Value, false)

		if exists {
			w.addFatalf(
				node.Token,
				"%v variable is already declared",
				node.Name.Value,
			)
		}

		ok := w.typesExists(node.Name.Type)

		if !ok {
			w.addFatalf(
				node.Token,
				"type %v is not defined", typeString(node.Name.Type),
			)
		}

		w.env.SetVariable(
			node.Name.Value,
			environment.Object{Token: node.Token, Type: node.Name.Type},
		)

		w.expectType(node.Value, node.Token, node.Name.Type)

		return w.Walk(node.Value)
}

func (w *Walker) walkConstStatement(types []*ast.Type, node *ast.Node) ([]*ast.Type, ast.Node) {
  if node.Value == nil {
    w.addFatalf(
      node.Token,
      "%v constant must be declared with a value",
      node.Name.Value,
    )
    return node.Name.Type, node
  }

  _, exists := w.env.GetVariable(node.Name.Value, false)

  if exists {
    w.addFatalf(
      node.Token,
      "%v constant is already declared",
      node.Name.Value,
    )
    return w.Walk(node.Value)
  }

  ok := w.typesExists(node.Name.Type)

  if !ok {
    w.addFatalf(
      node.Token,
      "type %v is not defined", typeString(node.Name.Type),
    )
  }

  if len(node.Name.Type) > 1 {
    w.addWarnf(
      node.Token,
      "constants can only be of a single type, got: %v", typeString(node.Name.Type),
    )
  }

  w.env.SetVariable(node.Name.Value, environment.Object{Token: node.Token, Const: true})

  w.state.inconst = true
  w.expectType(node.Value, node.Token, node.Name.Type)

  t, n := w.Walk(node.Value)
  w.state.inconst = false
  return t, n
}
