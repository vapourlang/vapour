package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func (v *vapour) readDir() error {
	err := filepath.WalkDir(*v.root, v.walk)

	if err != nil {
		return err
	}

	for _, f := range v.files {
		v.combined = append(v.combined, f.content...)
		v.combined = append(v.combined, byte(' '))
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

	rfl := File{
		path:    path,
		content: fl,
	}

	v.files = append(v.files, rfl)

	return nil
}
