package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (d *doctor) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: "Hello World",
		Type:    protocol.MessageTypeInfo,
	})
	return nil
}

func (d *doctor) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	//context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
	//	Message: "Hello World Init",
	//	Type:    protocol.MessageTypeInfo,
	//})

	var s protocol.DiagnosticSeverity = 1
	code := protocol.IntegerOrString{Value: "Doctor code"}

	ds := protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 1,
			},
			End: protocol.Position{
				Line:      0,
				Character: 3,
			},
		},
		Severity: &s,
		Code:     &code,
	}

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, ds)
	return nil
}
