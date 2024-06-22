package main

import protocol "github.com/tliron/glsp/protocol_3_16"

type Doctor struct {
	name     string
	protocol *protocol.Handler
}
