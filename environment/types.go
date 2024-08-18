package environment

import (
	"strings"

	"github.com/devOpifex/vapour/ast"
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
		if IsBaseType(typeName) {
			continue
		}

		class, classExists := e.GetClass(typeName)

		if classExists {
			code.add("@class(" + strings.Join(class.Class, ", ") + ")")
		}

		curlyLeft := "{"
		if len(typeObject.Attributes) == 0 {
			curlyLeft = ""
		}

		code.add("type " + typeName + ": " + collaseTypes(typeObject.Type) + " " + curlyLeft)

		if len(typeObject.Attributes) == 0 {
			continue
		}

		for i, a := range typeObject.Attributes {
			sep := ""
			if i > len(typeObject.Attributes)-1 {
				sep = ","
			}
			code.add("\t" + a.Name.Value + ": " + collaseTypes(a.Type) + sep)
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

func IsBaseType(name string) bool {
	for _, t := range baseTypes {
		if name == t {
			return true
		}
	}
	return false
}
