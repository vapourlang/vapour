package walker

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/environment"
	"github.com/devOpifex/vapour/r"
)

type Walker struct {
	errors diagnostics.Diagnostics
	env    *environment.Environment
	state  *state
}

type state struct {
	inconst   bool
	inmissing bool
	incall    int
}

func New() *Walker {
	return &Walker{
		env:   environment.New(environment.Object{}),
		state: &state{},
	}
}

func (w *Walker) Run(node ast.Node) {
	w.Walk(node)
	w.warnUnusedVariables()
	w.warnUnusedTypes()
	w.warnUnusedFunctions()
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

	case *ast.Decorator:
		w.walkDecorator(node)

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

	case *ast.Attrbute:
		return types, node

	case *ast.Identifier:
		return w.walkIdentifier(node)

	case *ast.Boolean:
		return node.Type, node

	case *ast.IntegerLiteral:
		return node.Type, node

	case *ast.FloatLiteral:
		return node.Type, node

	case *ast.VectorLiteral:
		return w.walkVectorLiteral(node)

	case *ast.SquareRightLiteral:
		return types, node

	case *ast.StringLiteral:
		return node.Type, node

	case *ast.PrefixExpression:
		return w.Walk(node.Right)

	case *ast.For:
		w.env = w.env.Enclose(w.env.Fn)
		w.Walk(node.Name)
		w.Walk(node.Vector)
		t, n := w.Walk(node.Value)
		w.warnUnusedVariables()
		w.env = w.env.Open()
		return t, n

	case *ast.While:
		w.Walk(node.Statement)
		w.env = w.env.Enclose(w.env.Fn)
		t, n := w.Walk(node.Value)
		w.warnUnusedVariables()
		w.env = w.env.Open()
		return t, n

	case *ast.InfixExpression:
		return w.walkInfixExpression(node)

	case *ast.IfExpression:
		w.Walk(node.Condition)

		w.env = w.env.Enclose(w.env.Fn)
		w.Walk(node.Consequence)
		w.warnUnusedVariables()
		w.env = w.env.Open()

		if node.Alternative != nil {
			w.env = w.env.Enclose(w.env.Fn)
			w.Walk(node.Alternative)
			w.warnUnusedVariables()
			w.env = w.env.Open()
		}

	case *ast.FunctionLiteral:
		return w.walkFunctionLiteral(node)

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
	w.state.incall += 1

	defer func() {
		w.state.incall -= 1
	}()

	token := node.Function.Item()

	fn, fnExists := w.env.GetFunction(token.Value, true)

	if fnExists {
		w.env.SetFunctionUsed(token.Value)
	}

	ty, tyExists := w.env.GetType(token.Value)

	// we don't have the type or function
	// it's an R function not declared in Vapour
	// we also currently ignore base R functions
	if !fnExists && !tyExists || (fnExists && fn.Name != "") {
		t, n := w.Walk(node.Function)

		switch n := n.(type) {
		case *ast.Identifier:
			if n.Value == "missing" {
				w.state.inmissing = true
			}
		}

		for _, v := range node.Arguments {
			w.Walk(v.Value)
		}

		w.state.inmissing = false

		return t, n
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
					"type `%v` is an object, all arguments must be named",
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
					"call `%v()` with argument #%v of type `%v`, expected parameter of type `%v`",
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

	// we are not declaring a constant
	// we check that we it is not a constant
	if !w.state.inconst && node.Operator == "=" {
		switch n := ln.(type) {
		case *ast.Identifier:
			v, exists := w.env.GetVariable(n.Value, false)

			if exists && v.Const {
				w.addFatalf(n.Token, "`%v` is a constant", n.Value)
			}
		}
	}

	if !w.state.inconst && w.state.incall == 0 && node.Operator != "::" && node.Operator != "<-" {
		switch n := ln.(type) {
		case *ast.Identifier:
			_, exists := w.env.GetVariable(n.Value, true)

			if !exists {
				_, exists = w.env.GetFunction(n.Value, true)
			}

			if !exists {
				w.addFatalf(n.Token, "variable `%v` does not exist", n.Value)
			}
		}
	}

	if node.Operator == "<-" {
		_, exists := w.env.GetVariable(ln.Item().Value, true)

		if !exists {
			w.addWarnf(
				ln.Item(),
				"`%v` does not exist in parent environment(s)",
				ln.Item().Value,
			)
		}
	}

	if node.Right != nil {
		// mathematical and comparators operations return the node and type of the right hand
		if len(lt) > 0 && node.Operator != "$" && node.Operator != "[[" && node.Operator != "::" && node.Operator != "[" {
			w.expectType(node.Right, node.Token, lt)
			return w.Walk(node.Right)
		}

		if node.Operator == "::" {
			installed, _ := r.PackageIsInstalled(ln.Item().Value)

			if !installed {
				w.addHintf(
					ln.Item(),
					"package `%v` is not installed",
					ln.Item().Value,
				)
			}

			rt, rn := w.Walk(node.Right)

			ok, err := r.PackageHasFunction(ln.Item().Value, rn.Item().Value)

			if !ok && err == nil {
				w.addHintf(
					rn.Item(),
					"function `%v` is not exported by package `%v`",
					rn.Item().Value,
					ln.Item().Value,
				)
			}
			return rt, rn
		}

		// we need to check if the attributes exist in the type
		if node.Operator == "$" {
			if len(lt) == 0 {
				return lt, ln
			}

			rt, rn := w.Walk(node.Right)

			ts, exists := w.env.GetType(lt[0].Name)

			if !exists {
				return lt, ln
			}

			found := false
			for _, a := range ts.Attributes {
				if rn.Item().Value == a.Name.Value {
					found = true
				}
			}

			// skip the any type, it can be any types
			if !found && lt[0].Name != "any" {
				w.addFatalf(
					rn.Item(),
					"attribute `%v` does not exist on type `%v`",
					rn.Item().Value,
					lt[0].Name,
				)
			}

			return rt, rn
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

func (w *Walker) walkLetStatement(node *ast.LetStatement) ([]*ast.Type, ast.Node) {
	// check that variables is not yet declared
	_, exists := w.env.GetVariable(node.Name.Value, false)

	if exists {
		w.addFatalf(
			node.Name.Token,
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
		environment.Object{
			Token: node.Token,
			Type:  node.Name.Type,
			Name:  node.Name.Value,
		},
	)

	w.expectType(node.Value, node.Token, node.Name.Type)

	return w.Walk(node.Value)
}

func (w *Walker) walkConstStatement(node *ast.ConstStatement) ([]*ast.Type, ast.Node) {
	if node.Value == nil {
		w.addFatalf(
			node.Name.Token,
			"%v constant must be declared with a value",
			node.Name.Value,
		)
		return node.Name.Type, node
	}

	_, exists := w.env.GetVariable(node.Name.Value, false)

	if exists {
		w.addFatalf(
			node.Name.Token,
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

func (w *Walker) walkReturnStatement(node *ast.ReturnStatement) ([]*ast.Type, ast.Node) {
	if node.ReturnValue == nil {
		return []*ast.Type{{Name: "null", List: false}}, node
	}

	t, n := w.Walk(node.ReturnValue)

	inFn, fn := w.env.GetFunctionEnvironment()

	if !inFn {
		return t, n
	}

	ok, _ := w.validReturnTypes(fn, t)

	if !ok {
		w.addFatalf(
			node.Token,
			"function expects %v, return %v",
			typeString(fn.Type),
			typeString(t),
		)
	}

	return t, n
}

func (w *Walker) walkDecorator(node *ast.Decorator) {
	_, exists := w.env.GetType(node.Type.Name.Value)

	if exists {
		w.addFatalf(node.Type.Name.Token, "type %v already defined", node.Type.Name.Value)
	}

	w.env.SetType(
		node.Type.Name.Value,
		environment.Object{
			Token:      node.Token,
			Type:       node.Type.Type,
			Name:       node.Type.Name.Value,
			Attributes: node.Type.Attributes,
			List:       node.Type.List,
			Object:     node.Type.Name.Type,
		},
	)
}

func (w *Walker) walkTypeStatement(node *ast.TypeStatement) {
	_, exists := w.env.GetType(node.Name.Value)

	if exists {
		w.addFatalf(node.Name.Token, "type %v already defined", node.Name.Value)
	}

	w.env.SetType(
		node.Name.Value,
		environment.Object{
			Token:      node.Token,
			Type:       node.Type,
			Name:       node.Name.Value,
			Attributes: node.Attributes,
			List:       node.List,
			Object:     node.Name.Type,
			Used:       false,
		},
	)
}

func (w *Walker) walkIdentifier(node *ast.Identifier) ([]*ast.Type, ast.Node) {
	v, exists := w.env.GetVariable(node.Value, true)

	if exists {
		// this probably could be improved
		// this logic is too simplistic
		w.env.SetVariableUsed(node.Value)

		if w.state.inmissing {
			w.env.SetVariableNotMissing(node.Value)
		}

		if !w.state.inmissing && v.CanMiss {
			w.addHintf(
				node.Token,
				"`%v` might be missing",
				node.Token.Value,
			)
		}

		return v.Type, node
	}

	fn, exists := w.env.GetFunction(node.Value, true)

	if exists {
		return fn.Type, node
	}

	_, exists = w.env.GetType(node.Value)

	if exists {
		w.env.SetTypeUsed(node.Value)
	}

	return node.Type, node
}

func (w *Walker) walkVectorLiteral(node *ast.VectorLiteral) ([]*ast.Type, ast.Node) {
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
}

func (w *Walker) walkFunctionLiteral(node *ast.FunctionLiteral) ([]*ast.Type, ast.Node) {
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
		if p.Default != nil {
			w.Walk(p.Default)
		}

		paramsObject := environment.Object{
			Token:   p.Token,
			Type:    p.Type,
			Name:    p.Token.Value,
			CanMiss: p.Default == nil && !p.Method,
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

	// we only warn on unused variables
	// if the function is named
	// anonymous functions may have to not use variables
	if node.Name.Value != "" {
		w.warnUnusedVariables()
	}
	w.env = w.env.Open()

	if node.Name.Value != "" {

		_, exists := w.env.GetFunction(node.Name.Value, false)

		// we don't flag if it's a method
		if exists && node.Method == "" {
			w.addFatalf(
				node.Token,
				"function `%v` is already defined",
				node.Name.Value,
			)
		}

		if !exists {
			w.env.SetFunction(
				node.Name.Value,
				environment.Object{
					Name:       node.Name.Value,
					Token:      node.Token,
					Type:       node.Type,
					Parameters: params,
					Used:       false,
					Method:     node.Method,
				},
			)
		}
	}

	fn, found := w.env.GetTypeFromSignature(node)

	if found {
		return []*ast.Type{{Name: fn}}, node
	}

	return []*ast.Type{{Name: "fn"}}, node
}

func (w *Walker) warnUnusedVariables() {
	vars, ok := w.env.AllVariablesUsed()

	if ok {
		return
	}

	for _, v := range vars {
		w.addInfof(
			v.Token,
			"variable `%v` is never used",
			v.Name,
		)
	}
}

func (w *Walker) warnUnusedFunctions() {
	fns, ok := w.env.AllFunctionsUsed()

	if ok {
		return
	}

	for _, v := range fns {
		w.addInfof(
			v.Token,
			"function `%v` is never called",
			v.Name,
		)
	}
}

func (w *Walker) warnUnusedTypes() {
	types, ok := w.env.AllTypesUsed()

	if ok {
		return
	}

	for _, v := range types {
		w.addInfof(
			v.Token,
			"type `%v` is never used",
			v.Name,
		)
	}
}

func (w *Walker) walkSquare(node *ast.Square) ([]*ast.Type, ast.Node) {
	var types []*ast.Type
	var n ast.Node
	for _, s := range node.Statements {
		types, n = w.Walk(s)
	}

	return types, n
}
