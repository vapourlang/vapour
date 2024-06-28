package lexer

import (
	"fmt"
	"os"
	"testing"
)

var itemName = map[itemType]string{
	itemError:         "error",
	itemIdent:         "identifier",
	itemDoubleQuote:   "\"",
	itemSingleQuote:   "'",
	itemAssign:        "assign",
	itemLeftCurly:     "{",
	itemRightCurly:    "}",
	itemLeftParen:     "(",
	itemRightParen:    ")",
	itemString:        "string",
	itemInteger:       "integer",
	itemFloat:         "float",
	itemNamespace:     "::",
	itemMathOperation: "mathOperation",
}

func print(l *lexer) {
	for _, v := range l.items {
		name := itemName[v.class]
		fmt.Fprintf(os.Stdout, "%v: %v - %v\n", v.class, v.val, name)
	}
}

func TestBasic(t *testing.T) {
	code := `x <- 1 + 2 - 1
					y <- 2`

	l := &lexer{
		input: code,
	}

	fmt.Fprintln(os.Stdout, "Running lexer")
	l.run()
	fmt.Fprintf(os.Stdout, "lexed %v tokens\n", len(l.items))

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}
