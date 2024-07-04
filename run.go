package main

import "github.com/tliron/glsp/server"

func (v *vapour) run() {
	v.server = server.NewServer(v.handler, v.name, false)
	v.server.RunStdio()
}
