package main

import (
	"github.com/mkideal/cli"
	"github.com/penberg/go-scylla-api/scylla"
)

var _ = app.Register(&cli.Command{
	Name: "flush",
	Desc: "Flush one or more column families",
	Argv: func() interface{} { return new(flushT) },
	Fn:   flush,
})

type flushT struct {
	cli.Helper
	Host     string `cli:"host" usage:"Node hostname or IP address" dft:"localhost"`
	Port     string `cli:"port" usage:"API server port number" dft:"10000"`
	Keyspace string `cli:"*k,keyspace" usage:"Keyspace to flush"`
	ColumnFamilies []string `cli:"c,columnfamily" usage:"Column families to flush"`
}

func flush(ctx *cli.Context) error {
	argv := ctx.Argv().(*flushT)
	client := scylla.NewClient(argv.Host, argv.Port)
	return client.Flush(argv.Keyspace, argv.ColumnFamilies)
}
