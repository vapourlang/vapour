package diagnostics

import (
	"bytes"
	"fmt"

	"github.com/vapourlang/vapour/cli"
	"github.com/vapourlang/vapour/token"
)

// To match LSP specs
type Severity int

const (
	Fatal Severity = iota
	Warn
	Hint
	Info
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

	for _, v := range d {
		out.WriteString(v.String())
	}

	return out.String()
}

func (v Diagnostic) String() string {
	var out bytes.Buffer
	out.WriteString("[" + prefix(v.Severity) + "]\t")
	out.WriteString(v.Token.File)
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%v", v.Token.Line))
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%v", v.Token.Char))
	out.WriteString(" " + v.Message + "\n")
	return out.String()
}

func (d Diagnostics) Print() {
	fmt.Printf("%v", d.String())
}

func prefix(s Severity) string {
	if s == Fatal {
		return cli.Red + "ERROR" + cli.Reset
	}

	if s == Warn {
		return cli.Yellow + "WARN" + cli.Reset
	}

	if s == Info {
		return cli.Blue + "INFO" + cli.Reset
	}

	return cli.Green + "HINT" + cli.Reset
}

func (ds Diagnostics) UniqueLine() Diagnostics {
	uniques := Diagnostics{}

	set := make(map[int]bool)
	for _, d := range ds {
		_, ok := set[d.Token.Line]

		if ok {
			continue
		}

		set[d.Token.Line] = true

		uniques = append(uniques, d)
	}

	return uniques
}

func (ds Diagnostics) Unique() Diagnostics {
	uniques := Diagnostics{}

	set := make(map[string]bool)
	for _, d := range ds {
		key := fmt.Sprintf("%d%d", d.Token.Line, d.Token.Char)
		_, ok := set[key]

		if ok {
			continue
		}

		set[key] = true

		uniques = append(uniques, d)
	}

	return uniques
}
