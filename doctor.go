package main

import (
	"fmt"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

type RFile struct {
	path    string
	content []byte
}

type RFiles []RFile

type doctor struct {
	name    string
	version string
	handler *protocol.Handler
	server  *server.Server
	root    *string
	rfiles  RFiles
}

func (rfls RFiles) print() {
	fmt.Printf("found %v R files\n", len(rfls))
	for _, fl := range rfls {
		fmt.Printf("%v\n", fl.path)
	}
}

func (d *doctor) print() {
	fmt.Printf("%v - %v", d.name, d.version)
	d.rfiles.print()
}
