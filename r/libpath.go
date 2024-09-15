package r

import (
	"encoding/json"
)

func LibPath() []string {
	var paths []string
	output, err := Callr(`paths <- paste0(.libPaths(), collapse = "\",\"")
	paths <- paste0("[\"", paths, "\"]")
	cat(paths)`)

	if err != nil {
		return paths
	}

	err = json.Unmarshal(output, &paths)

	if err != nil {
		return paths
	}

	return paths
}
