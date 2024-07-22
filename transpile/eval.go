package transpile

import (
	"github.com/devOpifex/vapour/ast"
)

func Transpile(node ast.Node, env *ast.Environment, prog string) ast.Node {
	switch node := node.(type) {

	case *ast.Program:
		return transpileProgram(node, env)

	case *ast.LetStatement:
		val := Transpile(node.Value, env, prog)
		env.Set(node.Name.Value, node)
		return val

	case *ast.ConstStatement:
		val := Transpile(node.Value, env, prog)
		env.Set(node.Name.Value, node)
		return val
	}

	return nil
}

func transpileProgram(program *ast.Program, env *ast.Environment) ast.Node {
	var prog string
	for _, statement := range program.Statements {
		Transpile(statement, env, prog)
	}

	return nil
}

func transpileBlockStatement(
	block *ast.BlockStatement,
	env *ast.Environment,
	prog string,
) ast.Node {
	var result ast.Node
	for _, statement := range block.Statements {
		result = Transpile(statement, env, prog)
	}

	return result
}

func transpileExpressions(
	exps []ast.Expression,
	env *ast.Environment,
	prog string,
) []ast.Node {
	var result []ast.Node

	for _, e := range exps {
		evaluated := Transpile(e, env, prog)
		result = append(result, evaluated)
	}

	return result
}
