package err

import (
	"bytes"
	"fmt"

	"github.com/devOpifex/vapour/token"
)

type Error struct {
	Token   token.Item
	Message string
}

type Errors []Error

func newError(token token.Item, message string) Error {
	return Error{
		Token:   token,
		Message: message,
	}
}

func hasErrors(errors []Error) bool {
	return len(errors) > 0
}

func (errs Errors) String() string {
	var out bytes.Buffer

	for _, v := range errs {
		out.WriteString("[ERROR]\t")
		out.WriteString("line ")
		out.WriteString(fmt.Sprintf("%v", v.Token.Line+1))
		out.WriteString(fmt.Sprintf(", character %v", v.Token.Pos+1))
		out.WriteString(": " + v.Message)
		out.WriteString("\n")
	}

	return out.String()
}

func (errs Errors) Print() {
	fmt.Println(errs.String())
}
