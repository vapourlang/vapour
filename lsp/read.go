package lsp

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/devOpifex/vapour/lexer"
)

func (l *LSP) readDir(root string) error {
	err := filepath.WalkDir(root, l.walk)

	if err != nil {
		return err
	}

	return nil
}

func (l *LSP) walk(path string, directory fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if directory.IsDir() {
		return nil
	}

	ext := filepath.Ext(path)

	if ext != ".vp" {
		return nil
	}

	fl, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	rfl := lexer.File{
		Path:    path,
		Content: fl,
	}

	l.files = append(l.files, rfl)

	return nil
}
