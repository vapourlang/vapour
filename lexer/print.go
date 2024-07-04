package lexer

import "fmt"

func (l *Lexer) Print() {
	fmt.Printf("Lexer with %v tokens\n", len(l.Items))
	l.Items.Print()
}
