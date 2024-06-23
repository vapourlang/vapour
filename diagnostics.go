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
	var sev protocol.DiagnosticSeverity = 0
	code := protocol.IntegerOrString{
		Value: "Error from doctor",
	}

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.Diagnostic{
		Severity: &sev,
		Code:     &code,
		Message:  "There is an errror from doctor",
	})

	return nil
}

func (d *doctor) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	return nil
}
