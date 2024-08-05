package lsp

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/walker"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

var handler protocol.Handler
var version string = "0.0.1"
var fatal protocol.DiagnosticSeverity = 1
var code protocol.IntegerOrString = protocol.IntegerOrString{Value: 2}
var src string = "Vapour"

func Run() {
	handler = protocol.Handler{
		Initialize:            initialize,
		Initialized:           initialized,
		Shutdown:              shutdown,
		SetTrace:              setTrace,
		TextDocumentDidOpen:   textDocumentDidOpen,
		TextDocumentDidChange: textDocumentDidChange,
	}

	server := server.NewServer(&handler, "Vapour", false)
	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "vapour",
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: "Vapour LSP running",
		Type:    protocol.MessageTypeInfo,
	})
	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: "Text document changed",
		Type:    protocol.MessageTypeInfo,
	})
	var diagnostics []protocol.Diagnostic

	file := strings.Replace(params.TextDocument.URI, "file://", "", 1)

	// check that it is a vapour file
	r, _ := regexp.Compile("vp$")
	if !r.MatchString(file) {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Message: "Not a vapour file",
			Type:    protocol.MessageTypeInfo,
		})
		return nil
	}

	content, err := os.ReadFile(file)

	if err != nil {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Message: fmt.Sprintf("Error reading file: %v", file),
			Type:    protocol.MessageTypeInfo,
		})
		return err
	}

	// lex
	l := lexer.NewCode(file, string(content))
	l.Run()

	if l.HasError() {
		diagnostics = addError(diagnostics, l.Errors)
		ds := protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: diagnostics,
		}

		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
		return nil
	}

	// parse
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		diagnostics = addError(diagnostics, p.Errors())
		ds := protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: diagnostics,
		}

		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
		return nil
	}

	// walk tree
	w := walker.New()
	w.Walk(prog)

	if w.HasError() {
		diagnostics = addError(diagnostics, w.Errors())
		ds := protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: diagnostics,
		}

		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
		return nil
	}

	return nil
}

func addError(ds []protocol.Diagnostic, ns diagnostics.Diagnostics) []protocol.Diagnostic {
	for _, e := range ns {
		ds = append(
			ds,
			protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(e.Token.Line + 1),
						Character: uint32(e.Token.Char + 1),
					},
					End: protocol.Position{
						Line:      uint32(e.Token.Line + 1),
						Character: uint32(e.Token.Char + len(e.Token.Value)),
					},
				},
				Severity: &fatal,
				Code:     &code,
				Source:   &src,
				Message:  e.Message,
			},
		)
	}

	return ds
}
