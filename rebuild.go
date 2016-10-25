package main

import (
	"github.com/mkideal/cli"
	"github.com/penberg/go-scylla-api/scylla"
)

var _ = app.Register(&cli.Command{
	Name: "rebuild",
	Desc: "Rebuild data by streaming from other nodes",
	Argv: func() interface{} { return new(rebuildT) },
	Fn:   rebuild,
})

type rebuildT struct {
	cli.Helper
	Host     string `cli:"host" usage:"Node hostname or IP address" dft:"localhost"`
	Port     string `cli:"port" usage:"API server port number" dft:"10000"`
	SourceDC string `cli:"source-dc" usage:"Name of datacenter to stream from. If not specified, picks any datacenter."`
}

func rebuild(ctx *cli.Context) error {
	argv := ctx.Argv().(*rebuildT)
	client := scylla.NewClient(argv.Host, argv.Port)
	return client.Rebuild(argv.SourceDC)
}
