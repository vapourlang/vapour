package environment

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

type Code struct {
	lines []string
}

func (c *Code) add(line string) {
	c.lines = append(c.lines, line)
}

func (c *Code) String() string {
	return strings.Join(c.lines, "\n")
}

func (e *Environment) GenerateTypes() *Code {
	code := &Code{}

	for typeName, typeObject := range e.types {
		if IsNativeType(typeName) || IsNativeObject(typeName) {
			continue
		}

		class, classExists := e.GetClass(typeName)

		if classExists {
			code.add("@class(" + strings.Join(class.Value.Classes, ", ") + ")")
		}

		curlyLeft := "{"
		if len(typeObject.Attributes) == 0 {
			curlyLeft = ""
		}

		if typeObject.Object != "impliedList" {
			code.add("type " + typeName + ": " + typeObject.Object + " " + curlyLeft)
		}

		if typeObject.Object == "impliedList" {
			code.add("type " + typeName + ": " + collaseTypes(typeObject.Type) + " " + curlyLeft)
		}

		if len(typeObject.Attributes) == 0 {
			continue
		}

		for i, a := range typeObject.Attributes {
			sep := ""
			if i < len(typeObject.Attributes)-1 {
				sep = ","
			}
			code.add("\t" + a.Name + ": " + collaseTypes(a.Type) + sep)
		}
		code.add("}")
	}

	return code
}

func collaseTypes(types []*ast.Type) string {
	var str []string

	for _, t := range types {
		typeString := ""
		if t.List {
			typeString += "[]"
		}

		typeString += t.Name

		str = append(str, typeString)
	}

	return strings.Join(str, " | ")
}

func IsNativeType(name string) bool {
	for _, t := range baseTypes {
		if name == t {
			return true
		}
	}
	return false
}

func IsNativeObject(name string) bool {
	for _, t := range baseObjects {
		if name == t {
			return true
		}
	}
	return false
}

var packagesLoaded []string

func isLoaded(library string) bool {
	for _, p := range packagesLoaded {
		if p == library {
			return true
		}
	}

	return false
}

func (env *Environment) LoadPackageTypes(pkg string) error {
	if library == "" {
		return nil
	}

	if isLoaded(pkg) {
		return nil
	}

	packagesLoaded = append(packagesLoaded, library)

	typeFile := path.Join(library, pkg, "types.vp")

	if _, err := os.Stat(typeFile); errors.Is(err, os.ErrNotExist) {
		return err
	}

	content, err := os.ReadFile(typeFile)

	if err != nil {
		return err
	}

	// lex
	l := lexer.NewCode(typeFile, string(content))
	l.Run()

	if l.HasError() {
		return errors.New("failed to lex types file")
	}

	// parse
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		return errors.New("failed to lex types file")
	}

	// range over the Statements
	// these should all be type declarations
	for _, p := range prog.Statements {
		switch node := p.(type) {
		case *ast.TypeStatement:
			env.SetType(
				Type{
					Token:      node.Token,
					Type:       node.Type,
					Attributes: node.Attributes,
					Object:     node.Object,
					Name:       node.Name,
					Package:    pkg,
					Used:       true,
				},
			)
		}
	}

	return nil
}
