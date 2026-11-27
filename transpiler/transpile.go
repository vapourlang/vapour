package transpiler

import (
	"strings"

	"github.com/vapourlang/vapour/ast"
	"github.com/vapourlang/vapour/environment"
)

type Transpiler struct {
	code []string
	env  *environment.Environment
	opts options
}

type options struct {
	inGeneric bool
	inDefault bool
}

func (t *Transpiler) Env() *environment.Environment {
	return t.env
}

func New() *Transpiler {
	env := environment.New()

	return &Transpiler{
		env: env,
	}
}

func (t *Transpiler) Transpile(node ast.Node) ast.Node {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return t.transpileProgram(node)

	case *ast.ExpressionStatement:
		if node.Expression != nil {
			return t.Transpile(node.Expression)
		}

	case *ast.LetStatement:
		t.env.SetVariable(
			node.Name,
			environment.Variable{
				Token: node.Token,
				Value: node.Type,
			},
		)
		if node.Value != nil {
			t.transpileLetStatement(node)
			t.Transpile(node.Value)
			t.addNewLine()
		}

	case *ast.NewLine:
		t.addNewLine()

	case *ast.Comma:
		t.trimNewLine()
		t.addCode(",")

	case *ast.ConstStatement:
		t.env.SetVariable(
			node.Name,
			environment.Variable{
				Token:   node.Token,
				Value:   node.Type,
				IsConst: true,
			},
		)
		if node.Value != nil {
			t.transpileConstStatement(node)
			t.Transpile(node.Value)
			t.addNewLine()
		}

	case *ast.ReturnStatement:
		t.addNewLine()
		t.addCode("return(")
		t.Transpile(node.ReturnValue)
		t.addCode(")")

	case *ast.DeferStatement:
		t.addCode("on.exit((")
		t.Transpile(node.Func)
		t.addCode(")())")

	case *ast.TypeStatement:
		t.env.SetType(
			environment.Type{
				Token:      node.Token,
				Type:       node.Type,
				Attributes: node.Attributes,
				Object:     node.Object,
				Name:       node.Name,
			},
		)

	case *ast.Null:
		t.addCode("NULL")

	case *ast.Keyword:
		t.addCode(node.Value)

	case *ast.CommentStatement:
		t.addCode(node.TokenLiteral())

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			t.Transpile(s)
		}

	case *ast.Attribute:
		t.addCode(node.Value)
		return node

	case *ast.Identifier:
		t.addCode(node.Value)
		return node

	case *ast.Boolean:
		t.addCode(strings.ToUpper(node.String()))

	case *ast.IntegerLiteral:
		t.addCode(node.Value)

	case *ast.FloatLiteral:
		t.addCode(node.Value)

	case *ast.VectorLiteral:
		t.addCode("c(")
		for i, s := range node.Value {
			t.Transpile(s)
			if i < len(node.Value)-1 {
				t.addCode(", ")
			}
		}
		t.addCode(")")

	case *ast.StringLiteral:
		t.addCode(node.Token.Value + node.Str + node.Token.Value)

	case *ast.PrefixExpression:
		t.addCode("(")
		t.addCode(node.Operator)
		t.Transpile(node.Right)
		t.trimNewLine()
		t.addCode(")")

	case *ast.For:
		t.addCode("for(")
		t.addCode(node.Name.Name)
		t.addCode(" in ")
		t.Transpile(node.Vector)
		t.addCode(") {")
		t.env = environment.Enclose(t.env, nil)
		t.Transpile(node.Value)
		t.addCode("}")
		t.addNewLine()
		t.env = environment.Open(t.env)

	case *ast.While:
		t.addCode("while(")
		t.Transpile(node.Statement)
		t.addCode(") {")
		t.env = environment.Enclose(t.env, nil)
		t.Transpile(node.Value)
		t.addCode("}")
		t.addNewLine()
		t.env = environment.Open(t.env)

	case *ast.InfixExpression:
		n := t.Transpile(node.Left)

		transpiled := false
		if node.Operator == "$" {
			switch n := n.(type) {
			case *ast.Identifier:
				v, exists := t.env.GetVariable(n.Value, true)

				if !exists {
					break
				}

				// the variable is typed with a custom type
				// e.g.: let s: st, we check whether st is a struct
				isStruct := false
				for _, ty := range v.Value {
					if ty.Name == "struct" {
						isStruct = true
						continue
					}

					typ, ok := t.env.GetType(ty.Package, ty.Name)
					if ok && typ.Object == "struct" {
						isStruct = true
					}
				}

				if !isStruct {
					break
				}

				t.popCode()
				transpiled = true
				t.addCode("attr(" + n.Value + ", \"")
				t.Transpile(node.Right)
				t.addCode("\")")
			}
		}

		if transpiled {
			break
		}

		t.trimNewLine()

		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Operator == ".." {
			node.Operator = ":"
		}

		if node.Operator == "<-" {
			node.Operator = "<<-"
		}

		if node.Operator != "+=" && node.Operator != "-=" {
			t.addCode(node.Operator)
		}

		if node.Operator == "+=" {
			t.addCode("=")
			t.Transpile(node.Left)
			t.addCode("+")
		}

		if node.Operator == "-=" {
			t.addCode("=")
			t.Transpile(node.Left)
			t.addCode("-")
		}

		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Right != nil {
			t.Transpile(node.Right)
			t.addNewLine()
		}

	case *ast.IndexExpression:
		t.Transpile(node.Left)
		t.trimNewLine()
		t.addCode(node.Open())
		for i, idx := range node.Indices {
			// nil is an empty slot, e.g.: x[1, ]
			if idx != nil {
				t.Transpile(idx)
				t.trimNewLine()
			}
			if i < len(node.Indices)-1 {
				t.addCode(", ")
			}
		}
		t.addCode(node.Close())

	case *ast.IfExpression:
		t.addCode("if(")
		t.Transpile(node.Condition)
		t.addCode("){")
		t.env = environment.Enclose(t.env, nil)
		t.Transpile(node.Consequence)
		t.env = environment.Open(t.env)
		t.addCode("}")

		if node.Alternative != nil {
			t.addCode(" else {")
			t.env = environment.Enclose(t.env, nil)
			t.Transpile(node.Alternative)
			t.env = environment.Open(t.env)
			t.addCode("}")
		}

		// the closing } must not be followed by
		// the next statement on the same line
		t.addNewLine()

	case *ast.FunctionLiteral:
		t.env = environment.Enclose(t.env, node.ReturnType)

		if node.Name != "" {
			t.addCode(node.Name)
		}

		if t.opts.inDefault {
			node.Method = &ast.Type{Name: "default"}
		}

		if node.Method != nil && node.Method.Name != "any" {
			t.addCode("." + node.Method.Name)
		}

		if node.Operator != "" {
			t.addCode(" " + node.Operator + " ")
		}

		t.addCode("function(")

		if node.MethodVariable != "" {
			t.addCode(node.MethodVariable)
			if len(node.Parameters) > 0 {
				t.addCode(", ")
			}
		}

		for i, p := range node.Parameters {
			t.env.SetVariable(
				p.Token.Value,
				environment.Variable{
					Token:   p.Token,
					Value:   p.Type,
					CanMiss: p.Default == nil,
					IsConst: false,
					Used:    true,
				},
			)

			t.addCode(p.Name)

			if p.Operator == "=" {
				t.addCode(" = ")
				t.Transpile(p.Default)
			}

			if i < len(node.Parameters)-1 {
				t.addCode(", ")
			}
		}

		t.addCode(") {")
		if node.Body != nil {
			t.Transpile(node.Body)
		}

		if node.Body == nil {
			t.addCode("UseMethod(\"" + node.Name + "\")")
		}

		t.env = environment.Open(t.env)
		t.addCode("}")
		t.addNewLine()

	case *ast.DecoratorEnvironment:
		t.env.SetEnv(
			node.Type.Name,
			environment.Env{
				Token: node.Token,
				Value: node,
			},
		)
		t.Transpile(node.Type)

	case *ast.DecoratorFactor:
		t.env.SetFactor(
			node.Type.Name,
			environment.Factor{
				Token: node.Token,
				Value: node,
			},
		)
		t.Transpile(node.Type)

	case *ast.DecoratorMatrix:
		t.env.SetMatrix(
			node.Type.Name,
			environment.Matrix{
				Token: node.Token,
				Value: node,
			},
		)
		t.Transpile(node.Type)

	case *ast.DecoratorClass:
		t.env.SetClass(
			node.Type.Name,
			environment.Class{
				Token: node.Token,
				Value: node,
			},
		)
		t.Transpile(node.Type)

	case *ast.DecoratorGeneric:
		t.opts.inGeneric = true
		t.Transpile(node.Func)
		t.opts.inGeneric = false

	case *ast.DecoratorDefault:
		t.opts.inDefault = true
		n := t.Transpile(node.Func)
		t.opts.inDefault = false
		return n

	case *ast.CallExpression:
		t.transpileCallExpression(node)
		//t.addNewLine() // avoid newline after return(x)
	}

	return node
}

