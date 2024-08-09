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
	inType    bool
	inClass   []string
	typeClass []string
}

func New() *Transpiler {
	env := environment.New(environment.Object{})

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
			node.Name.Value,
			environment.Object{Token: node.Token, Type: node.Name.Type},
		)
		if node.Value != nil {
			t.transpileLetStatement(node)
			t.Transpile(node.Value)
		}

	case *ast.NewLine:
		t.addCode("\n")

	case *ast.ConstStatement:
		t.env.SetVariable(
			node.Name.Value,
			environment.Object{Token: node.Token, Type: node.Name.Type},
		)
		if node.Value != nil {
			t.transpileConstStatement(node)
			t.Transpile(node.Value)
		}

	case *ast.ReturnStatement:
		t.addCode("\nreturn(")
		t.Transpile(node.ReturnValue)
		t.addCode(")")

	case *ast.TypeStatement:
		_, exists := t.env.GetType(node.Name.Value, node.List)

		if !exists {
			t.env.SetType(
				node.Name.Value,
				environment.Object{
					Token:      node.Token,
					Type:       node.Type,
					Name:       node.Name.Value,
					List:       node.List,
					Attributes: node.Attributes,
				},
			)
		}

		return node.Name

	case *ast.Null:
		t.addCode("NULL")

	case *ast.Keyword:
		t.addCode(node.Value)

	case *ast.CommentStatement:
		t.addCode(node.TokenLiteral())

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			t.Transpile(s)
			t.addCode("\n")
		}

	case *ast.Identifier:
		t.addCode(node.Value)
		return node

	case *ast.Boolean:
		t.addCode(node.String())

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
		t.Transpile(node.Statement)
		t.addCode(") {")
		t.env = t.env.Enclose(t.env.Fn)
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = t.env.Open()

	case *ast.While:
		t.addCode("while(")
		t.Transpile(node.Statement)
		t.addCode(") {\n")
		t.env = t.env.Enclose(t.env.Fn)
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = t.env.Open()

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
				for _, ty := range v.Type {
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

		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Operator == ".." {
			node.Operator = ":"
		}

		t.addCode(node.Operator)

		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Right != nil {
			t.Transpile(node.Right)

			if node.Operator == "[" {
				t.addCode("]")
			}

			if node.Operator == "[[" {
				t.addCode("]]")
			}
		}

	case *ast.IfExpression:
		t.addCode("if(")
		t.Transpile(node.Condition)
		t.addCode("){\n")
		t.env = t.env.Enclose(t.env.Fn)
		t.Transpile(node.Consequence)
		t.env = t.env.Open()
		t.addCode("}")

		if node.Alternative != nil {
			t.addCode(" else {\n")
			t.env = t.env.Enclose(t.env.Fn)
			t.Transpile(node.Alternative)
			t.env = t.env.Open()
			t.addCode("}")
		}

	case *ast.FunctionLiteral:
		t.env = t.env.Enclose(
			environment.Object{
				Token: node.Token,
				Name:  node.Name.Value,
				Type:  node.Type,
			},
		)

		if node.Name.String() != "" {
			t.addCode(node.Name.String())
		}

		if node.Method != "" {
			t.addCode("." + node.Method)
		}

		if node.Operator != "" {
			t.addCode(" " + node.Operator + " ")
		}

		t.addCode("function")
		t.addCode("(")

		for i, p := range node.Parameters {
			t.env.SetVariable(
				p.TokenLiteral(),
				environment.Object{
					Token: p.Token,
					Name:  p.Name,
					Type:  p.Type,
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
		t.Transpile(node.Body)
		t.env = t.env.Open()
		t.addCode("\n}")

	case *ast.Decorator:
		n := t.Transpile(node.Type)
		t.env.SetClass(n.Item().Value, environment.Object{Class: node.Classes})

	case *ast.CallExpression:
		tt, typeExists := t.env.GetType(node.Name, false)

		if !typeExists {
			tt, typeExists = t.env.GetType(node.Name, true)
		}

		if node.Name != "" && typeExists {
			if typeExists {
				name := tt.Type[0].Name

				if name == "struct" {
					name = "structure"
				}

				if name == "dataframe" {
					name = "data.frame"
				}

				if name == "int" || name == "num" || name == "char" {
					name = "c"
				}

				t.addCode(name)
			}
		} else {
			t.Transpile(node.Function)
		}
		t.addCode("(")

		for i, a := range node.Arguments {
			t.Transpile(a.Value)
			if i < len(node.Arguments)-1 {
				t.addCode(", ")
			}
		}

		class, hasClass := t.env.GetClass(node.Name)

		if hasClass {
			if len(node.Arguments) > 0 {
				t.addCode(", ")
			}
			t.addCode("class = c(\"" + strings.Join(class.Class, "\", \"") + "\")")
		}

		if typeExists && tt.Type[0].Name == "struct" && !hasClass {
			if len(node.Arguments) > 0 {
				t.addCode(", ")
			}
			t.addCode("class = \"" + node.Name + "\"")
			t.outType()
		}

		t.addCode(")")

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
	t.addCode(l.Name.Value + " = ")
}

func (t *Transpiler) transpileConstStatement(c *ast.ConstStatement) {
	t.addCode(c.Name.Value + " = ")
}

func (t *Transpiler) inType(name []string) {
	t.opts.inType = true
	t.opts.typeClass = name
}

func (t *Transpiler) outType() {
	t.opts.inType = false
	t.opts.typeClass = []string{}
}
