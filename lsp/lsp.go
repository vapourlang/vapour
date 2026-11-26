package lsp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/vapourlang/vapour/config"
	"github.com/vapourlang/vapour/diagnostics"
	"github.com/vapourlang/vapour/lexer"
	"github.com/vapourlang/vapour/parser"
	"github.com/vapourlang/vapour/walker"
)

var src string = "Vapour"

type LSP struct {
	files []lexer.File
	conf  *config.Config

	// content of documents open in the editor, by path
	// it may differ from the content on disk (unsaved changes)
	mu        sync.Mutex
	documents map[string]string

	// pending checks triggered by edits, by document
	// they are delayed so we do not check on every keystroke
	tmu    sync.Mutex
	timers map[protocol.DocumentUri]*time.Timer
}

// how long to wait after the last edit before checking
var changeDelay = 300 * time.Millisecond

// settings the client can pass in initializationOptions
// they override the user's configuration file
type initOptions struct {
	When     []string `json:"when"`
	Severity []string `json:"severity"`
}

type walkParams struct {
	TextDocument protocol.DocumentUri
}

var handler protocol.Handler
var version string = "0.0.1"
var code protocol.IntegerOrString = protocol.IntegerOrString{Value: 2}

func contains(value string, values []string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

func New() *LSP {
	return &LSP{
		files:     []lexer.File{},
		documents: make(map[string]string),
		timers:    make(map[protocol.DocumentUri]*time.Timer),
	}
}

func Run(conf *config.Config, tcp bool, port string) {
	l := New()

	l.conf = conf

	// handlers are always registered so we keep track of the
	// documents' content, `when` only decides when we check them
	// (it can be changed by the client on initialize)
	handler = protocol.Handler{
		Initialize:            l.initialize,
		Initialized:           l.initialized,
		Shutdown:              l.shutdown,
		SetTrace:              l.setTrace,
		TextDocumentDidOpen:   l.textDocumentDidOpen,
		TextDocumentDidSave:   l.textDocumentDidSave,
		TextDocumentDidClose:  l.textDocumentDidClose,
		TextDocumentDidChange: l.textDocumentDidChange,
	}

	server := server.NewServer(&handler, "Vapour", false)

	var err error

	if tcp {
		err = server.RunTCP(port)

		if err != nil {
			panic(err)
		}

		return
	}

	err = server.RunStdio()

	if err != nil {
		panic(err)
	}
}

func (l *LSP) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	l.applyInitOptions(params.InitializationOptions)

	capabilities := handler.CreateServerCapabilities()

	// the client sends the whole document on open and on change
	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindFull

	var supported bool = true

	capabilities.Workspace = &protocol.ServerCapabilitiesWorkspace{
		WorkspaceFolders: &protocol.WorkspaceFoldersServerCapabilities{
			Supported: &supported,
		},
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "vapour",
			Version: &version,
		},
	}, nil
}

func (l *LSP) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

// applyInitOptions overrides the configuration with the
// client's initializationOptions, e.g.: {"when": ["save"]}
func (l *LSP) applyInitOptions(raw any) {
	if raw == nil {
		return
	}

	data, err := json.Marshal(raw)

	if err != nil {
		return
	}

	var opts initOptions
	if err := json.Unmarshal(data, &opts); err != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if opts.When != nil {
		l.conf.Lsp.When = opts.When
	}

	if opts.Severity != nil {
		l.conf.Lsp.Severity = opts.Severity
	}
}

// enabled reports whether documents are checked on event
func (l *LSP) enabled(event string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if contains(event, l.conf.Lsp.When) {
		return true
	}

	// "text" is accepted for backward compatibility
	// it used to be the default but was never matched
	return event == "change" && contains("text", l.conf.Lsp.When)
}

// cancelCheck drops a pending check of the document
func (l *LSP) cancelCheck(uri protocol.DocumentUri) {
	l.tmu.Lock()
	defer l.tmu.Unlock()

	if t, ok := l.timers[uri]; ok {
		t.Stop()
		delete(l.timers, uri)
	}
}

// scheduleCheck checks the document once edits have stopped
// for changeDelay, each new edit restarts the wait
func (l *LSP) scheduleCheck(context *glsp.Context, uri protocol.DocumentUri) {
	l.tmu.Lock()
	defer l.tmu.Unlock()

	if t, ok := l.timers[uri]; ok {
		t.Stop()
	}

	var t *time.Timer
	t = time.AfterFunc(changeDelay, func() {
		l.tmu.Lock()
		// a newer edit replaced this timer
		if l.timers[uri] != t {
			l.tmu.Unlock()
			return
		}
		delete(l.timers, uri)
		l.tmu.Unlock()

		l.walkFiles(context, &walkParams{TextDocument: uri})
	})
	l.timers[uri] = t
}

