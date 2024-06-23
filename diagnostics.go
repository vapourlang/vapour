package main

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (d *doctor) publishDiagnostics(context *glsp.Context, params *protocol.PublishDiagnosticsParams) (interface{}, error) {
	var diagnostics []protocol.Diagnostic

	return diagnostics, nil
}

func (d *doctor) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	commonlog.NewInfoMessage(0, "Document open...")
	return nil
}

func (d *doctor) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	return nil
}
