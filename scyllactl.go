package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

const version = "v0.0.1"

var app = &cli.Command{
	Name: os.Args[0],
	Desc: "scyllactl",
	Text: `scyllactl is a tool for managing Scylla clusters`,
	Argv: func() interface{} { return new(scyllactlT) },
	Fn:   scyllactl,
}

type scyllactlT struct {
	cli.Helper
	Version bool `cli:"v,version" usage:"display version"`
}

func scyllactl(ctx *cli.Context) error {
	argv := ctx.Argv().(*scyllactlT)
	if argv.Version {
		ctx.String(version + "\n")
		return nil
	}
	ctx.String("try `%s --help for more information'\n", ctx.Path())
	return nil
}

func main() {
	cli.SetUsageStyle(cli.ManualStyle)
	if err := app.RunWith(os.Args[1:], os.Stderr, nil); err != nil {
		fmt.Printf("%v\n", err)
	}
}
