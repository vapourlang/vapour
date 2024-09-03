package lsp

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/walker"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

type LSP struct {
	files []lexer.File
	conf  *config
}

type walkParams struct {
	TextDocument protocol.DocumentUri
}

type config struct {
	When     []string `json:"when"`
	Severity []string `json:"severity"`
}

var handler protocol.Handler
var version string = "0.0.1"
var code protocol.IntegerOrString = protocol.IntegerOrString{Value: 2}
var src string = "Vapour"
var file string = ".vapour"

func makeConfigPath(conf string) string {
	dirname, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return path.Join(dirname, conf)
}

func hasConfig(conf string) bool {
	p := makeConfigPath(conf)

	if p == "" {
		return false
	}

	_, err := os.Stat(p)

	return !os.IsNotExist(err)
}

func readConfig(conf string) (*config, bool) {
	configuration := &config{}

	if !hasConfig(conf) {
		return configuration, false
	}

	p := makeConfigPath(conf)
	data, err := os.ReadFile(p)

	if err != nil {
		return configuration, false
	}

	err = json.Unmarshal(data, configuration)

	if err != nil {
		return configuration, false
	}

	return configuration, true
}

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
		files: []lexer.File{},
	}
}

func Run(tcp bool, port string) {
	l := New()

	conf, ok := readConfig(file)

	l.conf = conf

	handler = protocol.Handler{
		Initialize:  l.initialize,
		Initialized: l.initialized,
		Shutdown:    l.shutdown,
		SetTrace:    l.setTrace,
	}

	if !ok {
		l.conf = &config{
			When:     []string{"open", "save", "close", "text"},
			Severity: []string{"fatal", "warn", "info", "hint"},
		}
	}

	if contains("open", conf.When) {
		handler.TextDocumentDidOpen = l.textDocumentDidOpen
	}

	if contains("save", conf.When) {
		handler.TextDocumentDidSave = l.textDocumentDidSave
	}

	if contains("close", conf.When) {
		handler.TextDocumentDidClose = l.textDocumentDidClose
	}

	if contains("change", conf.When) {
		handler.TextDocumentDidChange = l.textDocumentDidChange
	}

	server := server.NewServer(&handler, "Vapour", false)

	var err error

	if tcp {
		err = server.RunTCP(port)

		if err != nil {
			panic(err)
		}
	}

	err = server.RunStdio()

	if err != nil {
		panic(err)
	}
}

func (l *LSP) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	// incremental changes + whole document on open
	capabilities.TextDocumentSync = 1

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
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: "Vapour initialised",
		Type:    protocol.MessageTypeInfo,
	})
	return nil
}

func (l *LSP) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (l *LSP) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (l *LSP) walkFiles(context *glsp.Context, params *walkParams) error {
	diagnostics := []protocol.Diagnostic{}

	// remove that in the future to leverage workspace
	// and only process files once
	l.files = []lexer.File{}

	// read directory
	file := strings.Replace(params.TextDocument, "file://", "", 1)
	root := filepath.Dir(file)
	err := l.readDir(root)

	if err != nil {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Message: fmt.Sprintf("Error reading files: %v", err.Error()),
			Type:    protocol.MessageTypeError,
		})
		return err
	}

	// lex
	le := lexer.New(l.files)
	le.Run()

	if le.HasError() {
		diagnostics = addError(diagnostics, le.Errors, file, l.conf.Severity)
		ds := protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument,
			Diagnostics: diagnostics,
		}
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
		return nil
	}

	// parse
	p := parser.New(le)
	prog := p.Run()

	if p.HasError() {
		diagnostics = addError(diagnostics, p.Errors(), file, l.conf.Severity)
		ds := protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument,
			Diagnostics: diagnostics,
		}
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
		return nil
	}

	// walk tree
	w := walker.New()
	w.Walk(prog)

	diagnostics = addError(diagnostics, w.Errors(), file, l.conf.Severity)
	ds := protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument,
		Diagnostics: diagnostics,
	}
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
	return nil
}

func (l *LSP) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	p := &walkParams{
		TextDocument: params.TextDocument.URI,
	}
	return l.walkFiles(context, p)
}

func (l *LSP) textDocumentDidSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	p := &walkParams{
		TextDocument: params.TextDocument.URI,
	}
	return l.walkFiles(context, p)
}

func (l *LSP) textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	p := &walkParams{
		TextDocument: params.TextDocument.URI,
	}
	return l.walkFiles(context, p)
}

func (l *LSP) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	p := &walkParams{
		TextDocument: params.TextDocument.URI,
	}
	return l.walkFiles(context, p)
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

func addError(ds []protocol.Diagnostic, ns diagnostics.Diagnostics, file string, include []string) []protocol.Diagnostic {
	for _, e := range ns {
		if e.Token.File != file {
			continue
		}

		if !validSeverity(e.Severity, include) {
			continue
		}

		s := protocol.DiagnosticSeverity(e.Severity)

		ds = append(
			ds,
			protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(e.Token.Line),
						Character: uint32(e.Token.Char - len(e.Token.Value)),
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
