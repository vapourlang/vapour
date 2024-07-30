package walker

import (
	"strings"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/token"
)

// Expect a type, where expectation is the left-side node
// node is the right-side node to traverse
func (w *Walker) expectType(node ast.Node, tok token.Item, expectation []*ast.Type) {
	actual, _ := w.Walk(node)
	ok, missing := w.typesIn(expectation, actual)

	if ok {
		return
	}

	w.addErrorf(
		tok,
		diagnostics.Fatal,
		"token `%v` type mismatch, assigning (%v) to (%v), missing (%v)",
		tok.Value,
		typeString(actual),
		typeString(expectation),
		typeString(missing),
	)
}

func (w *Walker) typesIn(expectation []*ast.Type, actual []*ast.Type) (bool, []*ast.Type) {
	var oks []bool
	var missing []*ast.Type

	for _, t := range actual {
		ok := w.typeIn(t, expectation)
		oks = append(oks, ok)

		if ok {
			continue
		}

		missing = append(missing, t)
	}

	return any(oks...), missing
}

// check if any value is false
func any(values ...bool) bool {
	for _, v := range values {
		if !v {
			return false
		}
	}

	return true
}

// Check that the actual type can be found in the list of expected types
func (w *Walker) typeIn(t *ast.Type, compare []*ast.Type) bool {
	for _, c := range compare {
		if c.Name == t.Name && c.List == t.List {
			return true
		}

		// check custom types
		a, exists := w.env.GetType(c.Name)

		// it's not a custom type, can't match
		if !exists {
			return false
		}

		// check whether type matches
		for _, at := range a.Type {
			if at.Name == t.Name && at.List == t.List {
				return true
			}
		}
	}

	return false
}

func typeString(t []*ast.Type) string {
	var types []string

	for _, v := range t {
		lst := ""
		if v.List {
			lst = "[]"
		}

		str := lst + v.Name
		types = append(types, str)
	}

	return strings.Join(types, ", ")
}

func (w *Walker) typeExists(t *ast.Type) bool {
	_, exists := w.env.GetType(t.Name)
	return exists
}

func (w *Walker) typesExists(t []*ast.Type) bool {
	var exist []bool

	for _, v := range t {
		_, exists := w.env.GetType(v.Name)
		exist = append(exist, exists)
	}

	for _, v := range exist {
		if !v {
			return false
		}
	}

	return true
}

func (w *Walker) allSameTypes(t []*ast.Type) bool {
	var previousTypes *ast.Type
	for i, v := range t {
		if i == 0 {
			previousTypes = v
			continue
		}

		if previousTypes != v {
			return false
		}
	}

	return true
}
