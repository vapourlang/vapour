package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vapourlang/vapour/lexer"
)

func (v *vapour) readDir() error {
	err := filepath.WalkDir(*v.root, v.walk)

	if err != nil {
		return err
	}

	return nil
}

func (v *vapour) walk(path string, directory fs.DirEntry, err error) error {
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

	v.files = append(v.files, rfl)

	return nil
}
