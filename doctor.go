package main

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

type Doctor struct {
	name    string
	version string
	handler *protocol.Handler
	server  *server.Server
}
