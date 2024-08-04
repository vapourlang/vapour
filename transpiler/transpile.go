package transpiler

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/environment"
)

type Transpiler struct {
	code string
	env  *environment.Environment
	opts options
}

type options struct {
	inType    bool
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
		t.transpileLetStatement(node)
		t.env.SetVariable(
			node.Name.Value,
			environment.Object{Token: node.Token, Type: node.Name.Type},
		)
		t.Transpile(node.Value)
		t.addCode("\n")

	case *ast.NewLine:
		t.addCode("\n")

	case *ast.ConstStatement:
		t.env.SetVariable(
			node.Name.Value,
			environment.Object{Token: node.Token, Type: node.Name.Type},
		)
		t.transpileConstStatement(node)
		t.Transpile(node.Value)
		t.addCode("\n")

	case *ast.ReturnStatement:
		t.addCode("\nreturn(")
		t.Transpile(node.ReturnValue)
		t.addCode(")")

	case *ast.TypeStatement:
		_, exists := t.env.GetType(node.Name.Value)

		if !exists {
			t.env.SetType(
				node.Name.Value,
				environment.Object{
					Token:  node.Token,
					Type:   node.Type,
					Object: node.Object,
					Name:   node.Name.Value,
					List:   node.List,
				},
			)
		}

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
		// check types
		_, varExists := t.env.GetVariable(node.Value, true)
		tt, typeExists := t.env.GetType(node.Value)

		if !varExists && typeExists {
			name := tt.Type[0].Name

			if name == "struct" {
				name = "structure"
				t.inType([]string{node.Value, tt.Object[0].Name})
			}

			if name == "dataframe" {
				name = "data.frame"
			}

			if name == "int" || name == "num" || name == "char" {
				name = "c"
			}

			t.addCode(name)
		} else {
			t.addCode(node.Value)
		}

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
		t.addCode("\n}")
		t.env = t.env.Open()

	case *ast.While:
		t.addCode("while(")
		t.Transpile(node.Statement)
		t.addCode(") {")
		t.env = t.env.Enclose(t.env.Fn)
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = t.env.Open()

	case *ast.InfixExpression:
		t.Transpile(node.Left)

		if node.Operator == "in" {
			t.addCode(" ")
		}

		t.addCode(node.Operator)

		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Right != nil {
			t.Transpile(node.Right)
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
					Token: node.Token,
					Name:  node.Name.Value,
				},
			)

			t.addCode(p.Name)

			if p.Operator == "=" {
				t.addCode(" = ")
				t.Transpile(p.Default)
			}

			if i < len(node.Parameters)-1 {
				t.addCode(",\n")
			}
		}
		t.addCode(") {")
		t.Transpile(node.Body)
		t.env = t.env.Open()
		t.addCode("\n}")

	case *ast.CallExpression:
		t.Transpile(node.Function)
		t.addCode("(")

		for i, a := range node.Arguments {
			t.Transpile(a)
			if i < len(node.Arguments)-1 {
				t.addCode(", ")
			}
		}

		t.addCode(")\n")
		if t.opts.inType {
			var classes string
			for i, v := range t.opts.typeClass {
				classes += "\"" + v + "\""
				if i < len(t.opts.typeClass)-1 {
					classes += ", "
				}
			}
			t.addCode(", class = c(" + classes + ")")
			t.outType()
		}
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
	return t.code
}

func (t *Transpiler) addCode(code string) {
	t.code += code
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
