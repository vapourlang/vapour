package transpiler

import (
	"strings"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/environment"
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
			t.addCode("\n")
		}

	case *ast.NewLine:
		t.addCode("\n")

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
			t.addCode("\n")
		}

	case *ast.ReturnStatement:
		t.addCode("\nreturn(")
		t.Transpile(node.ReturnValue)
		t.addCode(")")

	case *ast.DeferStatement:
		t.addCode("on.exit((")
		t.Transpile(node.Func)
		t.addCode(")())")

	case *ast.TypeStatement:
		t.env.SetType(
			node.Name,
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

	case *ast.Attrbute:
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
		t.env = environment.Open(t.env)

	case *ast.While:
		t.addCode("while(")
		t.Transpile(node.Statement)
		t.addCode(") {")
		t.env = environment.Enclose(t.env, nil)
		t.Transpile(node.Value)
		t.addCode("}")
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

				isStruct := false
				for _, ty := range v.Value {
					if ty.Name == "struct" {
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

		if t.code[len(t.code)-1] == "\n" {
			t.popCode()
		}

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
			t.addCode("\n")
		}

	case *ast.Square:
		t.popCode()
		if t.code[len(t.code)-1] == "\n" {
			t.popCode()
		}
		t.addCode(node.Token.Value)

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
				t.addCode(",")
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

	case *ast.DecoratorClass:
		t.Transpile(node.Type)
		t.env.SetClass(
			node.Type.Name,
			environment.Class{
				Token: node.Token,
				Value: node,
			},
		)

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
		t.addCode("\n")
	}

	return node
}

func (t *Transpiler) transpileProgram(program *ast.Program) ast.Node {
	var node ast.Node

	for _, statement := range program.Statements {
		node := t.Transpile(statement)

		switch n := node.(type) {
		case *ast.ReturnStatement:
			t.addCode("\nreturn(")
			if n.ReturnValue != nil {
				t.Transpile(n.ReturnValue)
			}
			t.addCode(")")
		}
	}

	return node
}

func (t *Transpiler) transpileCallExpression(node *ast.CallExpression) {
	typ, _ := t.env.GetType(node.Name)

	if typ.Object == "struct" {
		t.transpileCallExpressionStruct(node, typ)
		return
	}

	if typ.Object == "object" {
		t.transpileCallExpressionObject(node, typ)
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

	t.addCode(node.Function + "(")
	for i, a := range node.Arguments {
		t.Transpile(a.Value)
		if i < len(node.Arguments)-1 {
			t.addCode(", ")
		}
	}
	t.addCode(")")
}

func (t *Transpiler) transpileCallExpressionVector(node *ast.CallExpression, typ environment.Type) {
	t.addCode("c(")
	for i, a := range node.Arguments {
		t.Transpile(a.Value)
		if i < len(node.Arguments)-1 {
			t.addCode(", ")
		}
	}
	t.addCode(")")
}

func (t *Transpiler) transpileCallExpressionDataframe(node *ast.CallExpression, typ environment.Type) {
	names := []string{}
	t.addCode("structure(data.frame(")
	for i, a := range node.Arguments {
		t.Transpile(a.Value)
		names = append(names, a.Name)
		if i < len(node.Arguments)-1 {
			t.addCode(", ")
		}
	}
	t.addCode("), names = c(\"" + strings.Join(names, "\", \"") + "\")")

	cl, exists := t.env.GetClass(typ.Name)

	if exists {
		t.addCode(", class=c(\"" + strings.Join(cl.Value.Classes, "\", \"") + "\")")
		return
	}

	t.addCode(", class=c(\"" + typ.Name + "\", \"data.frame\")")

	t.addCode(")")
}

func (t *Transpiler) transpileCallExpressionObject(node *ast.CallExpression, typ environment.Type) {
	t.addCode("structure(list(")
	for i, a := range node.Arguments {
		t.Transpile(a.Value)
		if i < len(node.Arguments)-1 {
			t.addCode(", ")
		}
	}
	t.addCode(")")

	cl, exists := t.env.GetClass(typ.Name)

	if exists {
		t.addCode(", class=c(\"" + strings.Join(cl.Value.Classes, "\", \"") + "\")")
		return
	}

	t.addCode(", class=c(\"" + typ.Name + "\", \"list\")")

	t.addCode(")")
}

func (t *Transpiler) transpileCallExpressionStruct(node *ast.CallExpression, typ environment.Type) {
	t.addCode("structure(")
	for i, a := range node.Arguments {
		t.Transpile(a.Value)
		if i < len(node.Arguments)-1 {
			t.addCode(", ")
		}
	}
	cl, exists := t.env.GetClass(typ.Name)

	if exists {
		t.addCode(", class=c(\"" + strings.Join(cl.Value.Classes, "\", \"") + "\")")
		return
	}

	t.addCode(", class=\"" + typ.Name + "\"")

	t.addCode(")")
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

func (t *Transpiler) transpileLetStatement(l *ast.LetStatement) {
	t.addCode(l.Name + " = ")
}

func (t *Transpiler) transpileConstStatement(c *ast.ConstStatement) {
	t.addCode(c.Name + " = ")
}
