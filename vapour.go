package main

import (
	"fmt"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

type vpFile struct {
	path    string
	content []byte
}

type vpFiles []vpFile

type vapour struct {
	name    string
	version string
	handler *protocol.Handler
	server  *server.Server
	root    *string
	files   vpFiles
}

func (fls vpFiles) print() {
	fmt.Printf("found %v .vp files\n", len(fls))
	for _, fl := range fls {
		fmt.Printf("%v\n", fl.path)
	}
}

func (d *vapour) print() {
	fmt.Printf("%v - %v", d.name, d.version)
	d.files.print()
}