func (l *LSP) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (l *LSP) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (l *LSP) walkFiles(context *glsp.Context, params *walkParams) (err error) {
	// requests are handled one at a time on a single connection
	// but TCP serves multiple connections
	l.mu.Lock()
	defer l.mu.Unlock()

	// a bug in vapour must not bring the language server down
	defer func() {
		if r := recover(); r != nil {
			context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
				Message: fmt.Sprintf("Vapour failed to check the file: %v", r),
				Type:    protocol.MessageTypeError,
			})
			err = nil
		}
	}()

	diagnostics := []protocol.Diagnostic{}

	// remove that in the future to leverage workspace
	// and only process files once
	l.files = []lexer.File{}

	// read directory
	file, err := uriToPath(params.TextDocument)

	if err != nil {
		return err
	}

	root := filepath.Dir(file)
	err = l.readDir(root)

	if err != nil {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Message: fmt.Sprintf("Error reading files: %v", err.Error()),
			Type:    protocol.MessageTypeError,
		})
		return err
	}

	publish := func() {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument,
			Diagnostics: diagnostics,
		})
	}

	// lex
	le := lexer.New(l.files)
	le.Run()

	if le.HasError() {
		diagnostics = addError(diagnostics, le.Errors(), file, l.conf.Lsp.Severity)
		publish()
		return nil
	}

	// parse
	p := parser.New(le)
	prog := p.Run()

	if p.HasError() {
		diagnostics = addError(diagnostics, p.Errors(), file, l.conf.Lsp.Severity)
		publish()
		return nil
	}

	// walk tree
	w := walker.New()
	w.Walk(prog)

	diagnostics = addError(diagnostics, w.Errors(), file, l.conf.Lsp.Severity)
	publish()
	return nil
}

// uriToPath converts a file URI to a path on disk.
// e.g.: file:///home/me/my%20project/x.vp -> /home/me/my project/x.vp
// e.g.: file:///c%3A/project/x.vp -> c:\project\x.vp
func uriToPath(uri protocol.DocumentUri) (string, error) {
	u, err := url.Parse(string(uri))

	if err != nil {
		return "", fmt.Errorf("invalid document URI %v: %w", uri, err)
	}

	if u.Scheme != "" && u.Scheme != "file" {
		return "", fmt.Errorf("unsupported URI scheme %v", u.Scheme)
	}

	path := u.Path

	// windows: /c:/project/x.vp -> c:/project/x.vp
	if runtime.GOOS == "windows" && len(path) > 2 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}

	return filepath.Clean(filepath.FromSlash(path)), nil
}

func (l *LSP) setDocument(uri protocol.DocumentUri, content string) {
	path, err := uriToPath(uri)

	if err != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.documents[path] = content
}

func (l *LSP) removeDocument(uri protocol.DocumentUri) {
	path, err := uriToPath(uri)

	if err != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.documents, path)
}

func (l *LSP) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	l.setDocument(params.TextDocument.URI, params.TextDocument.Text)

	if !l.enabled("open") {
		return nil
	}

	p := &walkParams{
		TextDocument: params.TextDocument.URI,
	}
	return l.walkFiles(context, p)
}

func (l *LSP) textDocumentDidSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	if params.Text != nil {
		l.setDocument(params.TextDocument.URI, *params.Text)
	}

	if !l.enabled("save") {
		return nil
	}

	// the save supersedes a pending check
	l.cancelCheck(params.TextDocument.URI)

	p := &walkParams{
		TextDocument: params.TextDocument.URI,
	}
	return l.walkFiles(context, p)
}

func (l *LSP) textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	// once closed the content on disk is the reference
	l.removeDocument(params.TextDocument.URI)
	l.cancelCheck(params.TextDocument.URI)

	if !l.enabled("close") {
		return nil
	}

	p := &walkParams{
		TextDocument: params.TextDocument.URI,
	}
	return l.walkFiles(context, p)
}

func (l *LSP) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	// we use full sync: the last change holds the whole document
	for _, change := range params.ContentChanges {
		switch c := change.(type) {
		case protocol.TextDocumentContentChangeEventWhole:
			l.setDocument(params.TextDocument.URI, c.Text)
		}
	}

	if !l.enabled("change") {
		return nil
	}

	l.scheduleCheck(context, params.TextDocument.URI)
	return nil
}

func validSeverity(severity diagnostics.Severity, valid []string) bool {
	if severity == 0 && contains("fatal", valid) {
		return true
	}

	if severity == 1 && contains("warn", valid) {
		return true
	}

	if severity == 2 && contains("hint", valid) {
		return true
	}

	if severity == 3 && contains("info", valid) {
		return true
	}

	return false
}

// toLSPSeverity maps vapour severities to LSP severities
// vapour: Fatal = 0, Warn = 1, Hint = 2, Info = 3
// LSP: Error = 1, Warning = 2, Information = 3, Hint = 4
func toLSPSeverity(severity diagnostics.Severity) protocol.DiagnosticSeverity {
	switch severity {
	case diagnostics.Fatal:
		return protocol.DiagnosticSeverityError
	case diagnostics.Warn:
		return protocol.DiagnosticSeverityWarning
	case diagnostics.Info:
		return protocol.DiagnosticSeverityInformation
	case diagnostics.Hint:
		return protocol.DiagnosticSeverityHint
	}

	return protocol.DiagnosticSeverityError
}

func addError(ds []protocol.Diagnostic, ns diagnostics.Diagnostics, file string, include []string) []protocol.Diagnostic {
	for _, e := range ns {
		if filepath.Clean(e.Token.File) != file {
			continue
		}

		if !validSeverity(e.Severity, include) {
			continue
		}

		s := toLSPSeverity(e.Severity)

		// the token's character is where it ends
		start := e.Token.Char - len(e.Token.Value)
		if start < 0 {
			start = 0
		}

		ds = append(
			ds,
			protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(e.Token.Line),
						Character: uint32(start),
					},
					End: protocol.Position{
						Line:      uint32(e.Token.Line),
						Character: uint32(e.Token.Char),
					},
				},
				Severity: &s,
				Code:     &code,
				Source:   &src,
				Message:  e.Message,
			},
		)
	}
	return ds
}
