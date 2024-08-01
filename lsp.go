package main

import (
	"log"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

func (v *vapour) lspInit() {
	commonlog.NewInfoMessage(0, "Initialising vapour...")

	v.handler = &protocol.Handler{
		Initialize:            v.lspInitialize,
		Initialized:           v.lspInitialized,
		Shutdown:              v.lspShutdown,
		SetTrace:              v.lspSetTrace,
		TextDocumentDidOpen:   v.lspTextDocumentDidOpen,
		TextDocumentDidChange: v.lspTextDocumentDidChange,
	}
}

func (v *vapour) lspInitialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	v.root = params.RootPath
	capabilities := v.handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    v.name,
			Version: &v.version,
		},
	}, nil
}

func (v *vapour) lspInitialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (v *vapour) lspShutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (v *vapour) lspSetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (v *vapour) lspRun() {
	v.server = server.NewServer(v.handler, v.name, false)
	v.server.RunStdio()
}

func (v *vapour) lspTextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: "Hello doctor!",
		Type:    protocol.MessageTypeInfo,
	})
	return nil
}

func (v *vapour) lspTextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	var sev protocol.DiagnosticSeverity = 1
	var code protocol.IntegerOrString = protocol.IntegerOrString{Value: 2}
	src := "Doctor!"

	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: *v.root,
		Type:    protocol.MessageTypeInfo,
	})

	err := v.readDir()

	if err != nil {
		log.Fatal("error reading files")
	}

	ds := protocol.PublishDiagnosticsParams{
		URI: params.TextDocument.URI,
		Diagnostics: []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 2,
					},
					End: protocol.Position{
						Line:      0,
						Character: 4,
					},
				},
				Severity: &sev,
				Code:     &code,
				Source:   &src,
				Message:  "Diagnostic error from doctor",
			},
		},
	}

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
	return nil
}
