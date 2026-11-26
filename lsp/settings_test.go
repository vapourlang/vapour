package lsp

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vapourlang/vapour/config"
)

func newTestLSP(when ...string) *LSP {
	l := New()
	l.conf = config.ReadConfig()
	l.conf.Lsp.When = when
	return l
}

func TestEnabled(t *testing.T) {
	l := newTestLSP("open", "text")

	if !l.enabled("open") || l.enabled("save") {
		t.Fatal("unexpected open/save")
	}

	// "text" is the legacy name of "change"
	if !l.enabled("change") {
		t.Fatal("expected text to enable change")
	}
}

func TestApplyInitOptions(t *testing.T) {
	l := newTestLSP("open", "save")
	severity := l.conf.Lsp.Severity

	// decoded JSON as received from the client
	l.applyInitOptions(map[string]any{"when": []any{"save"}})

	if l.enabled("open") || !l.enabled("save") {
		t.Fatalf("when not overridden: %v", l.conf.Lsp.When)
	}

	if len(l.conf.Lsp.Severity) != len(severity) {
		t.Fatal("severity must be kept when not passed")
	}

	// invalid options are ignored
	l.applyInitOptions("nope")
	l.applyInitOptions(nil)

	if !l.enabled("save") {
		t.Fatal("invalid options changed the configuration")
	}
}

func TestChangeIsDebounced(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "x.vp")
	if err := os.WriteFile(file, []byte("let x: int = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	uri := protocol.DocumentUri("file://" + filepath.ToSlash(file))

	var mu sync.Mutex
	published := 0
	ctx := &glsp.Context{
		Notify: func(method string, params any) {
			if method == protocol.ServerTextDocumentPublishDiagnostics {
				mu.Lock()
				published++
				mu.Unlock()
			}
		},
	}

	l := newTestLSP("change")
	for i := 0; i < 5; i++ {
		err := l.textDocumentDidChange(ctx, &protocol.DidChangeTextDocumentParams{
			TextDocument: protocol.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: protocol.TextDocumentIdentifier{URI: uri},
			},
			ContentChanges: []any{
				protocol.TextDocumentContentChangeEventWhole{Text: "let x: int = 1\n"},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(changeDelay * 3)

	mu.Lock()
	defer mu.Unlock()
	if published != 1 {
		t.Fatalf("expected a single check for a burst of edits, got %v", published)
	}
}
