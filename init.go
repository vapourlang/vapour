package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (d *Doctor) init() {
	d.handler = &protocol.Handler{
		Initialize:  d.initialize,
		Initialized: d.initialized,
		Shutdown:    d.shutdown,
		SetTrace:    d.setTrace,
	}
}

func (d *Doctor) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := d.handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    d.name,
			Version: &d.version,
		},
	}, nil
}

func (d *Doctor) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (d *Doctor) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (d *Doctor) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
