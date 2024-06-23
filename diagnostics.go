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
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Message: "Hello World Init",
		Type:    protocol.MessageTypeInfo,
	})
	return nil
}
