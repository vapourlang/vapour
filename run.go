package main

import "github.com/tliron/glsp/server"

func (d *Doctor) run() {
	d.server = server.NewServer(d.handler, d.name, false)
}
