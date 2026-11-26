package lsp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vapourlang/vapour/lexer"
)

// findRoot returns the directory to check along with the file
// and whether to descend into its sub-directories:
// in a package (directory with a DESCRIPTION) the whole package,
// otherwise only the file's own directory, so a loose file
// e.g.: in ~/Projects does not pull in every project below it.
func findRoot(file string) (string, bool) {
	dir := filepath.Dir(file)

	for d := dir; ; {
		if info, err := os.Stat(filepath.Join(d, "DESCRIPTION")); err == nil && !info.IsDir() {
			return d, true
		}

		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}

	return dir, false
}

func (l *LSP) readDir(root string, recursive bool) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path == root {
				return nil
			}

			// e.g.: .git
			if !recursive || strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			return nil
		}

		return l.walk(path)
	})
}

func (l *LSP) walk(path string) error {
	var err error

	ext := filepath.Ext(path)

	if ext != ".vp" {
		return nil
	}

	path = filepath.Clean(path)

	// documents open in the editor may have unsaved changes
	var fl []byte
	if content, ok := l.documents[path]; ok {
		fl = []byte(content)
	} else {
		fl, err = os.ReadFile(path)

		if err != nil {
			return err
		}
	}

	rfl := lexer.File{
		Path:    path,
		Content: fl,
	}

	l.files = append(l.files, rfl)

	return nil
}
