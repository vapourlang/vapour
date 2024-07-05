package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func (v *vapour) readRoot() error {
	return filepath.WalkDir(*v.root, v.walk)
}

func (v *vapour) walk(path string, directory fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if directory.IsDir() {
		return nil
	}

	ext := filepath.Ext(path)

	if ext != ".cp" {
		return nil
	}

	fl, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	rfl := File{
		path:    path,
		content: fl,
	}

	v.files = append(v.files, rfl)

	return nil
}
