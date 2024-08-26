package walker

import (
	"fmt"

	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/token"
)

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

func (w *Walker) HasDiagnostic() bool {
	return len(w.errors) > 0
}

func (w *Walker) HasError() bool {
	for _, v := range w.errors {
		if v.Severity != diagnostics.Info && v.Severity != diagnostics.Hint {
			return true
		}
	}
	return false
}

func (w *Walker) Errors() diagnostics.Diagnostics {
	return w.errors
}
