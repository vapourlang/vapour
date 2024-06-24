package main

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (d *doctor) init() {
	commonlog.NewInfoMessage(0, "Initialising doctor...")

	d.handler = &protocol.Handler{
		Initialize:            d.initialize,
		Initialized:           d.initialized,
		Shutdown:              d.shutdown,
		SetTrace:              d.setTrace,
		TextDocumentDidOpen:   d.textDocumentDidOpen,
		TextDocumentDidChange: d.textDocumentDidChange,
	}
}

func (d *doctor) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := d.handler.CreateServerCapabilities()

	capabilities.TypeDefinitionProvider = &protocol.TypeDefinitionClientCapabilities{}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    d.name,
			Version: &d.version,
		},
	}, nil
}

func (d *doctor) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (d *doctor) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (d *doctor) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
