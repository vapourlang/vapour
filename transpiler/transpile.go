package transpiler

import (
	"strings"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/environment"
	"github.com/devOpifex/vapour/err"
	"github.com/devOpifex/vapour/token"
)

type Transpiler struct {
	code   string
	errors err.Errors
	env    *environment.Environment
	opts   options
}

type options struct {
	inType    bool
	typeClass []string
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
	case *ast.ExpressionStatement:
		if node.Expression != nil {
			return t.Transpile(node.Expression)
		}

	case *ast.Program:
		return t.transpileProgram(node)

	case *ast.LetStatement:
		t.env.SetVariable(node.Name.Value, environment.Object{Token: node.Token})
		t.transpileLetStatement(node)
		t.Transpile(node.Value)
		t.addCode("\n")

	case *ast.ConstStatement:
		t.env.SetVariable(node.Name.Value, environment.Object{Token: node.Token})
		t.transpileConstStatement(node)
		t.Transpile(node.Value)
		t.addCode("\n")

	case *ast.ReturnStatement:
		t.addCode("\nreturn(")
		t.Transpile(node.ReturnValue)
		t.addCode(")\n")

	case *ast.TypeStatement:
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

	case *ast.Keyword:
		t.addCode(node.Value)

	case *ast.CommentStatement:
		t.addCode("#")
		t.addCode(node.TokenLiteral())
		t.addCode("\n")

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			t.Transpile(s)
		}

	case *ast.Identifier:
		// check types
		_, varExists := t.env.GetVariable(node.Value, true)
		tt, typeExists := t.env.GetType(node.Value)

		if !varExists {

			if typeExists {
				fn := tt.Type[0].Name

				if fn == "struct" {
					fn = "structure"
					t.inType([]string{node.Value, tt.Object[0].Name})
				}

				if fn == "dataframe" {
					fn = "data.frame"
				}

				if fn == "int" || fn == "num" || fn == "char" {
					fn = "c"
				}

				t.addCode(fn)
			} else {
				t.addError(node.Token, node.Value+" does not exist")
			}

		}

		if varExists {
			t.addCode(node.Value)
		}

	case *ast.Boolean:
		t.addCode(node.String())

	case *ast.IntegerLiteral:
		t.addCode(node.Value)

	case *ast.VectorLiteral:
		t.addCode("c(")
		for _, s := range node.Value {
			t.Transpile(s)
		}
		t.addCode(")\n")

	case *ast.SquareRightLiteral:
		t.addCode("]")

	case *ast.StringLiteral:
		t.addCode(node.Token.Value + node.Str + node.Token.Value)

	case *ast.PrefixExpression:
		t.addCode("(")
		t.addCode(node.Operator)
		t.Transpile(node.Right)
		t.addCode(")")

	case *ast.InfixExpression:
		t.Transpile(node.Left)
		t.addCode(" " + node.Operator + " ")
		if node.Right != nil {
			return t.Transpile(node.Right)
		}

	case *ast.IfExpression:
		t.addCode("if(")
		t.Transpile(node.Condition)
		t.addCode("){\n")
		t.Transpile(node.Consequence)
		t.addCode("}")

		if node.Alternative != nil {
			t.addCode(" else {\n")
			t.Transpile(node.Alternative)
			t.addCode("\n}\n")
		}

	case *ast.FunctionLiteral:
		t.env = environment.NewEnclosed(t.env)

		params := []string{}
		for _, p := range node.Parameters {
			t.env.SetVariable(
				p.TokenLiteral(),
				environment.Object{
					Token: node.Token,
					Name:  node.Name.Value,
				},
			)
			params = append(params, p.String())
		}

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
		t.addCode(strings.Join(params, ", "))
		t.addCode(") {\n")
		t.Transpile(node.Body)
		t.addCode("}\n")

	case *ast.CallExpression:
		t.addCode("\n")
		args := []string{}
		for _, a := range node.Arguments {
			args = append(args, a.String())
		}

		t.Transpile(node.Function)
		t.addCode("(")
		t.addCode(strings.Join(args, ", "))

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

		t.addCode(")\n")
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
			t.addCode(")\n")
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

func (t *Transpiler) addError(tok token.Item, m string) {
	t.errors = append(t.errors, err.New(tok, m))
}

func (t *Transpiler) inType(name []string) {
	t.opts.inType = true
	t.opts.typeClass = name
}

func (t *Transpiler) outType() {
	t.opts.inType = false
	t.opts.typeClass = []string{}
}
