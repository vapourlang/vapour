package token

import (
	"fmt"
)

var ItemName = map[ItemType]string{
	ItemError:                "error",
	ItemDefer:                "defer",
	ItemIdent:                "identifier",
	ItemDoubleQuote:          "double quote",
	ItemSingleQuote:          "single quote",
	ItemAssign:               "assign",
	ItemLeftCurly:            "curly left",
	ItemRightCurly:           "curly right",
	ItemLeftParen:            "paren left",
	ItemRightParen:           "paren right",
	ItemLeftSquare:           "square left",
	ItemRightSquare:          "square right",
	ItemString:               "string",
	ItemInteger:              "integer",
	ItemFloat:                "float",
	ItemNamespace:            "namespace",
	ItemNamespaceInternal:    "namespace internal",
	ItemComment:              "comment",
	ItemDoubleEqual:          "double equal",
	ItemLessThan:             "less than",
	ItemGreaterThan:          "greater than",
	ItemNotEqual:             "not equal",
	ItemLessOrEqual:          "less or equal",
	ItemGreaterOrEqual:       "greater or equal",
	ItemBool:                 "boolean",
	ItemDollar:               "dollar sign",
	ItemComma:                "comma",
	ItemColon:                "colon",
	ItemQuestion:             "question mark",
	ItemBacktick:             "backtick",
	ItemInfix:                "infix",
	ItemIf:                   "if",
	ItemBreak:                "break",
	ItemElse:                 "else",
	ItemAnd:                  "ampersand",
	ItemOr:                   "vertical bar",
	ItemReturn:               "return",
	ItemNULL:                 "null",
	ItemNA:                   "NA",
	ItemNan:                  "NaN",
	ItemNAString:             "NA string",
	ItemNAReal:               "NA real",
	ItemNAComplex:            "NA complex",
	ItemNAInteger:            "NA integer",
	ItemPipe:                 "native pipe",
	ItemModulus:              "modulus",
	ItemDoubleLeftSquare:     "double left square",
	ItemDoubleRightSquare:    "double right square",
	ItemFor:                  "for loop",
	ItemRepeat:               "repeat",
	ItemWhile:                "while loop",
	ItemNext:                 "next",
	ItemIn:                   "in",
	ItemFunction:             "function",
	ItemPlus:                 "plus",
	ItemMinus:                "minus",
	ItemMultiply:             "multiply",
	ItemDivide:               "divide",
	ItemPower:                "power",
	ItemNewLine:              "new line",
	ItemEOF:                  "end of file",
	ItemTypes:                "type",
	ItemTypesPkg:             "type package",
	ItemTypesList:            "list type",
	ItemTypesDecl:            "type declaration",
	ItemRange:                "range",
	ItemLet:                  "let",
	ItemConst:                "const",
	ItemBang:                 "bang",
	ItemArrow:                "arrow function",
	ItemThreeDot:             "elipsis",
	ItemDecorator:            "decorator",
	ItemDecoratorClass:       "decorator class",
	ItemDecoratorGeneric:     "decorator generic",
	ItemDecoratorDefault:     "decorator default",
	ItemDecoratorFactor:      "decorator factor",
	ItemDecoratorEnvironment: "decorator env",
	ItemAttribute:            "attribute",
	ItemObjList:              "object list",
	ItemObjFunc:              "object function",
	ItemObjDataframe:         "object dataframe",
	ItemObjStruct:            "object struct",
	ItemObjObject:            "object object",
	ItemObjMatrix:            "object matrix",
	ItemObjFactor:            "object factor",
}

func (t ItemType) String() string {
	k, _ := ItemName[t]
	return k
}

func (item Item) String() string {
	return ItemName[item.Class]
}

func pad(str string, min int) string {
	out := str
	l := len(str)

	var i int
	for l < min {
		pad := "-"

		if i == 0 || i == min {
			pad = " "
		}
		out = out + pad
		l = len(out)
		i++
	}

	return out
}

func (i Item) Print() {
	name := i.String()
	val := i.Value
	if val == "\n" {
		val = "\\_n"
	}

	name = pad(name, 30)
	fmt.Printf(
		"%s `%v` \t [file: %v line: %v, char: %v]\n",
		name, val, i.File, i.Line, i.Char,
	)
}

func (i Items) Print() {
	for _, v := range i {
		v.Print()
	}
}
