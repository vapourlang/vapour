package lexer

import (
	"fmt"
	"os"
	"testing"
)

func print(l *lexer) {
	for _, v := range l.items {
		fmt.Fprintf(os.Stdout, "[%v - %v]\n", v.val, v.class)
	}
}

func TestBasic(t *testing.T) {
	code := "x <- 1 + 1"

	l := &lexer{
		input: code,
	}

	fmt.Fprintln(os.Stdout, "Running lexer")
	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}
}
