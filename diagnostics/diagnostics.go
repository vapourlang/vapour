package diagnostics

import (
	"bytes"
	"fmt"

	"github.com/devOpifex/vapour/token"
)

// To match LSP specs
type Severity int

const (
	Fatal Severity = iota
	Warn
	Info
	Hint
)

type Diagnostic struct {
	Token    token.Item
	Message  string
	Severity Severity
}

type Diagnostics []Diagnostic

func New(token token.Item, message string, severity Severity) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: severity,
	}
}

func (d Diagnostics) String() string {
	var out bytes.Buffer

	for _, v := range d {
		out.WriteString("[ERROR]\t")
		out.WriteString("line ")
		out.WriteString(fmt.Sprintf("%v", v.Token.Line+1))
		out.WriteString(fmt.Sprintf(", character %v", v.Token.Pos+1))
		out.WriteString(": " + v.Message)
		out.WriteString("\n")
	}

	return out.String()
}

func (d Diagnostics) Print() {
	fmt.Println(d.String())
}
