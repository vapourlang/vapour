package walker

import (
	"fmt"

	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/token"
)

func (w *Walker) addError(tok token.Item, s diagnostics.Severity, m string) {
	w.errors = append(w.errors, diagnostics.New(tok, m, s))
}

func (w *Walker) addErrorf(tok token.Item, s diagnostics.Severity, fm string, a ...interface{}) {
	str := fmt.Sprintf(fm, a...)
	w.errors = append(w.errors, diagnostics.New(tok, str, s))
}

func (w *Walker) addFatalf(tok token.Item, fm string, a ...interface{}) {
	str := fmt.Sprintf(fm, a...)
	w.errors = append(w.errors, diagnostics.New(tok, str, diagnostics.Fatal))
}

func (w *Walker) addWarnf(tok token.Item, fm string, a ...interface{}) {
	str := fmt.Sprintf(fm, a...)
	w.errors = append(w.errors, diagnostics.New(tok, str, diagnostics.Warn))
}

func (w *Walker) addInfof(tok token.Item, fm string, a ...interface{}) {
	str := fmt.Sprintf(fm, a...)
	w.errors = append(w.errors, diagnostics.New(tok, str, diagnostics.Info))
}

func (w *Walker) addHintf(tok token.Item, fm string, a ...interface{}) {
	str := fmt.Sprintf(fm, a...)
	w.errors = append(w.errors, diagnostics.New(tok, str, diagnostics.Hint))
}