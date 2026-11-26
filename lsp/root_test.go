package lsp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vapourlang/vapour/diagnostics"
	"github.com/vapourlang/vapour/token"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func paths(l *LSP) []string {
	var ps []string
	for _, f := range l.files {
		ps = append(ps, filepath.Base(filepath.Dir(f.Path))+"/"+filepath.Base(f.Path))
	}
	return ps
}

func TestFindRootLooseFile(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "script.vp"), "")
	write(t, filepath.Join(dir, "other", "broken.vp"), "")

	root, recursive := findRoot(filepath.Join(dir, "script.vp"))

	if root != dir || recursive {
		t.Fatalf("expected %v non recursive, got %v %v", dir, root, recursive)
	}

	l := New()
	if err := l.readDir(root, recursive); err != nil {
		t.Fatal(err)
	}

	if len(l.files) != 1 {
		t.Fatalf("expected only the file's directory, got %v", paths(l))
	}
}

func TestFindRootPackage(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "DESCRIPTION"), "Package: x")
	write(t, filepath.Join(dir, "vp", "x.vp"), "")
	write(t, filepath.Join(dir, "inst", "types.vp"), "")
	write(t, filepath.Join(dir, ".hidden", "skip.vp"), "")

	root, recursive := findRoot(filepath.Join(dir, "vp", "x.vp"))

	if root != dir || !recursive {
		t.Fatalf("expected %v recursive, got %v %v", dir, root, recursive)
	}

	l := New()
	if err := l.readDir(root, recursive); err != nil {
		t.Fatal(err)
	}

	if len(l.files) != 2 {
		t.Fatalf("expected the package's files, got %v", paths(l))
	}
}

func TestAddBlocked(t *testing.T) {
	file := filepath.Clean("/project/x.vp")
	other := filepath.Clean("/project/sub/y.vp")

	ds := addBlocked(nil, diagnostics.Diagnostics{
		diagnostics.NewError(token.Item{File: other, Line: 1}, "parse"),
		diagnostics.NewError(token.Item{File: other, Line: 3}, "parse"),
	}, file, "/project")

	if len(ds) != 1 || !strings.Contains(ds[0].Message, filepath.Join("sub", "y.vp")+":2") {
		t.Fatalf("expected a blocked diagnostic, got %+v", ds)
	}

	// the file's own errors explain it
	ds = addBlocked(nil, diagnostics.Diagnostics{
		diagnostics.NewError(token.Item{File: other}, "parse"),
		diagnostics.NewError(token.Item{File: file}, "parse"),
	}, file, "/project")

	if len(ds) != 0 {
		t.Fatalf("expected nothing, got %+v", ds)
	}
}
