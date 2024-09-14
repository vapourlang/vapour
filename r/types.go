package r

import (
	"errors"
	"os"
	"path"

	"github.com/devOpifex/vapour/ast"
)

func getPackagesTypes() ([]ast.Types, error) {
	var types []ast.Types
	entries, err := os.ReadDir(library)

	if err != nil {
		return types, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		typeFile := path.Join(e.Name(), "types.vp")

		if _, err := os.Stat(typeFile); errors.Is(err, os.ErrNotExist) {
			continue
		}
	}

	return types, nil
}
