package walker

import (
	"strings"

	"github.com/devOpifex/vapour/ast"
)

func (w *Walker) allTypesMatch(actual []*ast.Type, expected []*ast.Type) (bool, []*ast.Type, []*ast.Type) {
	var allMatch []bool
	var missing []*ast.Type

	// expected types that are not found in the actual type
	for _, t1 := range expected {
		matches := w.typeMatch(t1, expected)
		allMatch = append(allMatch, matches)
	}

	// all actual types can be expected
	matches := true
	for _, v := range allMatch {
		if !v {
			matches = v
		}
	}

	return matches, expected, missing
}

// Check that the actual type can be found in the list of expected types
func (w *Walker) typeMatch(actual *ast.Type, expected []*ast.Type) bool {
	for _, t1 := range expected {
		if t1.Name == actual.Name && t1.List == actual.List {
			return true
		}

		// check custom types
		a, exists := w.env.GetType(t1.Name)

		// it's not a custom type, can't match
		if !exists {
			return false
		}

		// check whether type matches
		for _, at := range a.Type {
			if at.Name == actual.Name && at.List == actual.List {
				return true
			}
		}

		e, exists := w.env.GetType(actual.Name)

		if !exists {
			return false
		}

		if a.Name == e.Name {
			return true
		}

		for _, at := range a.Type {
			for _, et := range e.Type {
				if at.Name == et.Name && at.List == et.List {
					return true
				}
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
