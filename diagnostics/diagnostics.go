package diagnostics

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/devOpifex/vapour/token"
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

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
var Bold = "\033[1m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

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
		v.Print()
	}

	return out.String()
}

func (v Diagnostic) Print() {
	var out bytes.Buffer
	out.WriteString("[" + prefix(v.Severity) + "]\t")
	out.WriteString(v.Token.File)
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%v", v.Token.Line))
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%v", v.Token.Char))
	out.WriteString(" " + v.Message)
	fmt.Println(out.String())
}

func (d Diagnostics) Print() {
	fmt.Println(d.String())
}

func prefix(s Severity) string {
	if s == Fatal {
		return Red + "ERROR" + Reset
	}

	if s == Warn {
		return Yellow + "WARN" + Reset
	}

	if s == Info {
		return Blue + "INFO" + Reset
	}

	return Green + "HINT" + Reset
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
