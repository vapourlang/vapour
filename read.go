package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func (d *doctor) readRoot() error {
	return filepath.WalkDir(*d.root, d.walk)
}

func (d *doctor) walk(path string, directory fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if directory.IsDir() {
		return nil
	}

	ext := filepath.Ext(path)

	if ext != ".R" || ext != ".r" {
		return nil
	}

	fl, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	rfl := RFile{
		path:    path,
		content: fl,
	}

	d.rfiles = append(d.rfiles, rfl)

	return nil
}
