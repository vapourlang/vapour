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
	typeClass []string
}

func (t *Transpiler) Env() *environment.Environment {
	return t.env
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
			t.addCode("\n")
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
			t.addCode("\n")
		}

	case *ast.ReturnStatement:
		t.addCode("\nreturn(")
		t.Transpile(node.ReturnValue)
		t.addCode(")")

	case *ast.DeferStatement:
		t.addCode("\non.exit((")
		t.Transpile(node.Func)
		t.addCode(")())")

	case *ast.TypeStatement:
		_, exists := t.env.GetType(node.Name.Value)

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

	case *ast.Attrbute:
		t.addCode(node.Value)
		return node

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
		t.addCode(")\n")

	case *ast.StringLiteral:
		t.addCode(node.Token.Value + node.Str + node.Token.Value)

	case *ast.PrefixExpression:
		t.addCode("(")
		t.addCode(node.Operator)
		t.Transpile(node.Right)
		t.addCode(")")

	case *ast.For:
		t.addCode("for(")
		t.addCode(node.Name.Name.Value)
		t.addCode(" in ")
		t.Transpile(node.Vector)
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

		if node.Operator == "<-" {
			node.Operator = "<<-"
		}

		t.addCode(node.Operator)

		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Right != nil {
			t.Transpile(node.Right)
		}

	case *ast.Square:
		t.popCode()
		t.addCode(node.Token.Value)
		for i, s := range node.Statements {
			t.Transpile(s)
			if i < len(node.Statements)-1 {
				t.addCode(", ")
			}
		}
		if node.Token.Value == "[" {
			t.addCode("]")
		}

		if node.Token.Value == "[[" {
			t.addCode("]]")
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

		if node.Method != "" && node.Method != "any" {
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
		if node.Body != nil {
			t.Transpile(node.Body)
		}

		if node.Body == nil {
			t.addCode("\nUseMethod(\"" + node.Name.Value + "\")")
		}

		t.env = t.env.Open()
		t.addCode("\n}")

	case *ast.DecoratorClass:
		n := t.Transpile(node.Type)
		t.env.SetClass(n.Item().Value, environment.Object{Class: node.Classes})

	case *ast.DecoratorGeneric:
		return t.Transpile(node.Func)

	case *ast.CallExpression:
		tt, typeExists := t.env.GetType(node.Name)

		// it has a name = it's an object
		if typeExists && tt.Name != "" {
			name := tt.Type[0].Name

			t.addCode("structure(")

			if tt.Type[0].List {
				name = "list"
			}

			if name == "struct" {
				name = ""
			}

			if name == "object" {
				name = "list"
			}

			if name == "dataframe" {
				name = "data.frame"
			}

			if name == "int" || name == "num" || name == "char" {
				name = "c"
			}

			t.addCode(name)

			if name != "" {
				t.addCode("(")
			}

			for i, a := range node.Arguments {
				t.Transpile(a.Value)
				if i < len(node.Arguments)-1 {
					t.addCode(", ")
				}
			}

			if name != "" {
				t.addCode(")")
			}

			class, hasClass := t.env.GetClass(node.Name)
			if hasClass {
				if len(node.Arguments) > 0 {
					t.addCode(", ")
				}
				t.addCode("class = c(\"" + strings.Join(class.Class, "\", \"") + "\")")
			}

			if typeExists && !hasClass {
				if len(node.Arguments) > 0 {
					t.addCode(", ")
				}
				t.addCode("class = c(\"" + node.Name + "\"")

				if name != "" {
					t.addCode(",\"" + name + "\"")
				}

				t.addCode(")")
				t.outType()
			}

			if name == "data.frame" {
				t.addCode(", names = c(")
				for i, v := range tt.Attributes {
					t.addCode("\"" + v.Name.Value + "\"")
					if i < len(tt.Attributes)-1 {
						t.addCode(", ")
					}
				}
				t.addCode(")")
			}

			t.addCode(")")

			return node
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
