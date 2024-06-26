package main

import (
	"log"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (d *doctor) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: "Hello doctor!",
		Type:    protocol.MessageTypeInfo,
	})
	return nil
}

func (d *doctor) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	var sev protocol.DiagnosticSeverity = 1
	var code protocol.IntegerOrString = protocol.IntegerOrString{Value: 2}
	src := "Doctor!"

	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: *d.root,
		Type:    protocol.MessageTypeInfo,
	})

	err := d.readRoot()

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