func (t *Transpiler) transpileProgram(program *ast.Program) ast.Node {
	var node ast.Node

	for _, statement := range program.Statements {
		node := t.Transpile(statement)

		switch n := node.(type) {
		case *ast.ReturnStatement:
			t.addCode("")
			t.addCode("return(")
			if n.ReturnValue != nil {
				t.Transpile(n.ReturnValue)
			}
			t.addCode(")")
		}
	}

	return node
}

func (t *Transpiler) transpileCallExpression(node *ast.CallExpression) {
	typ, _ := t.env.GetType("", node.Name)

	if typ.Object == "struct" {
		t.transpileCallExpressionStruct(node, typ)
		return
	}

	if typ.Object == "object" {
		t.transpileCallExpressionObject(node, typ)
		return
	}

	if typ.Object == "environment" {
		t.transpileCallExpressionEnvironment(node, typ)
		return
	}

	if typ.Object == "dataframe" {
		t.transpileCallExpressionDataframe(node, typ)
		return
	}

	if typ.Object == "list" {
		t.transpileCallExpressionObject(node, typ)
		return
	}

	if typ.Object == "vector" {
		t.transpileCallExpressionVector(node, typ)
		return
	}

	if typ.Object == "matrix" {
		t.transpileCallExpressionMatrix(node, typ)
		return
	}

	if typ.Object == "factor" {
		t.transpileCallExpressionFactor(node, typ)
		return
	}

	t.addCode(node.Function + "(")
	for i, a := range node.Arguments {
		t.Transpile(a.Value)
		if i < len(node.Arguments)-1 {
			t.addCode(", ")
		}
	}
	t.addCode(")")
}

