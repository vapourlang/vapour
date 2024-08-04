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

func NewError(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Fatal,
	}
}

func NewWarning(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Warn,
	}
}

func NewInfo(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Info,
	}
}

func NewHint(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Hint,
	}
}

func (d Diagnostics) String() string {
	var out bytes.Buffer

	for i, v := range d {
		out.WriteString("[" + prefix(v.Severity) + "]\t")
		out.WriteString("file ")
		out.WriteString(v.Token.File)
		out.WriteString(", line ")
		out.WriteString(fmt.Sprintf("%v", v.Token.Line+1))
		out.WriteString(fmt.Sprintf(", character %v", v.Token.Char+1))
		out.WriteString(": " + v.Message)
		if i < len(d)-1 {
			out.WriteString("\n")
		}
	}

	return out.String()
}

func (d Diagnostics) Print() {
	fmt.Println(d.String())
}

func prefix(s Severity) string {
	if s == Fatal {
		return "ERROR"
	}

	if s == Warn {
		return "WARN"
	}

	if s == Info {
		return "INFO"
	}

	return "HINT"
}
