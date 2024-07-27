package ast

type Type struct {
	Name string
	List bool
}

func allTypesMatch(types1 []*Type, types2 []*Type) (bool, []*Type, []*Type) {
	var allMatch []bool
	var expects []*Type

	for _, t1 := range types1 {
		matches := false
		for _, t2 := range types2 {
			if t1.Name == t2.Name && t1.List == t2.List {
				matches = true
			}
		}

		if !matches {
			expects = append(expects, t1)
		}

		allMatch = append(allMatch, matches)
	}

	matches := true
	for _, v := range allMatch {
		if !v {
			matches = v
		}
	}

	return matches, types2, expects
}
