package main

import "github.com/tliron/glsp/server"

func (d *doctor) run() {
	d.server = server.NewServer(d.handler, d.name, false)
	d.server.RunStdio()
}