// structure(list2env(list(...), envir = new.env(<@environment args>)), class = ...)
func (t *Transpiler) transpileCallExpressionEnvironment(node *ast.CallExpression, typ environment.Type) {
	t.addCode("structure(list2env(list(")
	t.transpileArguments(node.Arguments)
	t.addCode("), envir = new.env(")

	env, exists := t.env.GetEnv(typ.Name)
	if exists {
		t.transpileArguments(env.Value.Arguments)
	}

	t.addCode("))")
	t.addClass(typ, "environment")
	t.addCode(")")
}

// structure(factor(..., <@factor args>), class = ...)
func (t *Transpiler) transpileCallExpressionFactor(node *ast.CallExpression, typ environment.Type) {
	t.addCode("structure(factor(")
	t.transpileArguments(node.Arguments)

	fct, exists := t.env.GetFactor(typ.Name)
	if exists {
		t.transpileExtraArguments(node.Arguments, fct.Value.Arguments)
	}

	t.addCode(")")
	t.addClass(typ, "factor")
	t.addCode(")")
}

// structure(matrix(..., <@matrix args>), class = ...)
func (t *Transpiler) transpileCallExpressionMatrix(node *ast.CallExpression, typ environment.Type) {
	t.addCode("structure(matrix(")
	t.transpileArguments(node.Arguments)

	mat, exists := t.env.GetMatrix(typ.Name)
	if exists {
		t.transpileExtraArguments(node.Arguments, mat.Value.Arguments)
	}

	t.addCode(")")
	t.addClass(typ, "matrix")
	t.addCode(")")
}

