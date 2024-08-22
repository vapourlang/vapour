package walker

import (
	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/environment"
)

func (w *Walker) allTypesIdentical(types []*ast.Type) bool {
	if len(types) == 0 {
		return true
	}

	types, ok := w.getNativeTypes(types)

	if !ok {
		return false
	}

	var previousType *ast.Type
	for i, t := range types {
		if i == 0 {
			previousType = t
			continue
		}

		if t.Name != previousType.Name || t.List != previousType.List {
			return false
		}
	}

	return true
}

func typeIdentical(t1, t2 *ast.Type) bool {
	return t1.Name == t2.Name && t1.List == t2.List
}

func acceptAny(types ast.Types) bool {
	for _, t := range types {
		if t.Name == "any" {
			return true
		}
	}
	return false
}

func (w *Walker) typesValid(valid, actual ast.Types) bool {
	// we don't have the type
	if len(valid) == 0 {
		return true
	}

	if acceptAny(valid) {
		return true
	}

	for _, l := range actual {
		if w.typeValid(l, valid) {
			continue
		}

		return false
	}

	return true
}

func (w *Walker) typeValid(t *ast.Type, valid ast.Types) bool {
	// we just don't have the type
	// could be base R dataset
	if t.Name == "" {
		return true
	}

	for _, v := range valid {
		if typeIdentical(t, v) {
			return true
		}

		if v.Name == "int" && t.Name == "num" && v.List == t.List {
			return true
		}
	}

	return false
}

func (w *Walker) validMathTypes(types ast.Types) bool {
	types, ok := w.getNativeTypes(types)

	if !ok {
		return false
	}

	for _, t := range types {
		if !contains(t.Name, []string{"int", "num", "na"}) {
			return false
		}
	}
	return true
}

func contains(value string, arr []string) bool {
	for _, a := range arr {
		if value == a {
			return true
		}
	}
	return false
}

func (w *Walker) retrieveNativeTypes(types, nativeTypes ast.Types) (ast.Types, bool) {
	for _, t := range types {
		if environment.IsNativeType(t.Name) {
			nativeTypes = append(nativeTypes, t)
			continue
		}

		customType, exists := w.env.GetType(t.Name)

		if exists && customType.Object == "vector" {
			return w.retrieveNativeTypes(customType.Type, nativeTypes)
		}

		return nativeTypes, false
	}

	return nativeTypes, true
}

func (w *Walker) getNativeTypes(types ast.Types) (ast.Types, bool) {
	return w.retrieveNativeTypes(types, ast.Types{})
}

func (w *Walker) validIteratorTypes(types ast.Types) bool {
	var valid []bool
	for _, t := range types {
		if contains(t.Name, []string{"int", "num", "char"}) {
			valid = append(valid, true)
			continue
		}

		custom, exists := w.env.GetType(t.Name)
		if !exists {
			valid = append(valid, false)
			continue
		}

		if len(custom.Type) > 0 && allLists(custom.Type) {
			valid = append(valid, true)
			continue
		}

		if custom.Object == "list" {
			valid = append(valid, true)
			continue
		}

		valid = append(valid, false)
	}
	return allTrue(valid)
}

func allLists(types ast.Types) bool {
	for _, t := range types {
		if !t.List {
			return false
		}
	}

	return true
}

func allTrue(values []bool) bool {
	for _, b := range values {
		if !b {
			return false
		}
	}
	return true
}

func (w *Walker) checkIdentifier(node *ast.Identifier) {
	v, exists := w.env.GetVariable(node.Value, true)

	if exists {
		if v.CanMiss {
			w.addWarnf(
				node.Token,
				"`%v` might be missing",
				node.Token.Value,
			)
		}

		if v.IsConst {
			w.addFatalf(
				node.Token,
				"`%v` is a constant",
				node.Value,
			)
		}

		w.env.SetVariableUsed(node.Value)
		return
	}

	_, exists = w.env.GetType(node.Value)

	if exists {
		w.env.SetTypeUsed(node.Value)
		return
	}

	// we are actually declaring variable in a call
	if w.state != "call" {
		w.addWarnf(
			node.Token,
			"`%v` not found",
			node.Value,
		)
	}
}

func (w *Walker) checkIfIdentifier(node ast.Node) {
	switch n := node.(type) {
	case *ast.Identifier:
		w.checkIdentifier(n)
	}
}
