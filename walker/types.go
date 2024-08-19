package walker

import (
	"strings"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/environment"
	"github.com/devOpifex/vapour/token"
)

// Expect a type, where expectation is the left-side node
// node is the right-side node to traverse
func (w *Walker) expectType(node ast.Node, tok token.Item, expectation []*ast.Type) {
	actual, _ := w.Walk(node)
	ok, missing := w.validTypes(expectation, actual)

	if ok {
		return
	}

	w.addErrorf(
		tok,
		diagnostics.Fatal,
		"token `%v` type mismatch, left expects (%v) right returns (%v), missing (%v)",
		tok.Value,
		typeString(expectation),
		typeString(actual),
		typeString(missing),
	)
}

func (w *Walker) validTypes(expectation []*ast.Type, actual []*ast.Type) (bool, []*ast.Type) {
	var oks []bool
	var missing []*ast.Type

	expectation = w.replaceWithNativeTypes(expectation)
	actual = w.replaceWithNativeTypes(actual)

	for _, e := range expectation {
		ok := w.typeValid(e, actual)
		oks = append(oks, ok)

		if ok {
			continue
		}

		missing = append(missing, e)
	}

	return any(oks...), missing
}

func (w *Walker) validReturnTypes(expectation environment.Object, actual []*ast.Type) (bool, []*ast.Type) {
	var oks []bool
	var missing []*ast.Type

	for _, e := range expectation.Type {
		for _, a := range actual {
			types, exists := w.env.GetType(a.Name)

			isFn := false
			if exists && types.Object != nil {
				isFn = types.Object[0].Name == "func"
			}

			if isFn && expectation.Type[0].Name == e.Name {
				return true, missing
			}
		}
		ok := w.typeValid(e, actual)
		oks = append(oks, ok)

		if ok {
			continue
		}

		missing = append(missing, e)
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
func (w *Walker) typeValid(expecting *ast.Type, incoming []*ast.Type) bool {
	// we just don't have the type for this, we skip
	if len(incoming) == 0 {
		return true
	}

	// expects any(thing)
	if expecting.Name == "any" || expecting.Name == "default" {
		return true
	}

	for _, inc := range incoming {
		// we just don't have the type, we skip
		if expecting.Name == "" || inc.Name == "" {
			return true
		}

		// int can go into num
		if inc.Name == "int" && expecting.Name == "num" && inc.List == expecting.List {
			return true
		}

		if inc.Name != expecting.Name {
			return false
		}

		if expecting.Name == inc.Name && expecting.List == inc.List {
			return true
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

		if previousTypes.Name != v.Name && previousTypes.List == v.List {
			return false
		}
	}

	return true
}

func (w *Walker) replaceWithNativeTypes(types []*ast.Type) []*ast.Type {
	var nativeTypes []*ast.Type
	for _, t := range types {
		ts := w.GetNativeType(t)
		if ts != nil {
			nativeTypes = append(nativeTypes, ts...)
		} else {
			nativeTypes = append(nativeTypes, t)
		}
	}
	return nativeTypes
}

func (w *Walker) GetNativeType(t *ast.Type) []*ast.Type {
	if environment.IsBaseType(t.Name) {
		return []*ast.Type{t}
	}

	inherits, ok := w.env.GetType(t.Name)

	if !ok {
		return nil
	}

	if inherits.Object != nil {
		return nil
	}

	if len(inherits.Attributes) > 0 {
		return nil
	}

	var types []*ast.Type
	for _, t := range inherits.Type {
		t := w.GetNativeType(t)
		types = append(types, t...)
	}

	return types
}
