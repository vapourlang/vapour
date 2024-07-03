package lexer

import (
	"fmt"
	"os"
)

var ItemName = map[ItemType]string{
	ItemError:             "error",
	ItemIdent:             "identifier",
	ItemDoubleQuote:       "double quote",
	ItemSingleQuote:       "single quote",
	ItemAssign:            "assign",
	ItemLeftCurly:         "curly left",
	ItemRightCurly:        "curly right",
	ItemLeftParen:         "paren left",
	ItemRightParen:        "paren right",
	ItemLeftSquare:        "square left",
	ItemRightSquare:       "square right",
	ItemString:            "string",
	ItemInteger:           "integer",
	ItemFloat:             "float",
	ItemNamespace:         "namespace",
	ItemNamespaceInternal: "namespace internal",
	ItemComment:           "comment",
	ItemSpecialComment:    "special comment",
	ItemRoxygenTagAt:      "@roxygen",
	ItemRoxygenTag:        "roxygen tag",
	ItemRoxygenTagContent: "roxygen content",
	ItemTypeDef:           "type",
	ItemTypeVar:           "type variable",
	ItemDoubleEqual:       "double equal",
	ItemLessThan:          "less than",
	ItemGreaterThan:       "greater than",
	ItemNotEqual:          "not equal",
	ItemLessOrEqual:       "less or equal",
	ItemGreaterOrEqual:    "greater or equal",
	ItemBool:              "boolean",
	ItemDollar:            "dollar",
	ItemComma:             "comma",
	ItemColon:             "colon",
	ItemSemiColon:         "semicolon",
	ItemQuestion:          "question mark",
	ItemBacktick:          "backtick",
	ItemInfix:             "infix",
	ItemIf:                "if",
	ItemBreak:             "break",
	ItemElse:              "else",
	ItemAnd:               "ampersand",
	ItemOr:                "vertical bar",
	ItemReturn:            "return",
	ItemC:                 "C call",
	ItemCall:              "C++ call",
	ItemFortran:           "Fortan call",
	ItemNULL:              "null",
	ItemNA:                "NA",
	ItemNan:               "NaN",
	ItemNACharacter:       "NA character",
	ItemNAReal:            "NA real",
	ItemNAComplex:         "NA complex",
	ItemNAInteger:         "NA integer",
	ItemPipe:              "native pipe",
	ItemModulus:           "modulus",
	ItemDoubleLeftSquare:  "double left square",
	ItemDoubleRightSquare: "double right square",
	ItemFor:               "for loop",
	ItemRepeat:            "repeat",
	ItemWhile:             "while loop",
	ItemNext:              "next",
	ItemIn:                "in",
	ItemFunction:          "function",
	ItemPlus:              "plus",
	ItemMinus:             "minus",
	ItemMultiply:          "multiply",
	ItemDivide:            "divide",
	ItemPower:             "power",
	ItemEOL:               "end of line",
	ItemEOF:               "end of file",
}

func (item Item) String() string {
	return ItemName[item.class]
}

func (l *Lexer) Print() {
	fmt.Fprintln(os.Stdout, "\n----")
	for _, v := range l.Items {
		name := v.String()
		val := v.val
		if val == "\n" {
			val = "\\_n"
		}
		fmt.Fprintf(os.Stdout, "%v: %v [%v]\n", v.class, val, name)
	}
	fmt.Fprintf(os.Stdout, "=====> lexed %v tokens\n", len(l.Items))
}
