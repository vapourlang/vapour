package main

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (v *vapour) init() {
	commonlog.NewInfoMessage(0, "Initialising doctor...")

	v.handler = &protocol.Handler{
		Initialize:            v.initialize,
		Initialized:           v.initialized,
		Shutdown:              v.shutdown,
		SetTrace:              v.setTrace,
		TextDocumentDidOpen:   v.textDocumentDidOpen,
		TextDocumentDidChange: v.textDocumentDidChange,
	}
}

func (v *vapour) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
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

func (v *vapour) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (v *vapour) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (v *vapour) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