func (t *Transpiler) transpileCallExpressionVector(node *ast.CallExpression, typ environment.Type) {
	t.addCode("c(")
	t.transpileArguments(node.Arguments)
	t.addCode(")")
}

// structure(data.frame(...), names = c(...), class = ...)
func (t *Transpiler) transpileCallExpressionDataframe(node *ast.CallExpression, typ environment.Type) {
	names := []string{}
	for _, a := range node.Arguments {
		names = append(names, a.Name)
	}

	t.addCode("structure(data.frame(")
	t.transpileArguments(node.Arguments)
	t.addCode("), names = c(\"" + strings.Join(names, "\", \"") + "\")")
	t.addClass(typ, "data.frame")
	t.addCode(")")
}

// structure(list(...), class = ...)
func (t *Transpiler) transpileCallExpressionObject(node *ast.CallExpression, typ environment.Type) {
	t.addCode("structure(list(")
	t.transpileArguments(node.Arguments)
	t.addCode(")")
	t.addClass(typ, "list")
	t.addCode(")")
}

// structure(<value>, <attributes>, class = ...)
func (t *Transpiler) transpileCallExpressionStruct(node *ast.CallExpression, typ environment.Type) {
	t.addCode("structure(")
	t.transpileArguments(node.Arguments)
	t.addClass(typ)
	t.addCode(")")
}

// transpileArguments transpiles comma separated arguments.
func (t *Transpiler) transpileArguments(args []ast.Argument) {
	for i, a := range args {
		t.Transpile(a.Value)
		if i < len(args)-1 {
			t.addCode(", ")
		}
	}
}

// transpileExtraArguments appends decorator arguments
// after the arguments of the call (which may be empty).
func (t *Transpiler) transpileExtraArguments(callArgs, extra []ast.Argument) {
	if len(extra) == 0 {
		return
	}

	if len(callArgs) > 0 {
		t.addCode(", ")
	}

	for i, a := range extra {
		t.Transpile(a.Value)
		t.trimNewLine()
		if i < len(extra)-1 {
			t.addCode(", ")
		}
	}
}

// addClass adds the class argument to structure():
// the classes from @class if set, otherwise
// the type name followed by the base classes.
func (t *Transpiler) addClass(typ environment.Type, base ...string) {
	classes := append([]string{typ.Name}, base...)

	cl, exists := t.env.GetClass(typ.Name)
	if exists {
		classes = cl.Value.Classes
	}

	if len(classes) == 1 {
		t.addCode(", class=\"" + classes[0] + "\"")
		return
	}

	t.addCode(", class=c(\"" + strings.Join(classes, "\", \"") + "\")")
}

func (t *Transpiler) GetCode() string {
	return strings.Join(t.code, "")
}

func (t *Transpiler) addCode(code string) {
	t.code = append(t.code, code)
}

func (t *Transpiler) popCode() {
	t.code = t.code[:len(t.code)-1]
}

// trimNewLine removes a trailing new line, if any.
func (t *Transpiler) trimNewLine() {
	if len(t.code) > 0 && t.code[len(t.code)-1] == "\n" {
		t.popCode()
	}
}

func (t *Transpiler) transpileLetStatement(l *ast.LetStatement) {
	t.addCode(l.Name + " = ")
}

func (t *Transpiler) transpileConstStatement(c *ast.ConstStatement) {
	t.addCode(c.Name + " = ")
}

func (t *Transpiler) addNewLine() {
	if len(t.code) > 0 && t.code[len(t.code)-1] != "\n" {
		t.addCode("\n")
	}
}
