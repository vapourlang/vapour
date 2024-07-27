package ast

type Type struct {
	Name string
	List bool
}

func AllTypesMatch(actual []*Type, expected []*Type) (bool, []*Type, []*Type) {
	var allMatch []bool
	var missing []*Type

	for _, t1 := range actual {
		matches := false
		for _, t2 := range expected {
			if t1.Name == t2.Name && t1.List == t2.List {
				matches = true
			}
		}

		if !matches {
			missing = append(missing, t1)
		}

		allMatch = append(allMatch, matches)
	}

	matches := true
	for _, v := range allMatch {
		if !v {
			matches = v
		}
	}

	return matches, expected, missing
}
