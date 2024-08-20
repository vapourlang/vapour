package walker

import "github.com/devOpifex/vapour/ast"

func allTypesIdentical(types []*ast.Type) bool {
	if len(types) == 0 {
		return true
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
