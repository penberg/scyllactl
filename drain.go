package main

import (
	"github.com/mkideal/cli"
	"github.com/penberg/go-scylla-api/scylla"
)

var _ = app.Register(&cli.Command{
	Name: "drain",
	Desc: "Drain the node (stop accepting writes and flush all column families)",
	Argv: func() interface{} { return new(drainT) },
	Fn:   drain,
})

type drainT struct {
	cli.Helper
	Host string `cli:"host" usage:"Node hostname or IP address" dft:"localhost"`
	Port string `cli:"port" usage:"API server port number" dft:"10000"`
}

func drain(ctx *cli.Context) error {
	argv := ctx.Argv().(*drainT)
	client := scylla.NewClient(argv.Host, argv.Port)
	return client.Drain()
}
