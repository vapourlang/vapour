package lsp

import (
	"path/filepath"
	"runtime"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vapourlang/vapour/diagnostics"
	"github.com/vapourlang/vapour/token"
)

func TestToLSPSeverity(t *testing.T) {
	cases := map[diagnostics.Severity]protocol.DiagnosticSeverity{
		diagnostics.Fatal: protocol.DiagnosticSeverityError,
		diagnostics.Warn:  protocol.DiagnosticSeverityWarning,
		diagnostics.Info:  protocol.DiagnosticSeverityInformation,
		diagnostics.Hint:  protocol.DiagnosticSeverityHint,
	}

	for in, expected := range cases {
		if got := toLSPSeverity(in); got != expected {
			t.Fatalf("severity %v: expected %v, got %v", in, expected, got)
		}
	}
}

func TestUriToPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix paths")
	}

	cases := map[string]string{
		"file:///home/me/x.vp":              "/home/me/x.vp",
		"file:///home/me/my%20project/x.vp": "/home/me/my project/x.vp",
		"file:///home/me/caf%C3%A9/../x.vp": "/home/me/x.vp",
	}

	for uri, expected := range cases {
		got, err := uriToPath(protocol.DocumentUri(uri))
		if err != nil {
			t.Fatalf("%v: unexpected error %v", uri, err)
		}
		if got != expected {
			t.Fatalf("%v: expected %v, got %v", uri, expected, got)
		}
	}

	if _, err := uriToPath("untitled:Untitled-1"); err == nil {
		t.Fatal("expected an error for non file URI")
	}
}

func TestAddError(t *testing.T) {
	file := filepath.Clean("/project/x.vp")
	all := []string{"fatal", "warn", "info", "hint"}

	ds := addError(
		[]protocol.Diagnostic{},
		diagnostics.Diagnostics{
			diagnostics.NewError(token.Item{File: file, Line: 2, Char: 5, Value: "abc"}, "error"),
			// value longer than char: start must not underflow
			diagnostics.NewWarning(token.Item{File: file, Line: 0, Char: 1, Value: "abc"}, "warn"),
			// other file: dropped
			diagnostics.NewError(token.Item{File: "/project/y.vp", Line: 0, Char: 1}, "other"),
		},
		file,
		all,
	)

	if len(ds) != 2 {
		t.Fatalf("expected 2 diagnostics, got %v", len(ds))
	}

	if *ds[0].Severity != protocol.DiagnosticSeverityError || ds[0].Range.Start.Character != 2 {
		t.Fatalf("unexpected first diagnostic %+v", ds[0])
	}

	if *ds[1].Severity != protocol.DiagnosticSeverityWarning || ds[1].Range.Start.Character != 0 {
		t.Fatalf("unexpected second diagnostic %+v", ds[1])
	}

	// filtered by configured severity
	ds = addError(
		[]protocol.Diagnostic{},
		diagnostics.Diagnostics{
			diagnostics.NewWarning(token.Item{File: file}, "warn"),
		},
		file,
		[]string{"fatal"},
	)

	if len(ds) != 0 {
		t.Fatalf("expected warnings to be filtered out, got %v", len(ds))
	}
}
