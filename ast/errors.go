package ast

import (
	"bytes"
	"fmt"

	"github.com/devOpifex/vapour/token"
)

type astError struct {
	Token   token.Item
	Message string
}

type astErrors []astError

func newError(token token.Item, message string) astError {
	return astError{
		Token:   token,
		Message: message,
	}
}

func hasErrors(errors []astError) bool {
	return len(errors) > 0
}

func (errs astErrors) String() string {
	var out bytes.Buffer

	for _, v := range errs {
		out.WriteString("[ERROR]\t")
		out.WriteString(v.Message)
		out.WriteString(" at line ")
		out.WriteString(fmt.Sprintf("%v", v.Token.Line))
		out.WriteString(fmt.Sprintf(", character %v", v.Token.Pos))
		out.WriteString("\n")
	}

	return out.String()
}

func (errs astErrors) Print() {
	fmt.Printf("%v errors found\n", len(errs))
	fmt.Println(errs.String())
}
