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
}

func New() *Walker {
	return &Walker{
		env: environment.New(),
	}
}

func (w *Walker) Run(node ast.Node) {
	w.Walk(node)
	w.warnUnusedTypes()
	w.warnUnusedVariables()
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
		w.walkBlockStatement(node)

	case *ast.Attrbute:
		return types, node

	case *ast.Identifier:
		return w.walkIdentifier(node)

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
		w.walkFor(node)

	case *ast.While:
		w.Walk(node.Statement)
		w.env = environment.Enclose(w.env, nil)
		t, n := w.Walk(node.Value)
		w.env = environment.Open(w.env)
		return t, n

	case *ast.InfixExpression:
		return w.walkInfixExpression(node)

	case *ast.IfExpression:
		w.Walk(node.Condition)

		w.env = environment.Enclose(w.env, nil)
		w.Walk(node.Consequence)
		w.env = environment.Open(w.env)

		if node.Alternative != nil {
			w.env = environment.Enclose(w.env, nil)
			w.Walk(node.Alternative)
			w.env = environment.Open(w.env)
		}

	case *ast.FunctionLiteral:
		w.walkFunctionLiteral(node)

	case *ast.CallExpression:
		return w.walkCallExpression(node)
	}

	return types, node
}

func (w *Walker) walkProgram(program *ast.Program) (ast.Types, ast.Node) {
	var node ast.Node
	var types ast.Types

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

func (w *Walker) walkCallExpression(node *ast.CallExpression) (ast.Types, ast.Node) {
	fn, exists := w.env.GetFunction(node.Name, true)

	// we skip where there is no package, it's currently an indicator of external fn
	// we skip if it has elipsis, we can't check that
	if exists && fn.Package == "" {
		return w.walkKnownCallExpression(node, fn)
	}

	t, exists := w.env.GetType(node.Name)

	if exists {
		return w.walkKnownCallTypeExpression(node, t)
	}

	// we skip where there is no package, it's currently an indicator of external fn
	// we skip if it has elipsis, we can't check that
	if exists && fn.Package == "" {
		return w.walkKnownCallExpression(node, fn)
	}

	if node.Name == "missing" {
		return w.walkCallExpressionMissing(node)
	}

	for _, v := range node.Arguments {
		w.Walk(v.Value)
		w.checkIfIdentifier(v.Value)
	}

	return ast.Types{}, node
}

func (w *Walker) walkCallExpressionMissing(node *ast.CallExpression) (ast.Types, ast.Node) {
	for _, v := range node.Arguments {
		w.callIfIdentifier(v.Value, func(node *ast.Identifier) {
			w.env.SetVariableNotMissing(node.Value)
		})
		w.Walk(v.Value)
		w.checkIfIdentifier(v.Value)
	}

	return ast.Types{}, node
}

func (w *Walker) walkKnownCallTypeExpression(node *ast.CallExpression, t environment.Type) (ast.Types, ast.Node) {
	if t.Object == "object" {
		return w.walkKnownCallTypeObjectExpression(node, t)
	}

	if t.Object == "list" {
		return w.walkKnownCallTypeListExpression(node, t)
	}

	if t.Object == "struct" {
		return w.walkKnownCallTypeStructExpression(node, t)
	}

	if t.Object == "vector" {
		return w.walkKnownCallTypeVectorExpression(node, t)
	}

	if t.Object == "impliedList" {
		return w.walkKnownCallTypeImpliedListExpression(node, t)
	}

	for _, v := range node.Arguments {
		w.Walk(v.Value)
		w.checkIfIdentifier(v.Value)
	}

	return t.Type, node
}

func (w *Walker) walkKnownCallTypeImpliedListExpression(node *ast.CallExpression, t environment.Type) (ast.Types, ast.Node) {
	for _, v := range node.Arguments {
		rt, _ := w.Walk(v.Value)

		if v.Name != "" {
			w.addFatalf(
				v.Token,
				"`%v` expects unnamed arguments",
				t.Type[0].Name,
			)
			continue
		}

		if rt[0].Name != t.Type[0].Name {
			w.addFatalf(
				v.Token,
				"expects `%v`, got `%v`",
				t.Type[0].Name,
				rt[0].Name,
			)
		}
	}

	return t.Type, node
}

func (w *Walker) walkKnownCallTypeVectorExpression(node *ast.CallExpression, t environment.Type) (ast.Types, ast.Node) {
	for _, v := range node.Arguments {
		at, _ := w.Walk(v.Value)
		w.checkIfIdentifier(v.Value)
		ok := w.typesValid(t.Type, at)

		if !ok {
			w.addFatalf(
				node.Token,
				"`%v` expects `%v`, got `%v`",
				t.Name,
				t.Type,
				at,
			)
			continue
		}

		if v.Name != "" {
			w.addFatalf(
				v.Token,
				"vector expects unnamed arguments",
			)
			continue
		}
	}

	return t.Type, node
}

func (w *Walker) walkKnownCallTypeListExpression(node *ast.CallExpression, t environment.Type) (ast.Types, ast.Node) {
	for _, v := range node.Arguments {
		at, _ := w.Walk(v.Value)
		ok := w.typesValid(t.Type, at)

		if !ok {
			w.addFatalf(
				node.Token,
				"`%v` expects `%v`, got `%v`",
				t.Name,
				t.Type,
				at,
			)
			continue
		}

		if v.Name != "" {
			w.addFatalf(
				v.Token,
				"list expects unnamed arguments",
			)
			continue
		}
		w.checkIfIdentifier(v.Value)
	}

	return ast.Types{{Name: t.Name}}, node
}

func (w *Walker) walkKnownCallTypeStructExpression(node *ast.CallExpression, t environment.Type) (ast.Types, ast.Node) {
	for i, v := range node.Arguments {
		at, _ := w.Walk(v.Value)
		if i == 0 && v.Name != "" {
			w.addFatalf(
				v.Token,
				"struct expects unnamed first argument",
			)
			continue
		}

		if i > 0 && v.Name == "" {
			w.addFatalf(
				v.Token,
				"struct expects named arguments",
			)
			continue
		}

		if i == 0 {
			ok := w.typesValid(t.Type, at)
			if !ok {
				w.addFatalf(
					node.Token,
					"`%v` struct expects `%v`, got `%v`",
					t.Name,
					t.Type,
					at,
				)
				continue
			}
		}

		if i > 0 {
			w.attributeMatch(v, at, t)
		}
		w.checkIfIdentifier(v.Value)
	}

	return ast.Types{}, node
}

func (w *Walker) walkKnownCallTypeObjectExpression(node *ast.CallExpression, t environment.Type) (ast.Types, ast.Node) {
	for _, v := range node.Arguments {
		at, _ := w.Walk(v.Value)
		if v.Name == "" {
			w.addFatalf(
				v.Token,
				"object expects named arguments",
			)
			continue
		}
		w.checkIfIdentifier(v.Value)
		w.attributeMatch(v, at, t)
	}

	return ast.Types{{Name: t.Name}}, node
}

func (w *Walker) walkKnownCallExpression(node *ast.CallExpression, fn environment.Function) (ast.Types, ast.Node) {
	dots := hasElipsis(fn)
	for argumentIndex, argument := range node.Arguments {
		argumentType, _ := w.Walk(argument.Value)

		param, ok := getFunctionParameter(fn, argument.Name, argumentIndex)

		if !ok && argument.Name == "" && !dots {
			w.addFatalf(
				argument.Token,
				"could not find parameter #%v (too many arguments?)",
				argumentIndex+1,
			)
			continue
		}

		if !ok && argument.Name != "" && !dots {
			w.addFatalf(
				argument.Token,
				"could not find parameter `%v`",
				argument.Name,
			)
			continue
		}

		threedots := ""
		if !ok && dots {
			threedots = "(passed to ...)"
			param, _ = getFunctionElipsis(fn)
		}

		ok = w.typesValid(param.Type, argumentType)

		if !ok && argument.Name == "" {
			w.addFatalf(
				argument.Token,
				"argument #%v expects `%v`, got `%v` %v",
				argumentIndex+1,
				param.Type,
				argumentType,
				threedots,
			)
			continue
		}

		if !ok && argument.Name != "" {
			w.addFatalf(
				argument.Token,
				"argument `%v` expects `%v`, got `%v` %v",
				argument.Name,
				param.Type,
				argumentType,
				threedots,
			)
			continue
		}
	}

	return fn.Value.ReturnType, node
}

func hasElipsis(fn environment.Function) bool {
	for _, p := range fn.Value.Parameters {
		if p.Name == "..." {
			return true
		}
	}
	return false
}

func getFunctionParameter(fn environment.Function, name string, index int) (*ast.Parameter, bool) {
	for i, p := range fn.Value.Parameters {
		if p.Name == name {
			return p, true
		}

		if name == "" && i == index {
			return p, true
		}
	}

	return &ast.Parameter{}, false
}

func getFunctionElipsis(fn environment.Function) (*ast.Parameter, bool) {
	if !hasElipsis(fn) {
		return &ast.Parameter{}, false
	}

	for _, p := range fn.Value.Parameters {
		if p.Name == "..." {
			return p, true
		}
	}

	return &ast.Parameter{}, false
}

func (w *Walker) walkInfixExpression(node *ast.InfixExpression) ([]*ast.Type, ast.Node) {
	switch node.Operator {
	case "=":
		return w.walkInfixExpressionEqual(node)
	case "::":
		return w.walkInfixExpressionNamespace(node)
	case ":::":
		return w.walkInfixExpressionNamespaceInternal(node)
	case "+":
		return w.walkInfixExpressionMath(node)
	case "-":
		return w.walkInfixExpressionMath(node)
	case "/":
		return w.walkInfixExpressionMath(node)
	case "*":
		return w.walkInfixExpressionMath(node)
	case "+=":
		return w.walkInfixExpressionMath(node)
	case "-=":
		return w.walkInfixExpressionMath(node)
	case "<-":
		return w.walkInfixExpressionEqualParent(node)
	case "<":
		return w.walkInfixExpressionComparison(node)
	case ">":
		return w.walkInfixExpressionComparison(node)
	case "==":
		return w.walkInfixExpressionComparison(node)
	case "!=":
		return w.walkInfixExpressionComparison(node)
	case ">=":
		return w.walkInfixExpressionComparison(node)
	case "<=":
		return w.walkInfixExpressionComparison(node)
	case "|>":
		return w.walkInfixExpressionPipe(node)
	case "..":
		return w.walkInfixExpressionRange(node)
	case "$":
		return w.walkInfixExpressionDollar(node)
	default:
		return w.walkInfixExpressionDefault(node)
	}
}

func (w *Walker) walkFor(node *ast.For) {
	w.env = environment.Enclose(w.env, nil)
	w.Walk(node.Name)

	vectorType, vectorNode := w.Walk(node.Vector)
	ok := w.validIteratorTypes(vectorType)

	if !ok {
		w.addFatalf(
			vectorNode.Item(),
			"type `%v` is cannot be iterated",
			vectorType,
		)
	}

	w.walkBlockStatement(node.Value)
	w.env = environment.Open(w.env)
}

func (w *Walker) walkInfixExpressionDollar(node *ast.InfixExpression) (ast.Types, ast.Node) {
	_, ln := w.Walk(node.Left)

	w.checkIfIdentifier(ln)

	if node.Right == nil {
		w.addFatalf(
			node.Token,
			"expecting right",
		)
	}

	return w.Walk(node.Right)
}

func (w *Walker) walkInfixExpressionRange(node *ast.InfixExpression) (ast.Types, ast.Node) {
	lt, ln := w.Walk(node.Left)

	ok := w.validMathTypes(lt)
	if !ok {
		w.addFatalf(
			node.Token,
			"`%v`:`%v` is not valid",
			lt,
			lt,
		)
	}

	if node.Right != nil {
		rt, rn := w.Walk(node.Right)
		ok := w.validMathTypes(lt)
		if !ok {
			w.addFatalf(
				node.Token,
				"`%v`:`%v` is not valid",
				lt,
				lt,
			)
		}
		return rt, rn
	}

	return lt, ln
}

func (w *Walker) walkInfixExpressionPipe(node *ast.InfixExpression) (ast.Types, ast.Node) {
	w.Walk(node.Left)

	if node.Right == nil {
		w.addFatalf(
			node.Token,
			"pipe expect right-hand side",
		)
	}

	return w.Walk(node.Right)
}

func (w *Walker) walkInfixExpressionComparison(node *ast.InfixExpression) (ast.Types, ast.Node) {
	lt, ln := w.Walk(node.Left)

	w.checkIfIdentifier(ln)

	if node.Right != nil {
		rt, rn := w.Walk(node.Right)
		w.checkIfIdentifier(rn)

		ok := w.typesValid(lt, rt)
		if !ok {
			w.addInfof(
				node.Token,
				"comparison `%v` %v `%v` is not logical",
				lt,
				node.Operator,
				rt,
			)
		}
		return rt, rn
	}

	return lt, ln
}

func (w *Walker) walkInfixExpressionDefault(node *ast.InfixExpression) (ast.Types, ast.Node) {
	lt, ln := w.Walk(node.Left)

	if node.Right != nil {
		return w.Walk(node.Right)
	}

	return lt, ln
}

func (w *Walker) walkInfixExpressionMath(node *ast.InfixExpression) (ast.Types, ast.Node) {
	lt, ln := w.Walk(node.Left)

	ok := w.validMathTypes(lt)
	if !ok {
		w.addFatalf(
			node.Token,
			"`%v` %v `%v` is not valid",
			lt,
			node.Operator,
			lt,
		)
	}

	if node.Right != nil {
		rt, rn := w.Walk(node.Right)

		ok := w.validMathTypes(rt)
		if !ok {
			w.addFatalf(
				node.Token,
				"`%v` %v `%v` is not valid",
				lt,
				node.Operator,
				rt,
			)
		}
		return rt, rn
	}

	return lt, ln
}

func (w *Walker) walkInfixExpressionNS(node *ast.InfixExpression, operator string) (ast.Types, ast.Node) {
	_, ln := w.Walk(node.Left)

	exists, err := r.PackageIsInstalled(ln.Item().Value)

	if err != nil {
		w.addInfof(
			ln.Item(),
			"error checking if package `%v` is installed",
			ln.Item().Value,
		)
	}

	if !exists {
		w.addHintf(
			ln.Item(),
			"package `%v` is not installed",
			ln.Item().Value,
		)
	}

	if node.Right == nil {
		w.addFatalf(
			node.Token,
			"expects right hand side",
		)
	}

	rt, rn := w.Walk(node.Right)

	switch n := rn.(type) {
	case *ast.CallExpression:
		exists, err = r.PackageHasFunction(ln.Item().Value, operator, n.Function)

		if err != nil {
			w.addInfof(
				ln.Item(),
				"error checking `%v%v%v`",
				ln.Item().Value,
				operator,
				n.Function,
			)
		}

		if !exists {
			w.addHintf(
				ln.Item(),
				"`%v%v%v` not found",
				ln.Item().Value,
				operator,
				n.Function,
			)
		}
	}

	return rt, rn
}

func (w *Walker) walkInfixExpressionNamespace(node *ast.InfixExpression) (ast.Types, ast.Node) {
	return w.walkInfixExpressionNS(node, "::")
}

func (w *Walker) walkInfixExpressionNamespaceInternal(node *ast.InfixExpression) (ast.Types, ast.Node) {
	return w.walkInfixExpressionNS(node, ":::")
}

func (w *Walker) walkInfixExpressionEqual(node *ast.InfixExpression) (ast.Types, ast.Node) {
	lt, _ := w.Walk(node.Left)

	if node.Right == nil {
		w.addFatalf(
			node.Token,
			"expecting right hand side",
		)
	}

	rt, rn := w.Walk(node.Right)
	ok := w.typesValid(lt, rt)
	if !ok {
		w.addFatalf(
			node.Token,
			"left expects `%v`, right returns `%v`",
			lt,
			rt,
		)
	}
	return rt, rn
}

func (w *Walker) walkInfixExpressionEqualParent(node *ast.InfixExpression) ([]*ast.Type, ast.Node) {
	lt, ln := w.Walk(node.Left)

	w.callIfIdentifier(node.Left, func(n *ast.Identifier) {
		_, ok := w.env.GetVariableParent(n.Value)

		if ok {
			return
		}

		w.addFatalf(
			n.Token,
			"`%v` does not exist in parent environment",
			n.Value,
		)
	})

	if node.Right != nil {
		w.Walk(node.Right)
	}

	return lt, ln
}

func (w *Walker) walkLetStatement(node *ast.LetStatement) (ast.Types, ast.Node) {
	_, ok := w.env.GetVariable(node.Name, false)

	if ok {
		w.addFatalf(
			node.Token,
			"variable `%v` is already declared",
			node.Name,
		)

		return w.Walk(node.Value)
	}

	w.env.SetVariable(
		node.Name,
		environment.Variable{
			Token: node.Token,
			Value: node.Type,
			Name:  node.Name,
		},
	)

	rt, rn := w.Walk(node.Value)
	ok = w.typesValid(node.Type, rt)

	if !ok {
		w.addFatalf(
			node.Token,
			"`%v` expects `%v`, got `%v`",
			node.Name,
			node.Type,
			rt,
		)
	}

	return rt, rn
}

func (w *Walker) walkConstStatement(node *ast.ConstStatement) (ast.Types, ast.Node) {
	_, ok := w.env.GetVariable(node.Name, false)

	if ok {
		w.addFatalf(
			node.Token,
			"variable `%v` is already declared",
			node.Name,
		)

		return w.Walk(node.Value)
	}

	if len(node.Type) > 1 {
		w.addFatalf(
			node.Token,
			"constants may only have a single type",
		)
	}

	w.env.SetVariable(
		node.Name,
		environment.Variable{
			Token:   node.Token,
			Value:   node.Type,
			Name:    node.Name,
			IsConst: true,
		},
	)

	if node.Value == nil {
		w.addFatalf(
			node.Token,
			"constants without value",
		)
	}

	return w.Walk(node.Value)
}

func (w *Walker) walkReturnStatement(node *ast.ReturnStatement) (ast.Types, ast.Node) {
	t, n := w.Walk(node.ReturnValue)

	if w.env.ReturnType() != nil {
		ok := w.typesValid(w.env.ReturnType(), t)
		if !ok {
			w.addFatalf(
				node.Token,
				"return expects `%v`, got `%v`",
				w.env.ReturnType(),
				t,
			)
		}
	}

	return t, n
}

func (w *Walker) walkDecoratorDefault(node *ast.DecoratorDefault) (ast.Types, ast.Node) {
	if node.Func == nil {
		w.addFatalf(
			node.Token,
			"expecting function",
		)
	}
	return w.Walk(node.Func)
}

func (w *Walker) walkDecoratorGeneric(node *ast.DecoratorGeneric) {
	w.Walk(node.Func)
}

func (w *Walker) walkDecoratorClass(node *ast.DecoratorClass) (ast.Types, ast.Node) {
	w.env.SetClass(
		node.Type.Name,
		environment.Class{
			Token: node.Token,
			Value: node,
		},
	)

	return w.Walk(node.Type)
}

func (w *Walker) walkTypeStatement(node *ast.TypeStatement) {
	_, exists := w.env.GetType(node.Name)

	if exists {
		w.addFatalf(
			node.Token,
			"type `%v` already defined",
			node.Name,
		)
	}

	w.env.SetType(
		node.Name,
		environment.Type{
			Token:      node.Token,
			Type:       node.Type,
			Attributes: node.Attributes,
			Object:     node.Object,
			Name:       node.Name,
		},
	)
}

func (w *Walker) walkIdentifier(node *ast.Identifier) (ast.Types, ast.Node) {
	v, exists := w.env.GetVariable(node.Value, true)

	if exists {
		w.env.SetVariableUsed(node.Value)
		if v.CanMiss {
			w.addHintf(
				node.Token,
				"`%v` might be missing",
				v.Name,
			)
		}

		return v.Value, node
	}

	t, exists := w.env.GetType(node.Value)

	if exists {
		return t.Type, node
	}

	if !exists {
		w.checkIdentifier(node)
	}

	return node.Type, node
}

func (w *Walker) walkVectorLiteral(node *ast.VectorLiteral) ([]*ast.Type, ast.Node) {
	var ts ast.Types
	for _, s := range node.Value {
		t, _ := w.Walk(s)
		ts = append(ts, t...)
	}

	ok := w.allTypesIdentical(ts)

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
	if node.Name == "" {
		w.walkAnonymousFunctionLiteral(node)
		return
	}

	w.walkNamedFunctionLiteral(node)
}

func (w *Walker) walkNamedFunctionLiteral(node *ast.FunctionLiteral) {
	_, exists := w.env.GetFunction(node.Name, false)

	// we don't flag if it's a method
	if exists && node.Method == nil {
		w.addFatalf(
			node.NameToken,
			"function `%v` is already defined",
			node.Name,
		)
		return
	}

	_, exists = w.env.GetType(node.Name)

	// we don't flag if it's a method
	if exists && node.Method == nil {
		w.addFatalf(
			node.NameToken,
			"functions and types cannot share name (`%v`)",
			node.Name,
		)
		return
	}

	w.env.SetFunction(node.Name, environment.Function{Token: node.Token, Value: node})

	w.env = environment.Enclose(w.env, node.ReturnType)

	// we set the parameters in the environment
	// and check that we don't have duplicates
	paramsMap := make(map[string]bool)
	if node.Method != nil {
		paramsMap[node.MethodVariable] = true
		w.env.SetVariable(
			node.MethodVariable,
			environment.Variable{
				Token: node.Token,
				Value: ast.Types{node.Method},
				Name:  node.MethodVariable,
				Used:  true,
			},
		)
	}

	for _, p := range node.Parameters {
		if p.Default != nil {
			w.Walk(p.Default)
		}

		w.env.SetVariable(
			p.Token.Value,
			environment.Variable{
				Token:   p.Token,
				Value:   p.Type,
				CanMiss: p.Default == nil || p.Name == "...",
				Name:    p.Name,
				Used:    p.Name == "...",
			},
		)

		if p.Token.Value == "..." {
			continue
		}

		_, exists := paramsMap[p.Token.Value]

		if exists {
			w.addFatalf(p.Token, "duplicated function parameter `%v`", p.Token.Value)
		}

		paramsMap[p.Token.Value] = true
	}

	if node.Body != nil {
		for _, s := range node.Body.Statements {
			w.Walk(s)
		}
	}

	w.warnUnusedVariables()
	w.env = environment.Open(w.env)
}

func (w *Walker) walkAnonymousFunctionLiteral(node *ast.FunctionLiteral) {
	w.env = environment.Enclose(w.env, node.ReturnType)

	// we set the parameters in the environment
	// and check that we don't have duplicates
	paramsMap := make(map[string]bool)
	for _, p := range node.Parameters {
		if p.Default != nil {
			w.Walk(p.Default)
		}

		paramsObject := environment.Variable{
			Token:   p.Token,
			Value:   p.Type,
			CanMiss: p.Default == nil && p.Method,
			Name:    p.Token.Value,
			IsConst: false,
			Used:    false,
		}

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

	if node.Body != nil {
		for _, s := range node.Body.Statements {
			w.Walk(s)
		}
	}

	w.warnUnusedVariables()
	w.env = environment.Open(w.env)
}

func (w *Walker) walkSquare(node *ast.Square) (ast.Types, ast.Node) {
	var types []*ast.Type
	return types, node
}

func (w *Walker) walkBlockStatement(node *ast.BlockStatement) {
	for _, s := range node.Statements {
		w.Walk(s)
	}
}
